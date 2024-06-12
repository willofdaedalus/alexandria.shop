package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/lipgloss"
)

func (m *model) mainScreen() string {
	// Calculate itemsCount before rendering
	m.itemsCount = calculateItemsCount(m.termHeight)

	headers := []string{
		fmt.Sprint(m.itemsCount),
		fmt.Sprint("welcome, ", m.curUser.username),
		"? for details/help",
		fmt.Sprintf("c cart [%d] %.2f", len(m.c.items), m.c.booksTotal()),
	}

	if m.itemsCount > magicNum {
		m.curBooks, _ = getBooksForPage(m.db, m.itemsCount, m.prevOffset)
		if m.curBooks == nil {
			log.Fatalf("no books found")
		}
	}

	titles := extractTitles(m.curBooks)

	footerMsg := "ctrl+l to logout  |  / to search  |  up/down to navigate  "
	c := mainRenderContent{
		headerContents: headers,
		bookItems:      titles,
		footerMessage:  footerMsg,
	}

    m.content = c

	return m.mainBorderRender(m.content)
}

func calculateItemsCount(termHeight int) int {
	dynRenderHeight := termHeight - (termHeight / 4)
	actualRenderH := dynRenderHeight + 5
	innerH := (actualRenderH / 4) + (actualRenderH / 2)
	return innerH / 3
}

func renderHeader(w int, content string, border bool, margin ...int) string {
	border = !border
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false). //, false, true).//, false, false, false, false).
		Margin(margin...).
		Width(w).
		Align(lipgloss.Center).
		Render(content)
}

func renderItem(w int, content string, selected bool) string {
	s := lipgloss.NewStyle().Foreground(faded.GetForeground()).Border(lipgloss.NormalBorder())

	if selected {
		s = lipgloss.NewStyle().Foreground(magenta.GetForeground()).Border(lipgloss.ThickBorder())
	}

	return lipgloss.NewStyle().
		Border(s.GetBorder()).
		BorderForeground(s.GetForeground()).
		Padding(0, 1).
		Width(w).
		Align(lipgloss.Left).
		Render(content)
}

func (m model) mainBorderRender(content mainRenderContent) string {
	var (
		headers       []string = make([]string, 0)
		customMargins          = [][]int{
			{0, 2, 0, 0},
			{0, 2, 0, 2},
			{0, 2, 0, 2},
			{0, 0, 0, 2},
		}
	)

	dynRenderWidth := m.termWidth - (m.termWidth / 6)        // calculates the best width
	dynRenderHeight := m.termHeight - (m.termHeight / 4)     // calculates the best width
	actualRenderW := dynRenderWidth - 21                     // that's the width of the main border
	actualRenderH := dynRenderHeight + 5                     // that's the height of the main border
	innerW := actualRenderW - 5                              // fake padding subtracted from th main border width
	innerH := (actualRenderH / 4) + (actualRenderH / 2)      // fake padding subtracted from th main border height
	listSectionW := innerW / 3                               // listsection width
	bookDeetW := (innerW / 3) + (innerW / 4) + (innerW / 13) // bookdetails width

	// renders headers
	for i := 0; i < 4; i++ {
		headers = append(headers, renderHeader((actualRenderW-4)/5, content.headerContents[i], false, customMargins[i]...))
	}
	// joins all header related text
	innerHeaderRender := lipgloss.JoinHorizontal(lipgloss.Left, headers...)

	// renders the header section complete with header text
	header := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Margin(0, 1, 0, 1).
		Width(actualRenderW - 4).
		Align(lipgloss.Center).
		Render(innerHeaderRender)

	isHighlighted := func(index int) bool {
		return m.curItem == index
	}
	items := make([]string, 0)
	for i := 0; i < m.itemsCount; i++ {
		items = append(items, renderItem(listSectionW-4, content.bookItems[i], isHighlighted(i)))
	}

	itemsRender := lipgloss.JoinVertical(lipgloss.Center, items...)

	midSectionJoin := renderMidSections(listSectionW, innerH, m.itemsCount, bookDeetW, itemsRender, "book deets")

	footer := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Align(lipgloss.Center).
		Margin(1, 0, 0, 1).
		Width(actualRenderW - 4).
		Render(content.footerMessage)

	// this was the best render height across 5 different terminal emulators!
	h := 30
	if actualRenderH > 33 {
		h = 33
	}

	mainBorder := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(actualRenderW).
		Height(h-(innerH%m.itemsCount)).
		Render(header, midSectionJoin, footer)

	finalRender := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		mainBorder,
	)

	return finalRender
}

func renderMidSections(listSectionW, innerH, itemCount, bookDeetW int, listContent, deetContent string) string {
	listSection := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Margin(1, 1, 0, 1).
		Padding(0, 1).
		Width(listSectionW).
		Height(innerH - (innerH % itemCount)).
		Render(listContent)

	bookSection := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(bookDeetW).
		Height(innerH - (innerH % itemCount)).
		Render(deetContent)

	return lipgloss.JoinHorizontal(lipgloss.Center, listSection, bookSection)
}

func extractTitles(bs []book) []string {
	retval := make([]string, 0)
	for _, b := range bs {
		retval = append(retval, fmt.Sprintf("%s", b.Title))
	}

	return retval
}
