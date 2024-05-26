package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Remove("./users.db")

	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("open err: %s", err.Error())
	}
	defer db.Close()

	err = createUsersTable(db)
	if err != nil {
		log.Fatalf("table create err: %s", err.Error())
	}

	m := initialModel(db)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
