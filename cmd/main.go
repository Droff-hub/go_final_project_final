package main

import (
	"log"
	"os"

	"github.com/Droff-hub/go_final_project/pkg/db"
	"github.com/Droff-hub/go_final_project/pkg/server"
)

func main() {
	// Определяем путь к БД (стандартно scheduler.db)
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.InitDB(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.DB.Close()

	log.Println("Запуск планировщика задач...")
	server.Run()
}
