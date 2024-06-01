package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

func (m model) signUpScreen() string {
	var (
		layout            strings.Builder
		login, pwd, rePwd string
	)

	// the big "LOGIN" text at the top
	loginPrompt := lipgloss.NewStyle().
		Foreground(cyan.GetForeground()).
		Width(55).Height(5).
		Align(lipgloss.Center).
		Render(signUpText)

	// variables for the username and password prompt boxes
	loginRender := renderBoxDesc("username", 0, m.signupInputs)
	pwdRender := renderBoxDesc("password", 1, m.signupInputs)
	rePwdRender := renderBoxDesc("password", 2, m.signupInputs)

	// change the color of each render based on the current focus
	if m.signupCurField == 0 {
		login = magenta.PaddingLeft(8).Align(lipgloss.Left).Render(loginRender)
		pwd = faded.PaddingLeft(8).Align(lipgloss.Left).Render(pwdRender)
		rePwd = faded.PaddingLeft(8).Align(lipgloss.Left).Render(rePwdRender)
	} else if m.signupCurField == 1 {
		login = faded.PaddingLeft(8).Align(lipgloss.Left).Render(loginRender)
		pwd = magenta.PaddingLeft(8).Align(lipgloss.Left).Render(pwdRender)
		rePwd = faded.PaddingLeft(8).Align(lipgloss.Left).Render(rePwdRender)
	} else {
		login = faded.PaddingLeft(8).Align(lipgloss.Left).Render(loginRender)
		pwd = faded.PaddingLeft(8).Align(lipgloss.Left).Render(pwdRender)
		rePwd = magenta.PaddingLeft(8).Align(lipgloss.Left).Render(rePwdRender)
	}

	// footer/help message render
	helpText := "press ctrl+l to log in | press ctrl+c to quit"
	helpBox := noBorderStyle.PaddingTop(1).Width(50).Align(lipgloss.Bottom).Render(helpText)

	// join the various fields together;
	// first the input boxes and then those and the login prompt
	textFields := lipgloss.JoinVertical(lipgloss.Left, login, pwd, rePwd, helpBox)
	ui := lipgloss.JoinVertical(lipgloss.Left, loginPrompt, textFields)

	// render the actual with dialogBoxStyle but this simply "puts" the render
	// in the center of the screen no matter what
	dialog := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
	)

	// everything onscreen is a string so tie it up nice with a bow and return a string
	layout.WriteString(dialog)
	return layout.String()
}

// renders the login screen
func (m model) loginScreen() string {
	var (
		layout     strings.Builder
		login, pwd string
	)

	// the big "LOGIN" text at the top
	loginPrompt := lipgloss.NewStyle().
		Foreground(cyan.GetForeground()).
		Width(55).Height(5).
		Align(lipgloss.Center).
		Render(loginText)

		// variables for the username and password prompt boxes
	loginRender := renderBoxDesc("username", 0, m.loginInputs)
	pwdRender := renderBoxDesc("password", 1, m.loginInputs)

	// change the color of each render based on the current focus
	if m.loginCurField == 0 {
		login = magenta.PaddingLeft(8).Align(lipgloss.Left).Bold(true).Render(loginRender)
		pwd = faded.PaddingLeft(8).Align(lipgloss.Left).Bold(true).Render(pwdRender)
	} else {
		login = faded.PaddingLeft(8).Align(lipgloss.Left).Bold(true).Render(loginRender)
		pwd = magenta.PaddingLeft(8).Align(lipgloss.Left).Bold(true).Render(pwdRender)
	}

	// footer/help message render
	helpText := "press ctrl+s to sign up | press ctrl+c to quit"
	helpBox := noBorderStyle.Width(50).Height(1).Align(lipgloss.Bottom).Render(helpText)

	// join the various fields together;
	// first the input boxes and then those and the login prompt
	textFields := lipgloss.JoinVertical(lipgloss.Left, login, pwd, helpBox)
	ui := lipgloss.JoinVertical(lipgloss.Left, loginPrompt, textFields)

	// render the actual with dialogBoxStyle but this simply "puts" the render
	// in the center of the screen no matter what
	dialog := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
	)

	// everything onscreen is a string so tie it up nice with a bow and return a string
	layout.WriteString(dialog)
	return layout.String()
}

