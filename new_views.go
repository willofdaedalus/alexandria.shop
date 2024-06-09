package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// func (m model) mainScreen() string {
//     c := mainRenderContent {
//         footerMsg: "this is the main screen",
//     }
//
//     return m.mainFrameRender(c)
// }

func renderHeader(w int, content string, border bool, margin ...int) string {
	border = !border
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false). //, false, true).//, false, false, false, false).
		Margin(margin...).
		Width(w).
		Align(lipgloss.Center).
		Render(content)
}

func renderItem(w int, content string) string {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1).
		Width(w).
		Align(lipgloss.Left).
		Render(content)
}

//	func (m model) mainFrameRender(content mainRenderContent) string {
//		var headers []string = make([]string, 4)
//		// var headerContent = [...]string{
//		// 	"alexandria.shop",
//		// 	"welcome, userName",
//		// 	"? for details/help",
//		// 	"c cart [10] $100.10",
//		// }
//
//		var customMargins = [][]int{
//			{0, 2, 0, 0},
//			{0, 2, 0, 2},
//			{0, 2, 0, 2},
//			{0, 0, 0, 2},
//		}
//
//		dynRenderWidth := m.termWidth - (m.termWidth / 6)
//		dynRenderHeight := m.termHeight - (m.termHeight / 4)
//		actualRenderW := dynRenderWidth - 21
//		actualRenderH := dynRenderHeight + 5
//
//		for i := 0; i < 4; i++ {
//			headers[i] = renderHeader((actualRenderW-4)/5, content.headerMsgs[i], false, customMargins[i]...)
//		}
//
//		finalHeaders := lipgloss.JoinHorizontal(lipgloss.Left, headers...)
//		header := lipgloss.NewStyle().
//			Border(lipgloss.NormalBorder()).
//			Margin(0, 1, 0, 1).
//			Width(actualRenderW - 4).
//			Align(lipgloss.Center).
//			Render(finalHeaders)
//
//		innerW := actualRenderW - 5
//
//		listSection := lipgloss.NewStyle().
//			Border(lipgloss.NormalBorder()).
//			Margin(1, 1, 0, 1).
//			Width(innerW / 3).
//			Height((actualRenderH / 4) + (actualRenderH / 2)).
//			Render(content.listSection)
//
//		bookSection := lipgloss.NewStyle().
//			Border(lipgloss.NormalBorder()).
//			Width((innerW / 3) + (innerW / 4) + (innerW / 13)).
//			Height((actualRenderH / 4) + (actualRenderH / 2)).
//			Render(content.bookDetails)
//			// bookList
//
//		midSectionJoin := lipgloss.JoinHorizontal(lipgloss.Center, listSection, bookSection)
//
//		footer := lipgloss.NewStyle().
//			Border(lipgloss.NormalBorder()).
//			Margin(1, 0, 0, 1).
//			Width(actualRenderW - 4).
//			Render(fmt.Sprint(actualRenderH))
//
//		// this was the best render height across 5 different terminal emulators!
//		h := 30
//		if actualRenderH > 33 {
//			h = 33
//		}
//
//		mainBorder := lipgloss.NewStyle().
//			Border(lipgloss.NormalBorder()).
//			Width(actualRenderW).
//			Height(h).
//			Render(header, midSectionJoin, footer)
//
//		finalRender := lipgloss.Place(
//			m.termWidth, m.termHeight,
//			lipgloss.Center, lipgloss.Center,
//			mainBorder,
//		)
//
//		return finalRender
//	}
func (m model) mainBorderRender() string {
	var (
		headers []string = make([]string, 4)

		headerContent = [...]string{
			"alexandria.shop",
			"welcome, userName",
			"? for details/help",
			"c cart [10] $100.10",
		}

		customMargins = [][]int{
			{0, 2, 0, 0},
			{0, 2, 0, 2},
			{0, 2, 0, 2},
			{0, 0, 0, 2},
		}
	)

	dynRenderWidth := m.termWidth - (m.termWidth / 6)
	dynRenderHeight := m.termHeight - (m.termHeight / 4)
	actualRenderW := dynRenderWidth - 21
	actualRenderH := dynRenderHeight + 5

	for i := 0; i < 4; i++ {
		headers[i] = renderHeader((actualRenderW-4)/5, headerContent[i], false, customMargins[i]...)
	}
	finalHeaders := lipgloss.JoinHorizontal(lipgloss.Left, headers...)

	header := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Margin(0, 1, 0, 1).
		Width(actualRenderW - 4).
		Align(lipgloss.Center).
		Render(finalHeaders)

	innerW := actualRenderW - 5
	innerH := (actualRenderH / 4) + (actualRenderH / 2)
    listSectionW := innerW / 3
    bookDeetW :=(innerW / 3) + (innerW / 4) + (innerW / 13) 

	c := innerH / 3
	items := make([]string, 0)

	for i := 0; i < c; i++ {
		items = append(items, renderItem(listSectionW - 4, fmt.Sprint(innerH, " ", c)))
	}

	itemsRender := lipgloss.JoinVertical(lipgloss.Center, items...)

	listSection := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Margin(1, 1, 0, 1).
		Padding(0, 1).
		Width(listSectionW).
		Height(innerH - (innerH % c)).
		Render(itemsRender)

	bookSection := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(bookDeetW).
		Height(innerH - (innerH % c)).
		Render(fmt.Sprint("books section ", bookDeetW))
		// bookList

	midSectionJoin := lipgloss.JoinHorizontal(lipgloss.Center, listSection, bookSection)

	footer := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Margin(1, 0, 0, 1).
		Width(actualRenderW - 4).
		Render(fmt.Sprint(actualRenderH))

	// this was the best render height across 5 different terminal emulators!
	h := 30
	if actualRenderH > 33 {
		h = 33
	}

	mainBorder := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(actualRenderW).
		Height(h - (innerH % c)).
		Render(header, midSectionJoin, footer)

	finalRender := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		mainBorder,
	)

	return finalRender
}
