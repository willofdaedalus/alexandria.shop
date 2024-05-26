package main

import "database/sql"

func createUsersTable(db *sql.DB) {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(16) NOT NULL UNIQUE,
        password VARCHAR(16) NOT NULL
    );`

    _, err := db.Exec(query)

}
