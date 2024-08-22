package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/shev-dm/TODO-project/internal/database"
	"github.com/shev-dm/TODO-project/internal/hasher"
	"github.com/shev-dm/TODO-project/internal/models"
	"github.com/shev-dm/TODO-project/internal/parser"
)

type Handler struct {
	Store *database.Storage
}

func (h *Handler) GetNextDate(w http.ResponseWriter, r *http.Request) {
	dateNowString := r.FormValue("now")
	dateNow, err := time.Parse("20060102", dateNowString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := parser.NextDate(dateNow, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(nextDate))
}

func (h *Handler) PostTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var input models.Task
	var errorAnswer models.Err

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorAnswer.Err = err.Error()
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}

	input, errorAnswer = parser.CheckRulesAddOrUpdate(input)
	if errorAnswer.Err != "" {
		w.WriteHeader(http.StatusBadRequest)
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}

	id, err := h.Store.Add(input)
	if err != nil {
		errorAnswer.Err = err.Error()
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}

	input.Id = strconv.Itoa(int(id))

	response, err := json.Marshal(input)
	if err != nil {
		http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tasks := models.Tasks{Tasks: make([]models.Task, 0, 20)}
	var errorAnswer models.Err

	search := r.FormValue("search")

	tasks, err := h.Store.SearchTasks(search)
	if err != nil {
		errorAnswer.Err = err.Error()
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}
	response, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var errorAnswer models.Err
	taskIdStr := r.FormValue("id")

	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		errorAnswer.Err = "неверный формат id задачи"
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}

	task, err := h.Store.Get(taskId)
	if err != nil {
		errorAnswer.Err = err.Error()
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)

}

func (h *Handler) PutTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task models.Task
	var errorAnswer models.Err

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		errorAnswer.Err = "невозможно распарсить данные"
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}

	task, errorAnswer = parser.CheckRulesAddOrUpdate(task)
	if errorAnswer.Err != "" {
		w.WriteHeader(http.StatusBadRequest)
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}
	taskId, err := strconv.Atoi(task.Id)
	if err != nil {
		errorAnswer.Err = "неверный формат id задачи"
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}

	rowsAffected, err := h.Store.Update(task, taskId)
	if err != nil {
		errorAnswer.Err = err.Error()
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusBadRequest)
		errorAnswer.Err = "данного id нет в БД"
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}

	// Если программа дошла до этого места, значит ошибок нет, возвращаем пустую errorAnswer
	response, err := json.Marshal(errorAnswer)
	if err != nil {
		http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)

}

func (h *Handler) PostTaskDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var errorAnswer models.Err

	taskIdStr := r.FormValue("id")
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		errorAnswer.Err = "неверный формат id задачи"
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}

	task, err := h.Store.Get(taskId)
	if err != nil {
		errorAnswer.Err = err.Error()
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}

	if task.Repeat == "" {
		err = h.Store.Delete(taskId)
		if err != nil {
			errorAnswer.Err = err.Error()
			response, err := json.Marshal(errorAnswer)
			if err != nil {
				http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(response)
			return
		}
	} else {
		nextDate, err := parser.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			errorAnswer.Err = err.Error()
			response, err := json.Marshal(errorAnswer)
			if err != nil {
				http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(response)
			return
		}
		task.Date = nextDate

		_, err = h.Store.Update(task, taskId)
		if err != nil {
			errorAnswer.Err = err.Error()
			response, err := json.Marshal(errorAnswer)
			if err != nil {
				http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(response)
			return
		}
	}
	response, err := json.Marshal(errorAnswer)
	if err != nil {
		http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var errorAnswer models.Err

	taskIdStr := r.FormValue("id")
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		errorAnswer.Err = "неверный формат id задачи"
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}
	err = h.Store.Delete(taskId)
	if err != nil {
		errorAnswer.Err = err.Error()
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}

	response, err := json.Marshal(errorAnswer)
	if err != nil {
		http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)

}

func (h *Handler) PostSignin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var password models.Authentication
	var errorAnswer models.Err

	if err := json.NewDecoder(r.Body).Decode(&password); err != nil {
		errorAnswer.Err = "невозможно распарсить данные"
		response, err := json.Marshal(errorAnswer)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(response)
		return
	}
	todoPassword := os.Getenv("TODO_PASSWORD")

	if todoPassword == password.Password {
		signedToken, err := hasher.GenerateToken(todoPassword)
		if err != nil {
			return
		}

		token := models.Token{Token: signedToken}
		response, err := json.Marshal(token)
		if err != nil {
			http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
		return
	}

	errorAnswer.Err = "неверный пароль"

	response, err := json.Marshal(errorAnswer)
	if err != nil {
		http.Error(w, "ошибка преобразования данных", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(response)

}
