package main

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	// login/sign up text inputs
	authInputs   []textinput.Model
	authCurField int

	// scrTimer to transition to login
	scrTimer timer.Model

	view int

	termWidth  int
	termHeight int
}

func initialModel() model {
	m := model{
		scrTimer:   timer.NewWithInterval(startScrTimeout, time.Second),
		authInputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.authInputs {
		t = textinput.New()
		t.CharLimit = 16

		switch i {
		case 0:
			t.Focus()
		case 1:
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		case 2:
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		}

		m.authInputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	// Start the timer when the program initializes
	return tea.Batch(m.scrTimer.Init(), textinput.Blink)
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
		// debug purposes only
		case "0", "1", "2":
			m.view, _ = strconv.Atoi(msg.String())
		case "tab", "shift+tab", "up", "down":
			i := msg.String()

			if m.view == login {
				if i == "tab" || i == "down" {
					if m.authCurField < 1 {
						m.authCurField++
					} else {
						m.authCurField = 0
					}
				} else {
					if m.authCurField > 0 {
						m.authCurField--
					} else {
						m.authCurField = 1
					}

				}

				cmds := make([]tea.Cmd, len(m.authInputs))
				for i := 0; i <= len(m.authInputs)-1; i++ {
					if i == m.authCurField {
						// Set focused state
						cmds[i] = m.authInputs[i].Focus()
						m.authInputs[i].PromptStyle = magenta
						m.authInputs[i].TextStyle = magenta

						continue
					}
					// Remove focused state
					m.authInputs[i].Blur()
						m.authInputs[i].PromptStyle = faded
						m.authInputs[i].TextStyle = faded
				}

				return m, tea.Batch(cmds...)
			}
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

	cmd = m.updateInputs(msg)

	return m, cmd
}

func (m model) View() string {
	v := ""

	switch m.view {
	case 0:
		v = m.initialScreen()
	case 1:
		v = m.loginScreen()
	}

	return v
}
