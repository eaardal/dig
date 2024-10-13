package logentrieslist

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eaardal/dig/logentry"
	"github.com/eaardal/dig/viewcontroller"
)

type Model struct {
	viewEntries          []*viewcontroller.ViewEntry
	cursor               int
	showNearbyLogEntries bool
}

func NewModel(viewEntries []*viewcontroller.ViewEntry) Model {
	return Model{viewEntries, 0, false}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.viewEntries)-1 {
				m.cursor++
			}
		case "d":
			m.showNearbyLogEntries = !m.showNearbyLogEntries
		}
	}

	return m, nil
}

func (m Model) View() string {
	view := ""

	for index, entry := range m.viewEntries {
		cursor := " "
		if index == m.cursor {
			cursor = ">"
		}

		if m.showNearbyLogEntries && index == m.cursor {
			for _, entryBefore := range entry.LogEntriesBefore {
				view += fmt.Sprintf("    ") + formatLine(entry.Origin, entryBefore, " ", -1, -1, -1)
			}
		}

		numBefore := -1
		numAfter := -1
		if m.showNearbyLogEntries {
			numBefore = entry.NumLogEntriesToPreviousMatch
			numAfter = entry.NumLogEntriesToNextMatch
		}
		view += formatLine(entry.Origin, entry.LogEntry, cursor, index, numBefore, numAfter)

		if m.showNearbyLogEntries && index == m.cursor {
			for _, entryAfter := range entry.LogEntriesAfter {
				view += fmt.Sprintf("   ") + formatLine(entry.Origin, entryAfter, " ", -1, -1, -1)
			}
		}
	}

	return view
}

func formatLine(origin string, logEntry *logentry.LogEntry, cursor string, index int, numBefore int, numAfter int) string {
	app := origin
	time := logEntry.Time
	level := logEntry.Level
	msg := logEntry.Message

	ix := fmt.Sprintf("%d.", index)
	if numBefore > -1 || numAfter > -1 {
		ix = fmt.Sprintf("%d (%d/%d).", index, numBefore, numAfter)
	}

	if index > -1 {
		return fmt.Sprintf("%s %s %s - %s - %s - %s\n", cursor, ix, app, time, level, msg)
	}
	return fmt.Sprintf("%s %s - %s - %s - %s\n", cursor, app, time, level, msg)
}
