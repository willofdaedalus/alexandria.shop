package main

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
)

type user struct {
	username string
	password string
}

type book struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Year        int     `json:"year"`
	Genre       string  `json:"genre"`
}

type model struct {
	// LOGIN/SIGN UP VALUES
	// login text inputs
	loginInputs   []textinput.Model
	loginCurField int
	// sign up text inputs
	signupInputs   []textinput.Model
	signupCurField int
	authErr        error

	// scrTimer to transition to login
	scrTimer timer.Model

	// used to transition between different views
	view     int
	prevView int

	// get the current terminal's width and height
	termWidth  int
	termHeight int

	itemsOnDisplay []string
	curItem        int

	// current session user
	curUser user
	// db handler
	db *sql.DB
}
