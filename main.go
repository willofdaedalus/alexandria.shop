package main

import (
	"fmt"
	"os"
    "database/sql"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "./users.db")
    if err != nil {
        fmt.Println("there was an err: ", err)
        os.Exit(1)
    }
    defer db.Close()

    fmt.Println(db)

	m := initialModel(db)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
