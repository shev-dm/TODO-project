package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shev-dm/TODO-project/config"
	"github.com/shev-dm/TODO-project/internal/api/handlers"
	"github.com/shev-dm/TODO-project/internal/api/middleware"
	"github.com/shev-dm/TODO-project/internal/database"
)

func main() {
	initConfig := config.NewConfig()
	store, err := database.NewStorage(initConfig.DBFile)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	err = store.Init(initConfig.DBFile)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.Handler{Store: store}

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authentication(initConfig.Password))
		r.Post("/api/task", handler.PostTask)
		r.Get("/api/tasks", handler.GetTasks)
		r.Get("/api/task", handler.GetTask)
		r.Put("/api/task", handler.PutTask)
		r.Post("/api/task/done", handler.PostTaskDone)
		r.Delete("/api/task", handler.DeleteTask)
	})
	r.Get("/api/nextdate", handler.GetNextDate)
	r.Post("/api/signin", handler.PostSignin)

	r.Handle("/*", http.FileServer(http.Dir("./web")))

	log.Printf("Приложение запущено на порту %s\n", initConfig.Port)

	if err = http.ListenAndServe(initConfig.Port, r); err != nil {
		log.Fatal(err)
	}
}
