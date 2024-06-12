package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// the initial screen to display when the program first runs
func (m model) initialScreen() string {
	// center the ascii art within the terminal window
	return lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		cyan.Align(lipgloss.Center).Render(welcomeAscii),
	)
}

func (m model) renderAuthScreen(title, helpText string, inputs []textinput.Model, curField int) string {
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
		if i >= len(inputs) {
			break
		}
		inputRenders = append(inputRenders, m.renderInputBox(label, i, inputs, curField == i))
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

func (m model) signUpScreen() string {
	return m.renderAuthScreen(
		signUpText,
		"press ctrl+l to log in | press ctrl+c to quit",
		m.signupInputs,
		m.signupCurField,
	)
}

func (m model) loginScreen() string {
	return m.renderAuthScreen(
		loginText,
		"press ctrl+s to sign up | press ctrl+c to quit",
		m.loginInputs,
		m.loginCurField,
	)
}

func (m model) infoScreen(info string, w, h int) string {
	infoRender := noBorderStyle.
		Padding(0, 2, 0, 2).
		Width(w).Height(h).
		// Width(50).Height(3).
		Align(lipgloss.Center).Render(info)

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

func (m model) bookDetailsScreen(header1 string) string {
	var (
		footerMsg string
	)

	// Initialize variables
	renderWidth := (m.termWidth / 2) + 10
	if renderWidth < 0 {
		renderWidth = 0
	}

	booksDetailsRender := lipgloss.NewStyle().
		Foreground(yellow.GetForeground()).
		Width(50).
		MarginBottom(2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		PaddingLeft(1).
		Render(fmt.Sprintf("%s by %s\nPRICE: $%.2f\nYEAR: %d\n%s",
			selectedBook.Title,
			selectedBook.Author,
			selectedBook.Price,
			selectedBook.Year,
			selectedBook.Genre,
		),
		)

	longDescRender := lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder(), false, false, false, false).
		Height(4).
		Render(selectedBook.LongDesc)

	fBookRender := lipgloss.JoinVertical(
		lipgloss.Top,
		booksDetailsRender,
		longDescRender,
	)

	if m.c.bookInCart(selectedBook) {
		footerMsg = removeFromCartMsg
	} else {
		footerMsg = addToCartMsg
	}

	catalogueView := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		Padding(3, 5, 0, 5).
		BorderBottom(false).
		Width(renderWidth).
		Height(catalogueViewHeight - 8).
		Align(lipgloss.Left).
		Render(fBookRender)

	return m.mainScreenFrame(header1, footerMsg, catalogueView)
}

// this renders the entire catalogue view
// big thanks to @its_gaurav on the Charm CLI Discord!
func (m model) catalogueScreen(header1 string) string {
	// Initialize variables
	renderWidth := (m.termWidth / 2) + 10
	if renderWidth < 0 {
		renderWidth = 0
	}

	mainHeight := m.termHeight - 2
	offset := (mainHeight - 20) / 3

	// function to determine if an item is highlighted
	isHighlighted := func(index int) bool {
		return m.curItem == index
	}

	// render the top, mid, and bot items based on current item
	var itemsRender []string
	for i := 0; i < m.itemsCount; i++ {
		itemsRender = append(itemsRender, m.renderItemDisplay(renderWidth, offset, isHighlighted(i), m.curBooks[i]))
	}

	catalogueView := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderBottom(false).
		Width(renderWidth).
		Height(mainHeight - 20).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Center, itemsRender...))

	return m.mainScreenFrame(header1, catalogueHelpMsg, catalogueView)
}

func (m model) helpScreen() string {
	infoRender := noBorderStyle.
		PaddingTop(1).
		Width(80).Height(10).
		Padding(0, 2, 0, 2).
		Align(lipgloss.Left).Render(helpText)

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

func (m model) cartScreenView() string {
	var items []string

	for key, value := range m.c.items {
		// Format the key-value pair as "key - value"
		pair := key + " - " + fmt.Sprintf("%.2f", value)
		// Append the formatted string to the slice
		items = append(items, pair)
	}
    
    itemsRender := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Width(40).Render(strings.Join(items, "\n"))

    return m.mainScreenFrame("esc to go back", "esc to go back", itemsRender)
}
