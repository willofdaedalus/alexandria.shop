package main

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	magicNum          = 5 // the num of items to display, to query; imma be real with you I can't change this
	startScrTimeout   = time.Second * 5
	catalogueHelpMsg  = "ctrl+c to exit  |  ctrl+l to logout  |  tab/shift+tab or arrow keys to move"
	addToCartMsg      = "+ to add book to cart"
	removeFromCartMsg = "- to remove book from cart"
	helpText          = `alexandria.shop Copyright (c) 1996

Did this ever exist?
If you want it to then yes, yes it did

Who made this?
Built by @daedalus_here on Twitter

How to use?
Use tab/down key for downward selection and shift+tab/up for upward selection, escape to exit dialogue like this
`
)

// views tracker
const (
	vWelcome = iota
	vLogin
	vSignUp
	vCredErr
	vSuccess
	vCatalogue
	vBookDetails
	vHelp
	vCart
)

const (
	welcomeAscii = `
██     ██ ███████ ██       ██████  ██████  ███    ███ ███████     ████████  ██████  
██     ██ ██      ██      ██      ██    ██ ████  ████ ██             ██    ██    ██ 
██  █  ██ █████   ██      ██      ██    ██ ██ ████ ██ █████          ██    ██    ██ 
██ ███ ██ ██      ██      ██      ██    ██ ██  ██  ██ ██             ██    ██    ██ 
 ███ ███  ███████ ███████  ██████  ██████  ██      ██ ███████        ██     ██████  
                                                                                    
                                                                                    
 █████  ██      ███████ ██   ██  █████  ███    ██ ██████  ██████  ██  █████         
██   ██ ██      ██       ██ ██  ██   ██ ████   ██ ██   ██ ██   ██ ██ ██   ██        
███████ ██      █████     ███   ███████ ██ ██  ██ ██   ██ ██████  ██ ███████        
██   ██ ██      ██       ██ ██  ██   ██ ██  ██ ██ ██   ██ ██   ██ ██ ██   ██        
██   ██ ███████ ███████ ██   ██ ██   ██ ██   ████ ██████  ██   ██ ██ ██   ██        
`

	loginText = `
██       ██████   ██████  ██ ███    ██ 
██      ██    ██ ██       ██ ████   ██ 
██      ██    ██ ██   ███ ██ ██ ██  ██ 
██      ██    ██ ██    ██ ██ ██  ██ ██ 
███████  ██████   ██████  ██ ██   ████ 

`

	signUpText = `
███████ ██  ██████  ███    ██ ██    ██ ██████  
██      ██ ██       ████   ██ ██    ██ ██   ██ 
███████ ██ ██   ███ ██ ██  ██ ██    ██ ██████  
     ██ ██ ██    ██ ██  ██ ██ ██    ██ ██      
███████ ██  ██████  ██   ████  ██████  ██      
                                               
`
)

var (
	// CMYK VALUES
	cyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffff"))
	magenta = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff00ff"))
	yellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff00"))
	white  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))

	// style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Align(lipgloss.Center, lipgloss.Center)
	faded = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))

	noBorderStyle = lipgloss.NewStyle().
			Width(10).
			Border(lipgloss.HiddenBorder())

	textBoxStyle = lipgloss.NewStyle().
			Width(20).
			Border(lipgloss.NormalBorder())

	headerBoxStyle = lipgloss.NewStyle().
			Width(20).Height(1).
			Border(lipgloss.NormalBorder())

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
		// Padding(1, 0).
		BorderForeground(magenta.GetForeground())
)

var (
	selectedBook         book // this is updated in the function that handles the selection of the book items
	catalogueViewHeight  int  // hack to keep the height the same when displaying the details of a book
	validNavigationViews = []int{vLogin, vSignUp, vCatalogue, vCart}
	startingViews        = []int{vLogin, vSignUp, vWelcome, vHelp}
	mainViews            = []int{vCatalogue, vBookDetails}
)