func (m model) infoScreen(info string) string {
	infoRender := noBorderStyle.
		PaddingTop(1).
		Width(50).Height(3).
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

// this renders the entire catalogue view
// big thanks to @its_gaurav on the Charm CLI Discord!
func (m model) catalogueScreen(curUser string) string {
	var (
		top string
		mid string
		bot string
	)

	// renderWidth := m.termWidth - 22
	renderWidth := 110
	if renderWidth < 0 {
		renderWidth = 0
	}

	headerRender := m.renderHeaders(curUser, "29-05-24", "cart [16]")
	footer := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		Width(renderWidth).
		Height(2).
		Align(lipgloss.Center).
		Render("ctrl+c to exit  |  vim or arrow keys to navigate\nc for cart  |  ? for help/details")

	mainHeight := m.termHeight - lipgloss.Height(headerRender) - lipgloss.Height(footer) - 1
	offset := 4

	if m.curItem == 0 {
		top = renderItemDisplay(renderWidth, mainHeight/offset, true)
		mid = renderItemDisplay(renderWidth, mainHeight/offset, false)
		bot = renderItemDisplay(renderWidth, mainHeight/offset, false)
	} else if m.curItem == 1 {
		top = renderItemDisplay(renderWidth, mainHeight/offset, false)
		mid = renderItemDisplay(renderWidth, mainHeight/offset, true)
		bot = renderItemDisplay(renderWidth, mainHeight/offset, false)
	} else if m.curItem == 2 {
		top = renderItemDisplay(renderWidth, mainHeight/offset, false)
		mid = renderItemDisplay(renderWidth, mainHeight/offset, false)
		bot = renderItemDisplay(renderWidth, mainHeight/offset, true)
	}

	itemsRender := lipgloss.JoinVertical(lipgloss.Center, top, mid, bot)

	catalogueView := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderBottom(false).
		Width(renderWidth).
		Align(lipgloss.Center).
		Render(itemsRender)

	seperator := "├" + strings.Repeat("─", renderWidth) + "┤"

	headerRenderModified := strings.ReplaceAll(headerRender, "└", "├")
	headerRenderModified = strings.ReplaceAll(headerRenderModified, "┘", "┤")

	catalogueRender := lipgloss.JoinVertical(
		lipgloss.Bottom,
		headerRenderModified,
		catalogueView,
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

// renders an item in the catalogue view
func renderItemDisplay(renderWidth, offset int, selected bool) string {
	selectedColour := lipgloss.NewStyle().Foreground(faded.GetForeground())

	if selected {
		selectedColour = lipgloss.NewStyle().Foreground(magenta.GetForeground())
	}

	itemContent := lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder()).
		Foreground(selectedColour.GetForeground()).
		Render("book name here  |  author name here  |  pub date here  ")

	itemDesc := lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder()).
		Foreground(selectedColour.GetForeground()).
		Render("the description comes in here right now")

	contentRender := lipgloss.JoinVertical(lipgloss.Top, itemContent, itemDesc)

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(selectedColour.GetForeground()).
		Width(renderWidth - 5).Height(offset).
		Render(contentRender)
}

func (m model) renderHeaders(curUser, timeDate, cart string) string {

	tops := [][]string{
		{
			"alexandria.shop",
			fmt.Sprintf("welcome, %s", curUser),
			fmt.Sprintf("current date is %s", timeDate),
			cart,
		}, // actual headers
	}

	headerTable := table.New().
		Border(lipgloss.NormalBorder()).
		Width(112).
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
func renderBoxDesc(s string, idx int, inputs []textinput.Model) string {
	desc := noBorderStyle.Render(s)
	// this is the content from the input box as we type
	inputBox := textBoxStyle.Render(inputs[idx].View())
	finalRender := lipgloss.JoinHorizontal(lipgloss.Left, desc, inputBox)

	return finalRender
}
