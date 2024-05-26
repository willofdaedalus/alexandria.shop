package main

import (
	"database/sql"
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

	db *sql.DB
}

func initialModel(db *sql.DB) model {
	m := model{
		scrTimer:     timer.NewWithInterval(startScrTimeout, time.Second),
		loginInputs:  readyInputsFor(2),
		signupInputs: readyInputsFor(3),
		db:           db,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.scrTimer.Init(), textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd    tea.Cmd
		field  *int              // this determines the field depending on what auth screen we're on
		inputs []textinput.Model // this determines the input fields also depending on the auth screen

		uName  string // this stores the user name they entered
		uPwd   string // this stores the password the user entered
		uRePwd string // this stores the password confirmation the user entered at signup
	)

	// update the current inputs' focus based on the view
	if m.view == vLogin {
		field = &m.loginCurField
		inputs = m.loginInputs
	} else if m.view == vSignUp {
		field = &m.signupCurField
		inputs = m.signupInputs
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit

		case "ctrl+s":
			if m.view == vLogin {
				transitionView(&m, vSignUp)
			}
		case "ctrl+l":
			if m.view == vSignUp {
				transitionView(&m, vLogin)
			}

		case "enter":
			switch m.view {
			case vWelcome:
				transitionView(&m, vLogin)
				return m, m.scrTimer.Toggle()

			case vLogin, vSignUp:
				// this section could use a refactor if possible
				uName = inputs[0].Value()
				uPwd = inputs[1].Value()

				// pass the user's inputs for validation checks
				if m.view == vLogin {
					m.authErr = m.validateCreds(uName, uPwd)
				} else {
					uRePwd = inputs[2].Value() // we only need rePwd to check against pwd at signup
					m.authErr = m.validateCreds(uName, uPwd, uRePwd)
				}
				// transition to the authentication view on err
				if m.authErr != nil {
					transitionView(&m, vCredErr)
					return m, nil
				}

				// actual db authentication happens below
				m.authErr = m.authUser(uName, uPwd)
				// transition to the authentication view on err
				if m.authErr != nil {
					transitionView(&m, vCredErr)
					return m, nil
				}
				// only show the success screen when at the signup screen
				if m.view == vSignUp {
					transitionView(&m, vSuccess)
				}

			case vSuccess:
				transitionView(&m, vLogin)

			case vCredErr:
				transitionView(&m, m.prevView)
			}
		// debug purposes only
		case "0", "1", "2":
			m.view, _ = strconv.Atoi(msg.String())
		case "tab", "shift+tab", "up", "down":
			i := msg.String()

			if m.view == vLogin || m.view == vSignUp {
				if i == "tab" || i == "down" {
					nextInput(field, len(inputs))
				} else {
					prevInput(field, len(inputs))
				}
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
		transitionView(&m, vLogin)
	}

	cmd = m.updateInputs(msg, inputs)

	return m, cmd
}

func (m model) View() string {
	v := ""

	switch m.view {
	case vWelcome:
		v = m.initialScreen()
	case vLogin:
		v = m.loginScreen()
	case vSignUp:
		v = m.signUpScreen()
	case vCredErr:
		v = m.infoScreen(m.authErr.Error())
	case vSuccess:
		v = m.infoScreen("sign up successful!\n\npress enter to login now")
	}

	return v
}

// validates the input fields before sending to the db for authentication
func (m model) validateCreds(creds ...string) error {
	result := ""
	var err error = nil
	var emptyFields bool

	// check for empty fields
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
			err = fmt.Errorf("%s\npress enter to continue", result)
			emptyFields = true
		}
	}

	if emptyFields {
		return err
	}

	// we display the error message
	if m.view == vSignUp {
		if creds[1] != creds[2] {
			result = "your new passwords don't match"
			err = fmt.Errorf("%s\n\npress enter to continue", result)
		}
	}

	return err
}

func (m model) authUser(creds ...string) error {
	u := user{
		username: creds[0],
		password: creds[1],
	}

	if m.view == vLogin {
		return loginUser(m.db, u)
	} else if m.view == vSignUp {
		// since we've checked both new password and retype password fields are valid
		// there's no worry of sending the wrong information to the signup process
		return signupUser(m.db, u)
	}

	return nil
}
