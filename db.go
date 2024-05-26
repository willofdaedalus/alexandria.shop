package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// USER DB HANDLERS
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
	var (
		storedHash string
		query      = `SELECT passwordHash FROM users WHERE username=?`
	)

	err := db.QueryRow(query, authUser.username).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			// user doesn't exist
			return fmt.Errorf("user doesn't exist.\n\npress enter to continue")
		}
		// another error occurred
		return fmt.Errorf("error querying user: %w", err)
	}

	// compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(authUser.password))
	if err != nil {
		// incorrect password
		return fmt.Errorf("invalid password for user\n\npress enter to continue")
	}

	return nil
}

// BOOK DB HANDLERS

// create the books table
func createBooksTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL UNIQUE,
        author TEXT NOT NULL,
        description TEXT NOT NULL,
        price REAL NOT NULL,
        genre TEXT NOT NULL
    );
    `
	_, err := db.Exec(query)
	return err
}

func addBook(db *sql.DB, b book) error {
	query := `
    INSERT INTO books (
        title, author, description, price, genre
    ) VALUES (
        ?, ?, ?, ?, ?
    );`

	_, err := db.Exec(query, b.Title, b.Author, b.Description, b.Price, b.Genre)
	return err
}

func getAllBooks() ([]book, error) {
	// slices get resized by golang but this is a good start
	var books []book = make([]book, 20)

	// Open the JSON file
	jsonFile, err := os.Open("books.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	// Read the file content
	byteValues, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON content into a slice of structs
	err = json.Unmarshal(byteValues, &books)
	if err != nil {
		return nil, err
	}

	return books, nil
}
