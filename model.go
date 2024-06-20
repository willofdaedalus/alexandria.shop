package main

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

func initialModel(db *sql.DB) model {
	// cart initialisation
	c := initialCart()

	m := model{
		scrTimer:     timer.NewWithInterval(startScrTimeout, time.Second),
		loginInputs:  readyInputsFor(2),
		signupInputs: readyInputsFor(3),
		db:           db, curPage: 1,
		c:              c,
		itemsDispCount: magicNum,
	}

	books, _ := getBooksForPage(db, 7, m.mainOffset)
	if books == nil {
		log.Fatalf("no books found")
	}

	m.curBooks = books

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.scrTimer.Init(), textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd         tea.Cmd
		err         error
		scrCtxLen   int               // this stores the len of all input fields, items and such depending on the screen
		count       int               // this stores the number of items that can be iterated on the screen
		field       *int              // this determines the field depending on what auth screen we're on
		itemTracker *int              // this keeps items that checks out of bounds errors from happening
		uName       string            // this stores the name the user entered
		uPwd        string            // this stores the password the user entered
		uRePwd      string            // this stores the password confirmation the user entered at signup
		wrap        bool              // wraps input so that the selector go back to the start if at the end
		inputs      []textinput.Model // this determines the input fields also depending on the auth screen
	)

	// setup the dimensions stuff here
	m.spatials = setupDimensions(m.termHeight, m.termWidth)
	// Calculate itemsCount before rendering
	m.itemsDispCount = m.spatials.innerH / 3

	if m.view == vCatalogue {
		count, err = countBooks(m.db)
		if err != nil {
			log.Fatal(err)
		}
	} else if m.view == vCart {
		count = len(m.c.items) - 1
	}

	// ready the item tracker to be the mainItemsIterated
	itemTracker = &m.mainItemsIterated

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
		field = &m.mainItemsIter
		wrap = false
		scrCtxLen = m.itemsDispCount
		itemTracker = &m.mainItemsIterated
	} else if m.view == vCart {
		field = &m.cartItemIter
		wrap = false
		scrCtxLen = len(m.c.items)
		itemTracker = &m.cartItemsIterated
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit

		case "c":
			if m.view != vCart && slices.Contains(mainViews, m.view) {
				transitionView(&m, vCart)
			}

		case "ctrl+s":
			if m.view == vLogin {
				transitionView(&m, vSignUp)
			}

		case "ctrl+l":
			if m.view == vSignUp {
				transitionView(&m, vLogin)
			} else if m.view == vCatalogue {
				// simple logout
				m.curPage = 1
				m.mainItemsIter = 0
				m.curBooks, _ = getBooksForPage(m.db, 1, 4) // this could be an error
				m.resetFields()
				m.c.items = make([]cartItem, 0)
				*itemTracker = 0
				transitionView(&m, vLogin)
			}

		case "?", "/":
			if !slices.Contains(startingViews, m.view) {
				transitionView(&m, vHelp)
			}

		case "esc":
			// go back from the books view to the main catalogue view
			if m.view == vCart {
				// force the transition from book view to catalogue view to fix
				// the bug that forces the helpView/bookDetail view
				transitionView(&m, vCatalogue)
			} else if m.view == vHelp {
				transitionView(&m, m.prevView)
			}

		case "+", "-", "_", "=":
			s := msg.String()
			if m.view == vCatalogue {
				if s == "+" || s == "=" {
					m.c.addToCart(selectedBook)
				} else if s == "-" || s == "_" {
					m.c.removeFromCart(selectedBook)
				}
				// reset the iterators and iterated to keep the selector in check
				m.cartItemIter = 0
				m.cartItemsIterated = 0
			} else if m.view == vCart {
				if len(m.c.items) > 0 {
					if s == "-" || s == "_" {
						// curBook := m.c.allTitles()[*itemTracker]
						curBook := m.c.cartItemsToDisp(m.cartOffset, m.spatials)[*itemTracker]
                        m.c.removeFromCartStr(curBook)

						if m.cartItemIter == 0 {
							return m, nil
						}
						// move the selector to the previous book in cart
						m.cartItemIter--
						*itemTracker = m.cartItemIter
					}
				}
			}

		case "enter":
			switch m.view {
			case vWelcome:
				// skip the "WELCOME TO ALEXANDRIA" screen when we hit enter
				transitionView(&m, vLogin)
				return m, m.scrTimer.Toggle()

				// probably shouldn't be handling auth processing here...
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

		case "tab", "shift+tab", "up", "down":
			i := msg.String()

			// check that the current view can be navigated
			if slices.Contains(validNavigationViews, m.view) {
				if i == "tab" || i == "down" {
					// we're at the end of the list
					if (m.view == vCatalogue || m.view == vCart) && *itemTracker >= count {
						return m, nil
					}
					atBot := nextInput(field, scrCtxLen, wrap)
					*itemTracker++

					// check if we're at the end of the list and if we're, simply request
					// the next set of pages needed to render
					if atBot {
						if m.view == vCatalogue {
							books, err := getBooksForPage(m.db, m.itemsDispCount, m.mainOffset)
							if err != nil || books == nil {
								return m, nil
							}

							m.curBooks = books
							m.mainOffset++
						} else if m.view == vCart {
							m.cartOffset++
						}
					}
				} else if i == "shift+tab" || i == "up" {
					if (m.view == vCatalogue || m.view == vCart) && *itemTracker == 0 {
						return m, nil
					}
					atTop := prevInput(field, scrCtxLen, wrap)
					*itemTracker--

					// check if we're at the start of the list and if we're, simply request
					// the next set of pages needed to render
					if atTop {
						if m.view == vCatalogue {
							books, err := getBooksForPage(m.db, m.itemsDispCount, m.mainOffset)
							// books, err := getBooksForPage(m.db, m.curPage-1, 3)

							// make check to determine the incoming books are the same as the rendered
							// ones before moving the selector up to the next page
							if err != nil {
								return m, nil
							}

							// assign the new books to render and then display their length
							m.curBooks = books
							m.mainOffset--
						} else if m.view == vCart {
							m.cartOffset--
						}
					}
				}
			}

			// this is the controller for the input fields at login, signup
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
	// renderWidth := m.termWidth/2
	//    v := fmt.Sprintf("w: %d h: %d\n", m.termWidth, m.termHeight)
	// v +=  "├" + strings.Repeat("─", renderWidth) + "┤"

	// v = m.catalogueScreen("daedalus")
	// v = m.mainScreen()

	switch m.view {
	case vWelcome:
		v = m.initialScreen()
	case vLogin:
		v = m.loginScreen()
	case vSignUp:
		v = m.signUpScreen()
	case vCredErr:
		v = m.infoScreen(m.authErr.Error(), 50, 3)
	case vSuccess:
		v = m.infoScreen("sign up successful!\n\npress enter to login now", 50, 3)
	case vCatalogue:
		v = m.mainScreen()
	case vHelp:
		v = m.helpScreen()
	case vCart:
		v = m.cartScreen()
	}

	//    v := fmt.Sprintf("w: %d h: %d\n", m.termWidth, m.termHeight)
	return v
}

// validates the input fields before sending to the db for authentication
func (m *model) validateCreds(creds ...string) error {
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
func (m *model) authUser(creds ...string) error {
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

func (m *model) resetFields() {
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
