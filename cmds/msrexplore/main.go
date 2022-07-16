package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/u-root/u-root/pkg/msr"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		choices: []string{"Wi-Fi AP", "Wi-Fi Client", "DHCP"},
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

			// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

			// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			return m, tea.Quit
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

var style = gloss.NewStyle().
	Bold(true).
	Foreground(gloss.Color("#FAFAFA")).
	Background(gloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(22)

func (m model) View() string {
	// The header
	s := "Choose your setup\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return style.Render(s)
}

type smodel struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func (m smodel) Init() tea.Cmd {
	return m.spinner.Tick
}

type errMsg error

var (
	appStyle = gloss.NewStyle().Padding(1, 2)

	titleStyle = gloss.NewStyle().
			Foreground(gloss.Color("#FFFDF5")).
			Background(gloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = gloss.NewStyle().
				Foreground(gloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type mmodel struct {
	list    list.Model
	spinner spinner.Model
	loading bool
	err     error
}

func (m mmodel) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, tea.EnterAltScreen, m.spinner.Tick)
	return tea.Batch(cmds...)
}

func (m mmodel) View() string {
	if m.loading {
		str := fmt.Sprintf("\n\n %s ...... \n\n", m.spinner.View())
		return str
	}
	return appStyle.Render(m.list.View())
}

func (m mmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "s":
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			m.loading = !m.loading
			return m, cmd

		default:
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	// This will also call our delegate's update function.
	newListModel, lcmd := m.list.Update(msg)
	m.list = newListModel

	newSpinner, scmd := m.spinner.Update(msg)
	m.spinner = newSpinner

	cmds = append(cmds, scmd, lcmd)
	return m, tea.Batch(cmds...)
}

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

func newModel() mmodel {
	msrs := append(MSRS, msr.LockIntel...)
	items := make([]list.Item, len(msrs))
	for i, m := range msrs {
		cpus, _ := msr.AllCPUs()
		v, e := msr.MSR(m.Addr).Read(cpus)
		var title string
		if e != nil {
			title = fmt.Sprintf("--%v--", e)
		} else {
			title = fmt.Sprintf("%016x %064b", v, v)
		}
		items[i] = item{
			title:       title,
			description: fmt.Sprintf("%8x [%s]", uint(m.Addr), m.Name),
		}
	}
	d := list.NewDefaultDelegate()
	l := list.New(items, d, 0, 0)
	l.Title = "MSR explorer"
	l.Styles.Title = titleStyle
	// l.Styles.PaginationStyle = gloss.NewStyle().Foreground(gloss.Color("25"))

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = gloss.NewStyle().Foreground(gloss.Color("205"))

	return mmodel{
		list:    l,
		spinner: s,
	}
}

func main() {
	if err := tea.NewProgram(newModel()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	if false {
		p := tea.NewProgram(initialModel())
		if err := p.Start(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	}
}
