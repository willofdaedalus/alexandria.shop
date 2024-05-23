package main

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	// login/sign up text inputs
	authInputs []textinput.Model

	// scrTimer to transition to login
	scrTimer timer.Model

	view int

	termWidth  int
	termHeight int
}

func initialModel() model {
	m := model{
		scrTimer: timer.NewWithInterval(startScrTimeout, time.Second),
	}

	return m
}

func (m model) Init() tea.Cmd {
	// Start the timer when the program initializes
	return m.scrTimer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		case "enter":
			m.view = 1
		}

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

	case timer.TickMsg, timer.StartStopMsg:
		m.scrTimer, cmd = m.scrTimer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.view = 1
	}

	return m, nil
}

func (m model) View() string {
	v := ""

	switch m.view {
	case 0:
		v = m.initialScreen()
	case 1:
		v = m.authScreen()
	}

	return v
}
