package main

import (
	"database/sql"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./alexandria.db")
	if err != nil {
		log.Fatalf("open err: %s", err.Error())
	}
	defer db.Close()

	// check if the users table exists before creating it
	if !tableExists(db, "users") {
		err = createUsersTable(db)
		if err != nil {
			log.Fatalf("user table create err: %s", err.Error())
		}
	}

	// check if the books table exists before creating it
	if !tableExists(db, "books") {
		err = createBooksTable(db)
		if err != nil {
			log.Fatalf("book table create err: %s", err.Error())
		}

		// populate books table only if it's newly created
		books, err := readBooksFromJson()
		if err != nil {
			log.Fatal(err)
		}

		for _, book := range books {
			err = addBook(db, book)
			if err != nil {
				log.Fatalf("book insert err: %s", err.Error())
			}
		}
	}

	m := initialModel(db)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		log.Fatalf("tea program err: %s", err.Error())
	}
}


// function to check if a table exists in the database
func tableExists(db *sql.DB, tableName string) bool {
	// https://www.quora.com/How-can-you-check-if-a-table-exists-in-SQLite
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
	rows, err := db.Query(query, tableName)
	if err != nil {
		log.Fatalf("error checking if table exists: %s", err.Error())
	}
	defer rows.Close()

	return rows.Next()
}
