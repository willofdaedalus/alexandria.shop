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
	var userExists string

	// check if the user already exists
	query := `SELECT username FROM users WHERE username=?`
	err := db.QueryRow(query, newUser.username).Scan(&userExists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking if user exists: %w", err)
	}
	if userExists != "" {
		return fmt.Errorf("user %q already exists", newUser.username)
	}

	// generate the password hash from the user's password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newUser.password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// insert the new user into the database
	query = `INSERT INTO users (username, passwordHash) VALUES (?, ?)`
	_, err = db.Exec(query, newUser.username, hashedPwd)
	if err != nil {
		return fmt.Errorf("error inserting new user: %w", err)
	}

	return nil
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
        longDesc TEXT NOT NULL,
        year INTEGER NOT NULL,
        genre TEXT NOT NULL,
        price REAL NOT NULL
    );
    `
	_, err := db.Exec(query)
	return err
}

func addBook(db *sql.DB, b book) error {
	query := `
    INSERT INTO books (
        title, author, description, longDesc, year, genre, price
    ) VALUES (
        ?, ?, ?, ?, ?, ?, ?
    );`

	_, err := db.Exec(query, b.Title, b.Author, b.Description, b.LongDesc, b.Year, b.Genre, b.Price)
	return err
}

func getBooksForPage(db *sql.DB, pageItems, offset int) ([]book, error) {
	var (
		books []book
		b     book
	)

	query := `SELECT title, author, description, longDesc, year, genre, price FROM books LIMIT ? OFFSET ?`
	rows, err := db.Query(query, pageItems, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&b.Title, &b.Author, &b.Description, &b.LongDesc, &b.Year, &b.Genre, &b.Price)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	return books, nil
}

func readBooksFromJson() ([]book, error) {
	// slices get resized by golang but this is a good start
	var books []book = make([]book, 20)

	// Open the JSON file
	jsonFile, err := os.Open("small_list.json")
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
