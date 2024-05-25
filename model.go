package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

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

	// update the current inputs' focus based on the view
	if m.view == login {
		field = &m.loginCurField
		inputs = m.loginInputs
	} else if m.view == signUp {
		field = &m.signupCurField
		inputs = m.signupInputs
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		case "ctrl+s":
			if m.view == login {
				transitionView(&m, signUp)
			}
		case "ctrl+l":
			if m.view == signUp {
				transitionView(&m, login)
			}
		case "enter":
			switch m.view {
			case welcome:
				transitionView(&m, login)
				return m, m.scrTimer.Toggle()
			case login:
				m.authErr = m.checkUserCreds(inputs[0].Value(), inputs[1].Value())
			case signUp:
				m.authErr = m.checkUserCreds(inputs[0].Value(), inputs[1].Value(), inputs[2].Value())
			case credErr:
				transitionView(&m, m.prevView)
			}
		// debug purposes only
		case "0", "1", "2":
			m.view, _ = strconv.Atoi(msg.String())
		case "tab", "shift+tab", "up", "down":
			i := msg.String()

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
		// transition to the login screen on timer up
		transitionView(&m, login)
	}

	cmd = m.updateInputs(msg, inputs)

	return m, cmd
}

func (m model) View() string {
	v := ""

	if m.authErr != nil {
		transitionView(&m, credErr)
	}

	switch m.view {
	case welcome:
		v = m.initialScreen()
	case login:
		v = m.loginScreen()
	case signUp:
		v = m.signUpScreen()
	case credErr:
		v = m.errorScreen()
	}

	return v
}

func (m model) checkUserCreds(creds ...string) error {
	result := ""
	var err error = nil

	for i := range creds {
		if creds[i] == "" {
			switch i {
			case 0:
				result += "please enter a username\n"
			case 1:
				result += "please enter a password\n"
			case 2:
				result += "please re-enter your new password\n"
			}
			err = fmt.Errorf(result)
		}
	}

	return err
}
