package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/shev-dm/TODO-project/internal/models"
	_ "modernc.org/sqlite"
)

const maxTasksPerPage = 10

type Storage struct {
	db *sql.DB
}

func NewStorage(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	return &Storage{db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Init(appPath string) error {
	_, err := os.Stat(appPath)

	var install bool
	if err != nil {
		install = true
	}
	if !install {
		return nil
	}

	_, err = s.db.Exec("CREATE TABLE scheduler " +
		"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
		"date varchar(8)," +
		"title varchar(256)," +
		"comment TEXT," +
		"repeat varchar(128))")
	if err != nil {
		return err
	}

	_, err = s.db.Exec("CREATE INDEX idx_date ON scheduler (date)")
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Add(input models.Task) (int64, error) {
	var errorAnswer models.Err

	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment,repeat) values (:date, :title, :comment, :repeat)",
		sql.Named("date", input.Date),
		sql.Named("title", input.Title),
		sql.Named("comment", input.Comment),
		sql.Named("repeat", input.Repeat))
	if err != nil {
		errorAnswer.Err = "Ошибка при добавлении задачи"
		return 0, err
	}

	id, _ := res.LastInsertId()
	return id, nil
}

func (s *Storage) SearchTasks(search string) (models.Tasks, error) {
	tasks := models.Tasks{Tasks: make([]models.Task, 0, 20)}
	var row *sql.Rows
	var err error

	if search != "" {
		dateTime, err := time.Parse("02.01.2006", search)
		if err != nil {
			row, err = s.db.Query("SELECT id,date,title,comment,repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date ASC limit :limit",
				sql.Named("search", "%"+search+"%"),
				sql.Named("limit", maxTasksPerPage))
			if err != nil {
				return tasks, err
			}
		} else {
			dateTimeString := dateTime.Format("20060102")
			row, err = s.db.Query("SELECT id,date,title,comment,repeat FROM scheduler WHERE date LIKE :search ORDER BY date ASC limit :limit",
				sql.Named("search", "%"+dateTimeString+"%"),
				sql.Named("limit", maxTasksPerPage))
			if err != nil {
				return tasks, err
			}
		}
	} else {
		row, err = s.db.Query("SELECT id,date,title,comment,repeat from scheduler order by date ASC limit :limit",
			sql.Named("limit", maxTasksPerPage))
		if err != nil {
			return tasks, err
		}
	}

	defer row.Close()

	for row.Next() {
		var task models.Task
		err = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks.Tasks = append(tasks.Tasks, task)
	}
	if err = row.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}

func (s *Storage) Get(taskId int) (models.Task, error) {
	task := models.Task{}

	row := s.db.QueryRow("SELECT id,date,title,comment,repeat from scheduler where id = :id",
		sql.Named("id", taskId))
	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil || task.Id == "" {
		return models.Task{}, err
	}
	return task, nil
}

func (s *Storage) Update(task models.Task, taskId int) (int64, error) {

	result, err := s.db.Exec("UPDATE scheduler "+
		"SET date = :date, title = :title, comment = :comment, repeat = :repeat "+
		"where id = :id",
		sql.Named("id", taskId),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return rowsAffected, nil
}

func (s *Storage) Delete(taskId int) error {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", taskId))
	if err != nil {
		return err
	}
	return nil
}
