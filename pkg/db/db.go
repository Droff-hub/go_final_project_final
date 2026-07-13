package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// schema содержит SQL-запросы для создания таблицы и индекса
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(256) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`
// InitDB открывает базу данных и создаёт таблицу, если файл не существует.
func InitDB(dbFile string) error {
	// Проверяем, существует ли файл
	_, err := os.Stat(dbFile)
	needCreate := os.IsNotExist(err) // если файла нет, needCreate = true

	// Открываем базу данных (файл создастся, если его нет)
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

		// Проверяем подключение
	if err := db.Ping(); err != nil {
		return err
	}

	// Если файл не существовал, выполняем schema
	if needCreate {
		log.Printf("Файл БД %s не найден, создаём таблицу...", dbFile)
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
		log.Println("Таблица scheduler создана успешно")
	}

	DB = db
	log.Printf("Подключение к БД %s установлено", dbFile)
	return nil
}
