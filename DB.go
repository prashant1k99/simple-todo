package main

import (
	"database/sql"
	"os"
	"path/filepath"

	gap "github.com/muesli/go-app-paths"

	_ "github.com/mattn/go-sqlite3"
)

var dbQueries *sql.DB

func RunMigration() error {
	if dbQueries == nil {
		panic("Database connection is not initialized")
	}

	_, err := dbQueries.Exec(`
	CREATE TABLE IF NOT EXISTS todos (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT,
  is_closed BOOLEAN NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`)
	if err != nil {
		return err
	}

	return nil
}

func initDB() {
	// Initialize the database path
	scope := gap.NewScope(gap.User, "simple-todo")

	// Get the data directory for the application
	dbDir, err := scope.DataDirs()
	if err != nil {
		panic(err)
	}

	// Create the database path
	dbPath := filepath.Join(dbDir[0], "/st.db")

	// Check if the directory exists, otherwise create directory
	if _, err := os.Stat(filepath.Dir(dbPath)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			panic(err)
		}
	}

	// Open the database
	dbQueries, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	err = RunMigration()
	if err != nil {
		panic(err)
	}
}

func closeDB() {
	if dbQueries != nil {
		dbQueries.Close()
	}
}

// func testQuery() {
// 	if dbQueries == nil {
// 		panic("Database connection is not initialized")
// 	}

// 	err := dbQueries.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Database connection is alive")
// }
