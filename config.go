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

	textBoxStyle = lipgloss.NewStyle().
			Padding(0, 3).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(cyan.GetForeground())

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
		// BorderForeground(magenta.GetForeground()).
		BorderBackground(magenta.GetForeground()).
		Padding(1, 0)

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
