package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var defaultTimeout = time.Second * 5

var (
    cyan = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
    magenta = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff00ff"))
    yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
    style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Align(lipgloss.Center, lipgloss.Center)
)

type model struct {
	timer  timer.Model
	view   int
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	// Start the timer when the program initializes
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case timer.TickMsg, timer.StartStopMsg:
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

    case timer.TimeoutMsg:
        m.view = 1
	}

	return m, nil
}

func initialScreen(m model) string {
	ascii := `
              _                            _          _   _                          
             | |                          | |        | | | |                         
__      _____| | ___ ___  _ __ ___   ___  | |_ ___   | |_| |__   ___                 
\ \ /\ / / _ \ |/ __/ _ \| '_ ' _ \ / _ \ | __/ _ \  | __| '_ \ / _ \                
 \ V  V /  __/ | (_| (_) | | | | | |  __/ | || (_) | | |_| | | |  __/                
  \_/\_/ \___|_|\___\___/|_| |_| |_|\___|_ \__\___/ __\__|_| |_|\___|                
      | |                        | |    (_)        / _|                              
  __ _| | _____  ____ _ _ __   __| |_ __ _  __ _  | |_ ___  _ __ _   _ _ __ ___  ___ 
 / _' | |/ _ \ \/ / _' | '_ \ / _' | '__| |/ _' | |  _/ _ \| '__| | | | '_ ' _ \/ __|
| (_| | |  __/>  < (_| | | | | (_| | |  | | (_| | | || (_) | |  | |_| | | | | | \__ \
 \__,_|_|\___/_/\_\__,_|_| |_|\__,_|_|  |_|\__,_| |_| \___/|_|   \__,_|_| |_| |_|___/
                                                                                     
                                                                                     
`

	ascii += "\nwelcome fellow scholar\n"

	// Center the ASCII art within the terminal window
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, cyan.Render(ascii))
}

func (m model) View() string {
	v := ""

	switch m.view {
	case 0:
		v = initialScreen(m)
	case 1:
		v = "yet to be implemented"
	}

	return v
}

func main() {
	m := model{
		timer: timer.NewWithInterval(defaultTimeout, time.Second),
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
