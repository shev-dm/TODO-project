package config

import (
	"log"
	"os"
)

type Config struct {
	Port   string
	DBFile string
}

func NewConfig() *Config {
	c := Config{}
	c.init()
	return &c
}

func (c *Config) init() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		log.Fatal("не указана переменная окружения TODO_DBFILE")
	}
	c.DBFile = dbFile
	port := os.Getenv("TODO_PORT")
	if port == "" {
		log.Fatal("не указана переменная окружения TODO_PORT")
	}
	c.Port = ":" + port
}
