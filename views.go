package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// the initial screen to display when the program first runs
func (m *model) initialScreen() string {
	// center the ascii art within the terminal window
	return lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		cyan.Align(lipgloss.Center).Render(welcomeAscii),
	)
}

func (m *model) renderAuthScreen(title, helpText string, inputs []textinput.Model, curField int) string {
	var layout strings.Builder

	// render the title at the top
	titleRender := lipgloss.NewStyle().
		Foreground(cyan.GetForeground()).
		Width(55).Height(5).
		Align(lipgloss.Center).
		Render(title)

	// render input boxes
	var inputRenders []string
	for i, label := range []string{"username", "password", "password"} {
		// check for the iterator not to overlap in the case of login screen
		if i >= len(inputs) {
			break
		}
		inputRenders = append(inputRenders, m.renderInputBox(label, i, charLimitAuth+5, inputs, curField == i))
	}

	// render the help text at the bottom
	helpBox := noBorderStyle.PaddingTop(1).Width(50).Align(lipgloss.Bottom).Render(helpText)

	// join the input fields and help text
	textFields := lipgloss.JoinVertical(lipgloss.Left, inputRenders...)
	textFields = lipgloss.JoinVertical(lipgloss.Left, textFields, helpBox)

	// combine the title and input fields
	ui := lipgloss.JoinVertical(lipgloss.Left, titleRender, textFields)

	// place the ui in the center of the screen
	dialog := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
	)

	// return the final rendered layout
	layout.WriteString(dialog)
	return layout.String()
}

func (m *model) signUpScreen() string {
	return m.renderAuthScreen(
		signUpText,
		"press ctrl+l to log in | press ctrl+c to quit",
		m.signupInputs,
		m.signupCurField,
	)
}

func (m *model) loginScreen() string {
	return m.renderAuthScreen(
		loginText,
		"press ctrl+s to sign up | press ctrl+c to quit",
		m.loginInputs,
		m.loginCurField,
	)
}

func (m *model) infoScreen(info string, w, h int) string {
	infoRender := noBorderStyle.
		Padding(0, 2, 0, 2).
		Width(w).Height(h).
		// Width(50).Height(3).
		Align(lipgloss.Left).Render(info)

	// footer/help message render
	// helpText := "press enter"
	// helpBox := noBorderStyle.Width(50).Height(1).Align(lipgloss.Bottom).Render(helpText)

	dialog := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Width(50).Render(infoRender),
	)

	return dialog
}

func (m *model) helpScreen(msg string, w, h int, pos lipgloss.Position) string {
	infoRender := noBorderStyle.
		PaddingTop(1).
		Width(w).Height(h).
		Padding(0, 2, 0, 2).
		Align(pos).Render(msg)
	// Width(80).Height(10).
	// Padding(0, 2, 0, 2).
	// Align(lipgloss.Left).Render(helpText)

	// footer/help message render
	// helpText := "press enter"
	// helpBox := noBorderStyle.Width(50).Height(1).Align(lipgloss.Bottom).Render(helpText)

	dialog := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Width(80).Render(infoRender),
	)

	return dialog
}
