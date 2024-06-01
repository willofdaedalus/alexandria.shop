package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

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

		uName  string // this stores the name the user entered
		uPwd   string // this stores the password the user entered
		uRePwd string // this stores the password confirmation the user entered at signup

		wrap      bool // wraps input
		scrCtxLen int  // this stores the len of all input fields, items and such depending on the screen
	)

	// update the current inputs' focus based on the view
	if m.view == vLogin {
		field = &m.loginCurField
		inputs = m.loginInputs
		wrap = true
		scrCtxLen = len(m.loginInputs)
	} else if m.view == vSignUp {
		field = &m.signupCurField
		inputs = m.signupInputs
		wrap = true
		scrCtxLen = len(m.signupInputs)
	} else if m.view == vCatalogue {
		field = &m.curItem
		wrap = false
		scrCtxLen = 3
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
					return m, nil
				}
				// at this point user has been authenticated so we can safely do this
				m.curUser.username = m.loginInputs[0].Value()
				transitionView(&m, vCatalogue)

			case vSuccess:
				// m.resetFields()
				transitionView(&m, vLogin)

			case vCredErr:
				// m.resetFields()
				transitionView(&m, m.prevView)
			}
		// debug purposes only
		case "0", "1", "2":
			m.view, _ = strconv.Atoi(msg.String())
		case "tab", "shift+tab", "up", "down":
			i := msg.String()

			if m.view > vWelcome {
				if i == "tab" || i == "down" {
					nextInput(field, scrCtxLen, wrap)
				} else {
					prevInput(field, scrCtxLen, wrap)
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

	// v = m.catalogueScreen("daedalus")

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
	case vCatalogue:
		v = m.catalogueScreen(m.curUser.username)
	}

	return v
}

// validates the input fields before sending to the db for authentication
func (m model) validateCreds(creds ...string) error {
	// probably a better way to structure this function
	var (
		result      string = ""
		err         error  = nil
		emptyFields bool
		shortFields bool
	)

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

	// check for credentials shorter than 5 characters
	for i := range creds {
		if len(creds[i]) < 4 {
			switch i {
			case 0:
				result += "username must be at least 4 characters \n"
			case 1:
				result += "password must be at least 4 characters long\n"
			case 2:
				result += "re-entered password must be at least 4 characters\n"
			}
			err = fmt.Errorf("%s\npress enter to continue", result)
			shortFields = true
		}
	}

	// return if there are short fields
	if shortFields {
		return err
	}

	if strings.Contains(creds[0], " ") {
		return fmt.Errorf("usernames shouldn't contain any spaces\n\npress enter to continue")
	}

	// we display the error message
	if m.view == vSignUp {
		if creds[1] != creds[2] {
			result = "your new passwords don't match"
			err = fmt.Errorf("%s\n\npress enter to continue", result)
		} else if creds[0] == creds[1] {
			// make sure the password doesn't match the username; no need to check for the other field
			// because if they don't match the password will be rejected before the code here gets run
			result = "username and password are the same. please use a different password"
			err = fmt.Errorf("%s\n\npress enter to continue", result)
		}
	}

	return err
}

// actual db authentication handler
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

func (m model) resetFields() {
	for i := range m.loginInputs {
		m.loginInputs[i].Reset()
	}

	for i := range m.signupInputs {
		m.signupInputs[i].Reset()
	}

	// code below doesn't work some reason
	// find out why!
	m.loginCurField = 0
	m.signupCurField = 0

}
