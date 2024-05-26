package main

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
    _ "github.com/mattn/go-sqlite3"
)

func createUsersTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(16) NOT NULL UNIQUE,
        passwordHash TEXT NOT NULL
    );`

	_, err := db.Exec(query)
	return err
}

// add a new user to the database
func signupUser(db *sql.DB, newUser user) error {
	// generate the password hash from the user's password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newUser.password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO users(username, passwordHash) VALUES (?, ?);`
	_, err = db.Exec(query, newUser.username, hashedPwd)

	return err
}

// login handler
func loginUser(db *sql.DB, authUser user) error {
	// maybe this is overkill but we're really trying to preventing sql injection attacks
	// see more here: https://go.dev/doc/database/sql-injection
	query := `SELECT passwordHash FROM users WHERE username=?`

	var storedHash string
    err := db.QueryRow(query, authUser.username).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesn't exist. consider signing up, it's free")
		}
		return fmt.Errorf("error querying user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(authUser.password))
	if err != nil {
		return fmt.Errorf("invalid password or username")
	}

	return err
}
