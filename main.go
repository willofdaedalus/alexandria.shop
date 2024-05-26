package main

import (
	"database/sql"
	"log"
	"os"

	// tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Remove("./alexandria.db")

	db, err := sql.Open("sqlite3", "./alexandria.db")
	if err != nil {
		log.Fatalf("open err: %s", err.Error())
	}
	defer db.Close()

	err = createUsersTable(db)
	if err != nil {
		log.Fatalf("user table create err: %s", err.Error())
	}

	err = createBooksTable(db)
	if err != nil {
		log.Fatalf("book table create err: %s", err.Error())
	}

	books, err := getAllBooks()
	if err != nil {
		log.Fatal(err)
	}

	for _, book := range books {
		err = addBook(db, book)
		if err != nil {
			log.Fatalf("book insert err: %s", err.Error())
		}
	}

	// m := initialModel(db)
	//
	// if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}
