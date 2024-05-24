package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// the initial screen to display when the program first runs
func (m model) initialScreen() string {
	// Center the ASCII art within the terminal window
	return lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		cyan.Align(lipgloss.Center).Render(welcomeAscii),
	)
}

func (m model) authScreen() string {
	var doc strings.Builder

	return doc.String()
}

func (m model) loginScreen() string {
	var (
		layout     strings.Builder
		login, pwd string
	)

	helpText := "press ctrl+s to sign up | press ctrl+c to quit"
    helpBox := noBorderStyle.PaddingTop(2).Width(50).Align(lipgloss.Bottom).Render(helpText)

	loginTextt := noBorderStyle.Render("username")
	loginBox := textBoxStyle.
        Foreground(cyan.GetForeground()).
        Render(trimTextInput(m.authInputs[0].View()))
	loginRender := lipgloss.JoinHorizontal(lipgloss.Left, loginTextt, loginBox)

	pwdTextt := noBorderStyle.Render("password")
	pwdBox := textBoxStyle.Render(trimTextInput(m.authInputs[1].View()))
	pwdRender := lipgloss.JoinHorizontal(lipgloss.Left, pwdTextt, pwdBox)

	if m.authCurField == 0 {
		login = magenta.PaddingLeft(8).Align(lipgloss.Left).Bold(true).Render(loginRender)
		pwd = faded.PaddingLeft(8).Align(lipgloss.Left).Bold(true).Render(pwdRender)
	} else {
		login = faded.PaddingLeft(8).Align(lipgloss.Left).Bold(true).Render(loginRender)
		pwd = magenta.PaddingLeft(8).Align(lipgloss.Left).Bold(true).Render(pwdRender)
	}

	prompt := lipgloss.NewStyle().
		Foreground(cyan.GetForeground()).
		Width(55).Height(5).
		Align(lipgloss.Center).
		Render(loginText)

	textFields := lipgloss.JoinVertical(lipgloss.Left, login, pwd, helpBox)
	ui := lipgloss.JoinVertical(lipgloss.Left, prompt, textFields)

	dialog := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
	)

	layout.WriteString(dialog)

	return layout.String()
}
