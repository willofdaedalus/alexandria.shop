package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// the initial screen to display when the program first runs
func (m model) initialScreen() string {
	ascii := `
██     ██ ███████ ██       ██████  ██████  ███    ███ ███████     ████████  ██████  
██     ██ ██      ██      ██      ██    ██ ████  ████ ██             ██    ██    ██ 
██  █  ██ █████   ██      ██      ██    ██ ██ ████ ██ █████          ██    ██    ██ 
██ ███ ██ ██      ██      ██      ██    ██ ██  ██  ██ ██             ██    ██    ██ 
 ███ ███  ███████ ███████  ██████  ██████  ██      ██ ███████        ██     ██████  
                                                                                    
                                                                                    
 █████  ██      ███████ ██   ██  █████  ███    ██ ██████  ██████  ██  █████         
██   ██ ██      ██       ██ ██  ██   ██ ████   ██ ██   ██ ██   ██ ██ ██   ██        
███████ ██      █████     ███   ███████ ██ ██  ██ ██   ██ ██████  ██ ███████        
██   ██ ██      ██       ██ ██  ██   ██ ██  ██ ██ ██   ██ ██   ██ ██ ██   ██        
██   ██ ███████ ███████ ██   ██ ██   ██ ██   ████ ██████  ██   ██ ██ ██   ██        
`
	// Center the ASCII art within the terminal window
	return lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		cyan.Align(lipgloss.Center).Render(ascii),
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

	loginText := `
██       ██████   ██████  ██ ███    ██ 
██      ██    ██ ██       ██ ████   ██ 
██      ██    ██ ██   ███ ██ ██ ██  ██ 
██      ██    ██ ██    ██ ██ ██  ██ ██ 
███████  ██████   ██████  ██ ██   ████ 


`

    usrText := fmt.Sprintf("username \n%s\n", m.authInputs[0].View())
    pwdText := fmt.Sprintf("password \n%s", m.authInputs[1].View())


	if m.authCurField == 0 {
		login = magenta.Copy().PaddingLeft(10).Align(lipgloss.Left).Bold(true).Render(usrText)
		pwd = faded.Copy().PaddingLeft(10).Align(lipgloss.Left).Bold(true).Render(pwdText)
	} else {
		login = faded.Copy().PaddingLeft(10).Align(lipgloss.Left).Bold(true).Render(usrText)
		pwd = magenta.Copy().PaddingLeft(10).Align(lipgloss.Left).Bold(true).Render(pwdText)
	}

	prompt := lipgloss.NewStyle().
		Foreground(cyan.GetForeground()).
		Width(45).Height(5).
		Align(lipgloss.Center).
		Render(loginText)

	textFields := lipgloss.JoinVertical(lipgloss.Left, login, pwd)
	ui := lipgloss.JoinVertical(lipgloss.Left, prompt, textFields)

	dialog := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
	)

	layout.WriteString(dialog)

	return layout.String()
}
