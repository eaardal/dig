package interactiveselectlist

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type ListItem struct {
	Value      string
	IsSelected bool
}

type Model struct {
	choices    []string
	cursor     int
	selected   map[int]struct{}
	headerText string
}

func NewModel(listItems []ListItem, headerText string) Model {
	var choices []string
	selected := make(map[int]struct{})

	for i, listItem := range listItems {
		choices = append(choices, listItem.Value)
		if listItem.IsSelected {
			selected[i] = struct{}{}
		}
	}

	return Model{
		choices:    choices,
		selected:   selected,
		headerText: headerText,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, tea.Quit
		case "ctrl+c", "q":
			m.selected = make(map[int]struct{})
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	// Render the header
	view := fmt.Sprintf("%s\n\n", m.headerText)
	view += "Use up/down or k/j to navigate, space to select/unselect, enter to confirm choices and q to cancel.\n\n"

	for i, choice := range m.choices {
		// Is the cursor pointing at this choice?
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		// Is this choice selected?
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		// Render the row
		view += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// Render the footer
	view += "\nPress q to quit.\n"

	return view
}

func (m Model) GetSelectedChoices() []string {
	var selectedChoices []string
	for i := range m.selected {
		selectedChoices = append(selectedChoices, m.choices[i])
	}
	return selectedChoices
}
