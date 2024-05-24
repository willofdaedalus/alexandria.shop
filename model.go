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
	loginInputs   []textinput.Model
	loginCurField int
	// sign up text inputs
	signupInputs      []textinput.Model
	signupScrCurField int

	// scrTimer to transition to login
	scrTimer timer.Model

	view int

	termWidth  int
	termHeight int
}

func initialModel() model {
	m := model{
		scrTimer:     timer.NewWithInterval(startScrTimeout, time.Second),
		loginInputs:  readyInputsFor(2),
		signupInputs: readyInputsFor(3),
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.scrTimer.Init(), textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd    tea.Cmd
		field  *int
		inputs []textinput.Model
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		case "enter":
			if m.view == welcome {
				m.view = 1
			} else if m.view == login || m.view == signUp {
				// validate the inputs here
			}
		// debug purposes only
		case "0", "1", "2":
			m.view, _ = strconv.Atoi(msg.String())
		case "tab", "shift+tab", "up", "down":
			i := msg.String()

			if m.view == login {
				field = &m.loginCurField
				inputs = m.loginInputs
			} else if m.view == signUp {
				field = &m.signupScrCurField
				inputs = m.signupInputs
			}

			if i == "tab" || i == "down" {
				nextInput(field, len(inputs))
			} else {
				prevInput(field, len(inputs))
			}

			cmds := focusFields(field, inputs)

			return m, tea.Batch(cmds...)
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

	// update the current inputs' focus based on the view
	if m.view == login {
		inputs = m.loginInputs
	} else if m.view == signUp {
		inputs = m.signupInputs
	}
	cmd = m.updateInputs(msg, inputs)

	return m, cmd
}

func (m model) View() string {
	v := ""

	switch m.view {
	case 0:
		v = m.initialScreen()
	case 1:
		v = m.loginScreen()
	case 2:
		v = m.signUpScreen()
	}

	return v
}
