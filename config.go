package main

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	// CMYK VALUES
	cyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffff"))
	magenta = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff00ff"))
	yellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff00"))

	// style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Align(lipgloss.Center, lipgloss.Center)
	faded = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))

	topHeaderStyle = lipgloss.NewStyle().
			Width(50).Height(5)

	noBorderStyle = lipgloss.NewStyle().
			Width(10).
			Border(lipgloss.HiddenBorder())

	textBoxStyle = lipgloss.NewStyle().
			Width(20).
			Border(lipgloss.NormalBorder())

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 0).
			BorderForeground(magenta.GetForeground())
		// BorderBackground(magenta.GetForeground())

	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

const (
	startScrTimeout = time.Second * 5
	width           = 96
)

const (
	welcome = iota
	login
	signUp
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
