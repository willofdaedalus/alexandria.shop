package main

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	startScrTimeout = time.Second * 5
)

// views tracker
const (
	vWelcome = iota
	vLogin
	vSignUp
	vCredErr
	vSuccess
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
	username   string
	password   string
	rePassword string
)
