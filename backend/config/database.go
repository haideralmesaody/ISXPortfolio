package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./data/isxportfolio.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	log.Println("Database initialized successfully")
}
