package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)


func trimTextInput(s string) string {
    nS := ""
    if strings.Contains(s, "> ") {
        nS, _ = strings.CutSuffix(s, "> ")
    }

    return nS
}


func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.authInputs))

	// Only text authInputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.authInputs {
		m.authInputs[i], cmds[i] = m.authInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}