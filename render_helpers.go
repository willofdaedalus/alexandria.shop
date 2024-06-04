package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// renders a text in a particular way with colours and styling
func styleTextWith(t string, col lipgloss.TerminalColor, bold bool) string {
	return lipgloss.NewStyle().
		Bold(bold).
		Foreground(col).
		Render(t)
}

func (m model) renderInputBox(label string, index int, inputs []textinput.Model, focused bool) string {
	render := renderBoxDesc(label, index, focused, inputs)
	if focused {
		return magenta.PaddingLeft(8).Align(lipgloss.Left).Render(render)
	}
	return faded.PaddingLeft(8).Align(lipgloss.Left).Render(render)
}

// renders an item in the catalogue view
func (m model) renderItemDisplay(renderWidth, renderHeight int, selected bool, b book) string {
    var bookInCart = "[in cart]"

	selectProperties := lipgloss.NewStyle().
		Foreground(faded.GetForeground()).
		Border(lipgloss.NormalBorder())
		// Border(lipgloss.HiddenBorder())

	if selected {
		selectProperties = lipgloss.NewStyle().
			Foreground(magenta.GetForeground()).
			Border(lipgloss.ThickBorder())

		// assign the current selected book to the global selected book
		selectedBook = b
	}

    if !m.c.bookInCart(b) {
        bookInCart = " "
    }

	itemContent := lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder(), false, true, true, true).
		Foreground(selectProperties.GetForeground()).
        Render(fmt.Sprintf("%s  |  %s  |  $%.2f  \t%s", b.Title, b.Author, b.Price, bookInCart))

	itemDesc := lipgloss.NewStyle().
		Inline(true).
		Foreground(selectProperties.GetForeground()).
		Border(lipgloss.HiddenBorder(), true, true, false, true).
		Render(b.Description)

	contentRender := lipgloss.JoinVertical(lipgloss.Top, itemContent, itemDesc)

	return lipgloss.NewStyle().
		Border(selectProperties.GetBorder()).
		BorderForeground(selectProperties.GetForeground()).
		Padding(0, 1, 0, 1).
		Width(renderWidth - 5).Height(renderHeight + 3).
		Render(contentRender)
}

// render the headers at the top of the catalogue page
func (m model) renderHeaders(header1, cart string, renderWidth int) string {
	tops := [][]string{
		{
			header1,
			styleTextWith(fmt.Sprintf("welcome, %s", m.curUser.username), cyan.GetForeground(), true),
			styleTextWith(fmt.Sprint("? for help and details"), cyan.GetForeground(), true),
			styleTextWith(cart, cyan.GetForeground(), true),
		}, // actual headers
	}

	headerTable := table.New().
		Border(lipgloss.NormalBorder()).
		Width(renderWidth + 2).
		StyleFunc(table.StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)
		})).
		Rows(tops...)

	return headerTable.Render()
}

func renderHeaderBox(s string) string {
	return headerBoxStyle.
		Width(20).Margin(-1).
		Align(lipgloss.Center).
		Render(s)
}

// function to return a nicely formatted description and input box
func renderBoxDesc(s string, idx int, focused bool, inputs []textinput.Model) string {
	desc := noBorderStyle.Bold(focused).Render(s)
	// this is the content from the input box as we type
	// side not find a way to render the textbox thicker
	inputBox := textBoxStyle.Render(inputs[idx].View())
	finalRender := lipgloss.JoinHorizontal(lipgloss.Left, desc, inputBox)

	return finalRender
}

func (m model) mainScreenFrame(header1, footerMsg, content string) string {
	// Initialize variables
	renderWidth := (m.termWidth / 2) + 10
	if renderWidth < 0 {
		renderWidth = 0
	}

	headerRender := m.renderHeaders(
		styleTextWith(header1, magenta.GetForeground(), true),
		fmt.Sprintf("c cart [%d] $%.2f", len(m.c.items), m.c.booksTotal()),
		renderWidth,
	)

	footer := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		Width(renderWidth).
		Align(lipgloss.Center).
		Render(footerMsg)

	mainHeight := m.termHeight - lipgloss.Height(headerRender) - lipgloss.Height(footer)
	catalogueViewHeight = mainHeight
	// offset := (mainHeight - 20) / 3

	seperator := "├" + strings.Repeat("─", renderWidth) + "┤"

	headerRenderModified := strings.ReplaceAll(headerRender, "└", "├")
	headerRenderModified = strings.ReplaceAll(headerRenderModified, "┘", "┤")

	catalogueRender := lipgloss.JoinVertical(
		lipgloss.Center,
		headerRenderModified,
		content,
		seperator,
		footer,
	)

	cFinalRender := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		catalogueRender,
	)

	return cFinalRender
}
