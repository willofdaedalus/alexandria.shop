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
	LongDesc    string  `json:"long_description"`
	Year        int     `json:"year"`
	Genre       string  `json:"genre"`
	Price       float64 `json:"price"`
}

type cartItem struct {
	title string
	price float64
}

type cart struct {
	items []cartItem
}

type model struct {
	// LOGIN/SIGN UP VALUES
	authErr           error
	loginInputs       []textinput.Model
	signupInputs      []textinput.Model
	checkoutInputs    []textinput.Model
	loginCurField     int
	signupCurField    int
	checkoutCurField  int
	view              int
	prevView          int
	curPage           int
	termWidth         int
	termHeight        int
	itemsDispCount    int
	mainItemsIterated int
	mainItemsIter     int
	cartItemsIterated int
	cartItemIter      int
	mainOffset        int
	cartOffset        int
	scrTimer          timer.Model
	curBooks          []book            // books that are being displayed
	curCartItems      []string          // all items that can be displayed from the cart
	curUser           user              // information about the current user
	db                *sql.DB           //db handler
	c                 cart              // cart system
	content           mainRenderContent // things to pass to the main frame to render
	spatials          dimensions        // essential variables needed to ease responsiveness
}

type dimensions struct {
	dynRenderWidth  int // calculates the best width
	dynRenderHeight int // calculates the best height
	actualRenderW   int // that's the width of the main border
	actualRenderH   int // that's the height of the main border
	innerW          int // padding subtracted from the main width
	innerH          int // padding subtracted from the main height
	listSectionW    int // list section width
	bookDeetsW      int // book details width
}

type mainRenderContent struct {
	headerContents []string
	bookItems      []book
	bookDetails    string
	footerMessage  string
}
