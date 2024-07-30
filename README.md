# Финальный проект TODO-LIST

Выполнены все задания со звездочкой, кроме создания докер образа.

Структура проекта
TODO-project/
├── cmd/
│   └── main.go
├── config/
│   └── config.go
├── internal/
│   ├──api/
│   │   ├── handlers/
│   │   │      └── handlers.go
│   │   └── middleware
│   │         └── middleware.go
│   ├── database/
│   │   ├── database.go
│   │   └── scheduler.db
│   ├── hasher/
│   │   └── hasher.go
│   ├── parser/
│   │   └── parser.go
│   └── models/
│       └── models.go
├── go.mod
└── go.sum

Для корректной работы укажите переменные окружения:
TODO_DBFILE - абсолютный путь к вашей БД. Например- C:\Users\user\scheduler.db.
TODO_PORT - порт, на котором будет создан локальный сервер. Например- 8080.
TODO_PASSWORD - пароль, который у вас запросят при входе.

Аналогичные данные нужно прописать в файле settings.go для запуска тестов.