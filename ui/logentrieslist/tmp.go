package logentrieslist

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	list          list.Model
	itemGenerator *randomItemGenerator
	keys          *listKeyMap
	delegateKeys  *delegateKeyMap
}

func newModel() model {
	var (
		itemGenerator randomItemGenerator
		delegateKeys  = newDelegateKeyMap()
		listKeys      = newListKeyMap()
	)

	// Make initial list of items
	const numItems = 24
	items := make([]list.Item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = itemGenerator.next()
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	groceryList := list.New(items, delegate, 0, 0)
	groceryList.Title = "Groceries"
	groceryList.Styles.Title = titleStyle
	groceryList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return model{
		list:          groceryList,
		keys:          listKeys,
		delegateKeys:  delegateKeys,
		itemGenerator: &itemGenerator,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.keys.insertItem):
			m.delegateKeys.remove.SetEnabled(true)
			newItem := m.itemGenerator.next()
			insCmd := m.list.InsertItem(0, newItem)
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))
			return m, tea.Batch(insCmd, statusCmd)
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

func Example() {
	rand.Seed(time.Now().UTC().UnixNano())

	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return m.NewStatusMessage(statusMessageStyle("You chose " + title))

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return m.NewStatusMessage(statusMessageStyle("Deleted " + title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}

type randomItemGenerator struct {
	titles     []string
	descs      []string
	titleIndex int
	descIndex  int
	mtx        *sync.Mutex
	shuffle    *sync.Once
}

func (r *randomItemGenerator) reset() {
	r.mtx = &sync.Mutex{}
	r.shuffle = &sync.Once{}

	r.titles = []string{
		"Artichoke",
		"Baking Flour",
		"Bananas",
		"Barley",
		"Bean Sprouts",
		"Bitter Melon",
		"Black Cod",
		"Blood Orange",
		"Brown Sugar",
		"Cashew Apple",
		"Cashews",
		"Cat Food",
		"Coconut Milk",
		"Cucumber",
		"Curry Paste",
		"Currywurst",
		"Dill",
		"Dragonfruit",
		"Dried Shrimp",
		"Eggs",
		"Fish Cake",
		"Furikake",
		"Garlic",
		"Gherkin",
		"Ginger",
		"Granulated Sugar",
		"Grapefruit",
		"Green Onion",
		"Hazelnuts",
		"Heavy whipping cream",
		"Honey Dew",
		"Horseradish",
		"Jicama",
		"Kohlrabi",
		"Leeks",
		"Lentils",
		"Licorice Root",
		"Meyer Lemons",
		"Milk",
		"Molasses",
		"Muesli",
		"Nectarine",
		"Niagamo Root",
		"Nopal",
		"Nutella",
		"Oat Milk",
		"Oatmeal",
		"Olives",
		"Papaya",
		"Party Gherkin",
		"Peppers",
		"Persian Lemons",
		"Pickle",
		"Pineapple",
		"Plantains",
		"Pocky",
		"Powdered Sugar",
		"Quince",
		"Radish",
		"Ramps",
		"Star Anise",
		"Sweet Potato",
		"Tamarind",
		"Unsalted Butter",
		"Watermelon",
		"Weißwurst",
		"Yams",
		"Yeast",
		"Yuzu",
		"Snow Peas",
	}

	r.descs = []string{
		"A little weird",
		"Bold flavor",
		"Can’t get enough",
		"Delectable",
		"Expensive",
		"Expired",
		"Exquisite",
		"Fresh",
		"Gimme",
		"In season",
		"Kind of spicy",
		"Looks fresh",
		"Looks good to me",
		"Maybe not",
		"My favorite",
		"Oh my",
		"On sale",
		"Organic",
		"Questionable",
		"Really fresh",
		"Refreshing",
		"Salty",
		"Scrumptious",
		"Delectable",
		"Slightly sweet",
		"Smells great",
		"Tasty",
		"Too ripe",
		"At last",
		"What?",
		"Wow",
		"Yum",
		"Maybe",
		"Sure, why not?",
	}

	r.shuffle.Do(func() {
		shuf := func(x []string) {
			rand.Shuffle(len(x), func(i, j int) { x[i], x[j] = x[j], x[i] })
		}
		shuf(r.titles)
		shuf(r.descs)
	})
}

func (r *randomItemGenerator) next() item {
	if r.mtx == nil {
		r.reset()
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	i := item{
		title:       r.titles[r.titleIndex],
		description: r.descs[r.descIndex],
	}

	r.titleIndex++
	if r.titleIndex >= len(r.titles) {
		r.titleIndex = 0
	}

	r.descIndex++
	if r.descIndex >= len(r.descs) {
		r.descIndex = 0
	}

	return i
}

//
//import (
//	"fmt"
//	"os"
//
//	"github.com/charmbracelet/bubbles/list"
//	tea "github.com/charmbracelet/bubbletea"
//	"github.com/charmbracelet/lipgloss"
//)
//
//var docStyle = lipgloss.NewStyle().Margin(1, 1)
//
//type item struct {
//	title, descrip string
//}
//
//func (i item) Title() string       { return i.title }
//func (i item) Description() string { return i.descrip }
//func (i item) FilterValue() string { return i.title }
//
//type model struct {
//	list list.Model
//}
//
//func (m model) Init() tea.Cmd {
//	return nil
//}
//
//func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.KeyMsg:
//		if msg.String() == "ctrl+c" {
//			return m, tea.Quit
//		}
//	case tea.WindowSizeMsg:
//		h, v := docStyle.GetFrameSize()
//		m.list.SetSize(msg.Width-h, msg.Height-v)
//	}
//
//	var cmd tea.Cmd
//	m.list, cmd = m.list.Update(msg)
//	return m, cmd
//}
//
//func (m model) View() string {
//	return docStyle.Render(m.list.View())
//}
//
//func Example() {
//	items := []list.Item{
//		item{title: "Raspberry Pi’s", descrip: "I have ’em all over my house"},
//		item{title: "Nutella", descrip: "It's good on toast"},
//		item{title: "Bitter melon", descrip: "It cools you down"},
//		item{title: "Nice socks", descrip: "And by that I mean socks without holes"},
//		item{title: "Eight hours of sleep", descrip: "I had this once"},
//		item{title: "Cats", descrip: "Usually"},
//		item{title: "Plantasia, the album", descrip: "My plants love it too"},
//		item{title: "Pour over coffee", descrip: "It takes forever to make though"},
//		item{title: "VR", descrip: "Virtual reality...what is there to say?"},
//		item{title: "Noguchi Lamps", descrip: "Such pleasing organic forms"},
//		item{title: "Linux", descrip: "Pretty much the best OS"},
//		item{title: "Business school", descrip: "Just kidding"},
//		item{title: "Pottery", descrip: "Wet clay is a great feeling"},
//		item{title: "Shampoo", descrip: "Nothing like clean hair"},
//		item{title: "Table tennis", descrip: "It’s surprisingly exhausting"},
//		item{title: "Milk crates", descrip: "Great for packing in your extra stuff"},
//		item{title: "Afternoon tea", descrip: "Especially the tea sandwich part"},
//		item{title: "Stickers", descrip: "The thicker the vinyl the better"},
//		item{title: "20° Weather", descrip: "Celsius, not Fahrenheit"},
//		item{title: "Warm light", descrip: "Like around 2700 Kelvin"},
//		item{title: "The vernal equinox", descrip: "The autumnal equinox is pretty good too"},
//		item{title: "Gaffer’s tape", descrip: "Basically sticky fabric"},
//		item{title: "Terrycloth", descrip: "In other words, towel fabric"},
//	}
//
//	delg := list.NewDefaultDelegate()
//	delg.SetSpacing(1)
//
//	ls := list.New(items, delg, 0, 0)
//
//	m := model{list: ls}
//	m.list.Title = "My Fave Things"
//
//	p := tea.NewProgram(m, tea.WithAltScreen())
//
//	if _, err := p.Run(); err != nil {
//		fmt.Println("Error running program:", err)
//		os.Exit(1)
//	}
//}
