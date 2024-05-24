package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// go back to the next field
func nextInput(field *int, size int) {
	if *field < size - 1 {
		*field++
	} else {
		*field = 0
	}
}

// go back to the previous field
func prevInput(field *int, size int) {
    if *field > 0 {
        *field--
    } else {
        *field = size - 1
    }
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

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.loginInputs))

	// Only text authInputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.loginInputs {
		m.loginInputs[i], cmds[i] = m.loginInputs[i].Update(msg)
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
		case 2:
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		}

		ret[i] = t
	}

	return ret
}
