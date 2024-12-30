package config

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	// Ensure data directory exists
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	dbPath := filepath.Join(dataDir, "isxportfolio.db")
	log.Printf("Initializing database at: %s", dbPath)

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Database connection successful")

	// Create users table
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}
	log.Println("Users table created/verified successfully")

	// Verify table exists
	var tableName string
	err = DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users';").Scan(&tableName)
	if err != nil {
		log.Fatal("Failed to verify users table exists:", err)
	}
	log.Printf("Verified users table exists: %s", tableName)

	log.Println("Database initialized successfully")
}

func TestDatabaseWrite() {
	_, err := DB.Exec(`
		INSERT INTO users (email, name) 
		VALUES ('test@example.com', 'Test User')
		ON CONFLICT(email) DO NOTHING
	`)
	if err != nil {
		log.Printf("Test write failed: %v", err)
	} else {
		log.Println("Test write successful")
	}
}
