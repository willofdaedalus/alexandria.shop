package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// the initial screen to display when the program first runs
func (m model) initialScreen() string {
	ascii := `
              _                            _        
             | |                          | |       
__      _____| | ___ ___  _ __ ___   ___  | |_ ___  
\ \ /\ / / _ \ |/ __/ _ \| '_ ' _ \ / _ \ | __/ _ \ 
 \ V  V /  __/ | (_| (_) | | | | | |  __/ | || (_) |
  \_/\_/ \___|_|\___\___/|_| |_| |_|\___|_ \__\___/ 
      | |                        | |    (_)         
  __ _| | _____  ____ _ _ __   __| |_ __ _  __ _    
 / _' | |/ _ \ \/ / _' | '_ \ / _' | '__| |/ _' |   
| (_| | |  __/>  < (_| | | | | (_| | |  | | (_| |   
 \__,_|_|\___/_/\_\__,_|_| |_|\__,_|_|  |_|\__,_|   
                                                    
                                                    
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

	okButton := activeButtonStyle.Render("Yes")
	cancelButton := buttonStyle.Render("Maybe")

	question := lipgloss.NewStyle().
        Width(100).Height(m.termHeight / 2).
        Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	dialog := lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("ox"),
		lipgloss.WithWhitespaceForeground(subtle),
	)

    doc.WriteString(dialog)

	return doc.String()
}
