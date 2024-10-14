package logentrieslist

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eaardal/dig/logentry"
	"github.com/eaardal/dig/ui"
	"github.com/eaardal/dig/utils"
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

	for index, entry := range m.viewEntries {
		cursor := renderCursor(index == m.cursor, ui.Styles.CursorStyle)
		origin := renderOrigin(entry.Origin, ui.Styles.LogEntryStyles.OriginStyle)

		if m.showNearbyLogEntries && index == m.cursor {
			for _, entryBefore := range entry.LogEntriesBefore {
				view += formatNearbyLine(entry.Origin, entryBefore)
			}
		}

		numBefore, numAfter := getNumberOfNearbyLogEntries(entry, m.showNearbyLogEntries)
		view += formatLine(entry.LogEntry, origin, cursor, index, numBefore, numAfter)

		if m.showNearbyLogEntries && index == m.cursor {
			for _, entryAfter := range entry.LogEntriesAfter {
				view += formatNearbyLine(entry.Origin, entryAfter)
			}
		}
	}

	return view
}

func formatLine(logEntry *logentry.LogEntry, origin string, cursor string, index int, numBefore int, numAfter int) string {
	timestamp := renderTimestamp(logEntry.Time, ui.Styles.LogEntryStyles.TimestampStyle)
	level := renderLevel(logEntry.Level, ui.Styles.LogEntryStyles.LevelStyle)
	msg := renderMessage(logEntry.Message, ui.Styles.LogEntryStyles.MessageStyle)
	lineNumber := renderLineNumber(index, numBefore, numAfter, ui.Styles.LogEntryStyles.LineNumberStyle)

	return fmt.Sprintf("%s %s %s - %s - %s - %s\n", cursor, lineNumber, origin, timestamp, level, msg)
}

func formatNearbyLine(origin string, logEntry *logentry.LogEntry) string {
	timestamp := renderTimestamp(logEntry.Time, ui.Styles.NearbyLogEntryStyles.TimestampStyle)
	level := renderLevel(logEntry.Level, ui.Styles.NearbyLogEntryStyles.LevelStyle)
	msg := renderMessage(logEntry.Message, ui.Styles.NearbyLogEntryStyles.MessageStyle)

	return fmt.Sprintf("\t%s - %s - %s - %s\n", origin, timestamp, level, msg)
}

func renderTimestamp(timestamp string, style lipgloss.Style) string {
	parsedTime, _ := time.Parse(time.RFC3339, timestamp)
	formattedTime := parsedTime.Format("2006-01-02 15:04:05")

	if isToday(parsedTime) {
		formattedTime = parsedTime.Format("15:04:05")
	}

	return style.Render(formattedTime)
}

func renderOrigin(origin string, style lipgloss.Style) string {
	if utils.IsValueKubernetesPodID(origin) {
		return renderKubernetesPodNameOrigin(origin, ui.Styles.LogEntryStyles.OriginStyle)
	}

	color := ui.GetPastelColorForValue(origin)
	return style.Foreground(color).Render(origin)
}

func renderKubernetesPodNameOrigin(podID string, style lipgloss.Style) string {
	_, deploymentID, replicaSetID := utils.SplitIntoKubernetesPodIDParts(podID)
	deploymentColor := ui.GetPastelColorForValue(deploymentID)
	replicaSetColor := ui.GetPastelColorForValue(replicaSetID)
	return style.Render(fmt.Sprintf("%s-%s", style.Foreground(deploymentColor).Render(deploymentID), style.Foreground(replicaSetColor).Render(replicaSetID)))
}

func renderLevel(level string, style lipgloss.Style) string {
	return style.Render(level)
}

func renderMessage(message string, style lipgloss.Style) string {
	return style.Render(message)
}

func renderCursor(cursor bool, style lipgloss.Style) string {
	if cursor {
		return style.Render(">")
	}
	return " "
}

func renderLineNumber(index int, numBefore, numAfter int, style lipgloss.Style) string {
	formattedIndex := fmt.Sprintf("%d.", index)

	if numBefore > -1 || numAfter > -1 {
		formattedIndex = fmt.Sprintf("%d (%d/%d)", index, numBefore, numAfter)
	}

	return style.Render(formattedIndex)
}

func getNumberOfNearbyLogEntries(entry *viewcontroller.ViewEntry, showNearby bool) (int, int) {
	numBefore := -1
	numAfter := -1

	if showNearby {
		numBefore = entry.NumLogEntriesToPreviousMatch
		numAfter = entry.NumLogEntriesToNextMatch
	}

	return numBefore, numAfter
}

func isToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.YearDay() == now.YearDay()
}
