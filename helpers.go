package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// go back to the next field
func nextInput(field *int, size int, wrap bool) bool {
	if *field < size-1 {
		*field++
	} else {
		if wrap {
			*field = 0
		}
		return true
	}
	return false
}

// go back to the previous field
func prevInput(field *int, size int, wrap bool) bool {
	if *field > 0 {
		*field--
	} else {
		if wrap {
			*field = size - 1
		}
		return true
	}
	return false
}

func focusFields(field *int, inputs []textinput.Model) []tea.Cmd {
	cmds := make([]tea.Cmd, len(inputs))
	for i := 0; i <= len(inputs)-1; i++ {
		if i == *field {
			// Set focused state
			cmds[i] = inputs[i].Focus()
			inputs[i].PromptStyle = magenta
			inputs[i].TextStyle = magenta

			continue
		}
		// Remove focused state
		inputs[i].Blur()
		inputs[i].PromptStyle = faded
		inputs[i].TextStyle = faded
	}

	return cmds
}

func (m *model) updateInputs(msg tea.Msg, inputs []textinput.Model) tea.Cmd {
	cmds := make([]tea.Cmd, len(inputs))

	// Only text authInputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range inputs {
		inputs[i], cmds[i] = inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func readyInputsFor(size int) []textinput.Model {
	ret := make([]textinput.Model, size)

	var t textinput.Model
	for i := 0; i < size; i++ {
		t = textinput.New()
		t.CharLimit = 16
		t.Prompt = "" // added prompt here but renders no input

		switch i {
		case 0:
			t.Focus()
		case 1:
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		case 2: // doesn't matter if login won't reach here; it's a switch-case statement after all
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		}

		ret[i] = t
	}

	return ret
}

func transitionView(m *model, to int) {
	prev := m.view
	m.view = to
	m.prevView = prev
}
