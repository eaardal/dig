package logentrieslist

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eaardal/dig/logentry"
	"github.com/eaardal/dig/ui"
	"github.com/eaardal/dig/viewcontroller"
	"time"
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

func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func (m Model) View() string {
	view := ""

	distinctOrigins := []string{}
	for _, entry := range m.viewEntries {
		if !contains(distinctOrigins, entry.Origin) {
			distinctOrigins = append(distinctOrigins, entry.Origin)
		}
	}
	originColorCodes := []lipgloss.Color{ui.PastelOrange, ui.PastelPurple, ui.PastelTeal, ui.PastelPink, ui.PastelLavender}
	originColors := map[string]lipgloss.Color{}

	colorIndex := 0
	for i, origin := range distinctOrigins {
		if i >= len(originColorCodes) {
			colorIndex = 0
		} else {
			colorIndex = i
		}
		originColors[origin] = originColorCodes[colorIndex]
	}

	for index, entry := range m.viewEntries {
		cursor := " "
		if index == m.cursor {
			cursor = ">"
		}

		if m.showNearbyLogEntries && index == m.cursor {
			for _, entryBefore := range entry.LogEntriesBefore {
				view += fmt.Sprintf("    ") + formatNearbyLine(entry.Origin, originColors[entry.Origin], entryBefore, -1, -1, -1)
			}
		}

		numBefore := -1
		numAfter := -1
		if m.showNearbyLogEntries {
			numBefore = entry.NumLogEntriesToPreviousMatch
			numAfter = entry.NumLogEntriesToNextMatch
		}

		view += formatLine(entry.Origin, originColors[entry.Origin], entry.LogEntry, cursor, index, numBefore, numAfter)

		if m.showNearbyLogEntries && index == m.cursor {
			for _, entryAfter := range entry.LogEntriesAfter {
				view += fmt.Sprintf("   ") + formatNearbyLine(entry.Origin, originColors[entry.Origin], entryAfter, -1, -1, -1)
			}
		}
	}

	return view
}

func isToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.YearDay() == now.YearDay()
}

func formatLine(origin string, originColor lipgloss.Color, logEntry *logentry.LogEntry, cursor string, index int, numBefore int, numAfter int) string {
	parsedTime, _ := time.Parse(time.RFC3339, logEntry.Time)
	formattedTime := parsedTime.Format("2006-01-02 15:04:05")
	if isToday(parsedTime) {
		formattedTime = parsedTime.Format("15:04:05")
	}

	app := ui.Styles.LogMessageStyles.OriginStyle.Foreground(originColor).Render(origin)
	styledTime := ui.Styles.LogMessageStyles.TimestampStyle.Render(formattedTime)
	level := ui.Styles.LogMessageStyles.LevelStyle.Render(logEntry.Level)
	msg := ui.Styles.LogMessageStyles.MessageStyle.Render(logEntry.Message)
	styledCursor := ui.Styles.LogMessageStyles.CursorStyle.Render(cursor)

	ix := fmt.Sprintf("%d.", index)
	if numBefore > -1 || numAfter > -1 {
		ix = fmt.Sprintf("%d (%d/%d)", index, numBefore, numAfter)
	}
	styledIndex := ui.Styles.LogMessageStyles.LineCountStyle.Render(ix)

	if index > -1 {
		return fmt.Sprintf("%s %s %s - %s - %s - %s\n", styledCursor, styledIndex, app, styledTime, level, msg)
	}
	return fmt.Sprintf("%s %s - %s - %s - %s\n", styledCursor, app, styledTime, level, msg)
}

func formatNearbyLine(origin string, originColor lipgloss.Color, logEntry *logentry.LogEntry, index int, numBefore int, numAfter int) string {
	parsedTime, _ := time.Parse(time.RFC3339, logEntry.Time)
	formattedTime := parsedTime.Format("2006-01-02 15:04:05")
	if isToday(parsedTime) {
		formattedTime = parsedTime.Format("15:04:05")
	}

	app := ui.Styles.LogMessageStyles.OriginStyle.Foreground(originColor).Render(origin)
	styledTime := ui.Styles.NearbyLogEntryStyles.TimestampStyle.Render(formattedTime)
	level := ui.Styles.NearbyLogEntryStyles.LevelStyle.Render(logEntry.Level)
	msg := ui.Styles.NearbyLogEntryStyles.MessageStyle.Render(logEntry.Message)

	ix := fmt.Sprintf("%d.", index)
	if numBefore > -1 || numAfter > -1 {
		ix = fmt.Sprintf("%d (%d/%d)", index, numBefore, numAfter)
	}

	if index > -1 {
		return fmt.Sprintf("  %s %s - %s - %s - %s\n", ix, app, styledTime, level, msg)
	}
	return fmt.Sprintf("  %s - %s - %s - %s\n", app, styledTime, level, msg)
}
