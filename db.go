package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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
	// check if the user exists
	var storedHash string
	query := `SELECT passwordHash FROM users WHERE username=?`
	err := db.QueryRow(query, authUser.username).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			// user doesn't exist
			return fmt.Errorf("user doesn't exist. signup it's free")
		}
		// another error occurred
		return fmt.Errorf("error querying user: %w", err)
	}

	// compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(authUser.password))
	if err != nil {
		// incorrect password
		return fmt.Errorf("invalid password")
	}

	return nil
}
