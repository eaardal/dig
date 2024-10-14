package logentrieslist

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eaardal/dig/logentry"
	"github.com/eaardal/dig/ui"
	"github.com/eaardal/dig/unicode"
	"github.com/eaardal/dig/utils"
	"github.com/eaardal/dig/viewcontroller"
	"time"
)

type Model struct {
	viewEntries                 []*viewcontroller.ViewEntry
	cursor                      int
	showClosestNearbyLogEntries bool
	showAllNearbyLogEntries     bool
}

func NewModel(viewEntries []*viewcontroller.ViewEntry) Model {
	return Model{viewEntries, len(viewEntries) - 1, false, false}
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
			m.showClosestNearbyLogEntries = !m.showClosestNearbyLogEntries
		case "f":
			m.showAllNearbyLogEntries = !m.showAllNearbyLogEntries
		}
	}

	return m, nil
}

func (m Model) View() string {
	view := ""

	for index, entry := range m.viewEntries {
		cursor := renderCursor(index == m.cursor, ui.Styles.Cursor)
		origin := renderOrigin(entry.Origin, ui.Styles.LogEntry.Origin)
		numBefore, numAfter := getNumberOfNearbyLogEntries(entry, m.showClosestNearbyLogEntries)

		if m.showClosestNearbyLogEntries && index == m.cursor {
			view += renderLogEntriesBefore(view, entry, m.showAllNearbyLogEntries)
		}

		view += formatLine(entry.LogEntry, origin, cursor, index, numBefore, numAfter)

		if m.showClosestNearbyLogEntries && index == m.cursor {
			view += renderLogEntriesAfter(view, entry, m.showAllNearbyLogEntries)
		}
	}

	view += "Press ctrl+c or q to quit"

	return view
}

func renderLogEntriesBefore(view string, entry *viewcontroller.ViewEntry, showAllNearbyLogEntries bool) string {
	nearbyLogEntriesBefore := getNearbyLogEntriesBefore(entry, showAllNearbyLogEntries)

	view += "\n"
	view += renderLogEntriesBeforeBracketTop(len(nearbyLogEntriesBefore), entry.NumLogEntriesToPreviousMatch)

	for _, entryBefore := range nearbyLogEntriesBefore {
		view += renderLogEntriesBeforeLine(formatNearbyLine(entry.Origin, entryBefore))
	}

	view += renderLogEntriesBeforeBracketBottom()
	view += "\n"

	return view
}

func renderLogEntriesAfter(view string, entry *viewcontroller.ViewEntry, showAllNearbyLogEntries bool) string {
	nearbyLogEntriesAfter := getNearbyLogEntriesAfter(entry, showAllNearbyLogEntries)

	view += "\n"
	view += renderLogEntriesAfterBracketTop(len(nearbyLogEntriesAfter), entry.NumLogEntriesToNextMatch)

	for _, entryAfter := range nearbyLogEntriesAfter {
		view += renderLogEntriesAfterLine(formatNearbyLine(entry.Origin, entryAfter))
	}

	view += renderLogEntriesAfterBracketBottom()
	view += "\n"

	return view
}

func renderLogEntriesBeforeBracketTop(showing int, total int) string {
	if showing > total {
		showing = total
	}
	return unicode.BracketTopWithText(ui.Styles.NearbyLogEntriesBracket, ui.Styles.NearbyLogEntriesBracketText, "Log entries before (showing %d/%d):", showing, total)
}

func renderLogEntriesBeforeBracketBottom() string {
	return unicode.BracketBottomWithText(ui.Styles.NearbyLogEntriesBracket, ui.Styles.NearbyLogEntriesBracketText, "End")
}

func renderLogEntriesBeforeLine(line string) string {
	return unicode.PrefixVerticalLine(ui.Styles.NearbyLogEntriesBracket, line)
}

func renderLogEntriesAfterBracketTop(showing int, total int) string {
	if showing > total {
		showing = total
	}
	return unicode.BracketTopWithText(ui.Styles.NearbyLogEntriesBracket, ui.Styles.NearbyLogEntriesBracketText, "Log entries after (showing %d/%d):", showing, total)
}

func renderLogEntriesAfterBracketBottom() string {
	return unicode.BracketBottomWithText(ui.Styles.NearbyLogEntriesBracket, ui.Styles.NearbyLogEntriesBracketText, "End")
}

func renderLogEntriesAfterLine(line string) string {
	return unicode.PrefixVerticalLine(ui.Styles.NearbyLogEntriesBracket, line)
}

func getNearbyLogEntriesBefore(entry *viewcontroller.ViewEntry, showAll bool) []*logentry.LogEntry {
	if showAll {
		return entry.LogEntriesBefore
	}

	take := entry.NumPreviousLogEntriesToShow
	if len(entry.LogEntriesBefore) < take {
		take = len(entry.LogEntriesBefore)
	}

	return entry.LogEntriesBefore[:take]
}

func getNearbyLogEntriesAfter(entry *viewcontroller.ViewEntry, showAll bool) []*logentry.LogEntry {
	if showAll {
		return entry.LogEntriesAfter
	}

	take := entry.NumNextLogEntriesToShow
	if len(entry.LogEntriesAfter) < take {
		take = len(entry.LogEntriesAfter)
	}

	return entry.LogEntriesAfter[:take]
}

func formatLine(logEntry *logentry.LogEntry, origin string, cursor string, index int, numBefore int, numAfter int) string {
	timestamp := renderTimestamp(logEntry.Time, ui.Styles.LogEntry.Timestamp)
	level := renderLevel(logEntry.Level, ui.Styles.LogEntry.Level)
	msg := renderMessage(logEntry.Message, ui.Styles.LogEntry.Message)
	lineNumber := renderLineNumber(index, numBefore, numAfter, ui.Styles.LogEntry.LineNumber)

	return fmt.Sprintf("%s %s %s - %s - %s - %s\n", cursor, lineNumber, origin, timestamp, level, msg)
}

func formatNearbyLine(origin string, logEntry *logentry.LogEntry) string {
	timestamp := renderTimestamp(logEntry.Time, ui.Styles.NearbyLogEntry.Timestamp)
	level := renderLevel(logEntry.Level, ui.Styles.NearbyLogEntry.Level)
	msg := renderMessage(logEntry.Message, ui.Styles.NearbyLogEntry.Message)

	return fmt.Sprintf("%s - %s - %s - %s\n", origin, timestamp, level, msg)
}

func renderTimestamp(timestamp string, style lipgloss.Style) string {
	parsedTime, _ := time.Parse(time.RFC3339, timestamp)
	formattedTime := parsedTime.Format("2006-01-02 15:04:05")

	if utils.IsToday(parsedTime) {
		formattedTime = parsedTime.Format("15:04:05")
	}

	return style.Render(formattedTime)
}

func renderOrigin(origin string, style lipgloss.Style) string {
	if utils.IsValueKubernetesPodID(origin) {
		return renderKubernetesPodIDOrigin(origin, ui.Styles.LogEntry.Origin)
	}

	color := ui.GetPastelColorForValue(origin)
	return style.Foreground(color).Render(origin)
}

func renderKubernetesPodIDOrigin(podID string, style lipgloss.Style) string {
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
