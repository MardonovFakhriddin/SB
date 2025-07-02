package main

import (
	"SB/internal/configs"
	"SB/internal/controller"
	"SB/internal/db"
	"SB/logger"
	"log"
)

func main() {
	// Reading configs
	if err := configs.ReadSettings(); err != nil {
		log.Fatalf("Ошибка чтения настроек: %v", err)
	}

	// Initializing logger
	if err := logger.Init(); err != nil {
		log.Fatalf("Ошибка инициализации логера: %v", err)
	}
	logger.Info.Println("Loggers initialized successfully!")

	// Connecting to db
	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Ошибка подключения к бд: %v", err)
	}
	logger.Info.Println("Connection to database established successfully!")

	// Initializing db-migrations
	if err := db.InitMigrations(); err != nil {
		log.Fatalf("Ошибка миграции к бд: %v", err)
	}
	logger.Info.Println("Migrations initialized successfully!")

	// Running http-server
	if err := controller.RunServer(); err != nil {
		return
	}

}
