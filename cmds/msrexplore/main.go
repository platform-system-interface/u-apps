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

func initialSModel() smodel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = gloss.NewStyle().Foreground(gloss.Color("205"))
	return smodel{spinner: s}
}

func (m smodel) Init() tea.Cmd {
	return m.spinner.Tick
}

type errMsg error

func (m smodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

}

func (m smodel) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	return str
}

const (
	EFER                           msr.MSR = 0xc0000080
	STAR                           msr.MSR = 0xc0000081
	LSTAR                          msr.MSR = 0xc0000082
	CSTAR                          msr.MSR = 0xc0000083
	SYSCALL_MASK                   msr.MSR = 0xc0000084
	FS_BASE                        msr.MSR = 0xc0000100
	GS_BASE                        msr.MSR = 0xc0000101
	KERNEL_GS_BASE                 msr.MSR = 0xc0000102
	TSC_AUX                        msr.MSR = 0xc0000103
	HV_X64_RESET                   msr.MSR = 0x40000003
	HV_X64_TSC_FREQUENCY           msr.MSR = 0x40000022
	HV_X64_APIC_FREQUENCY          msr.MSR = 0x40000023
	HV_X64_REENLIGHTENMENT_CONTROL msr.MSR = 0x40000106
	HV_X64_TSC_EMULATION_CONTROL   msr.MSR = 0x40000107
	HV_X64_TSC_EMULATION_STATUS    msr.MSR = 0x40000108
)

var MSRS = []msr.MSRVal{
	{Name: "EFER", Addr: EFER, Set: 0},
	{Name: "STAR", Addr: STAR, Set: 0},
	{Name: "LSTAR", Addr: LSTAR, Set: 0},
	{Name: "CSTAR", Addr: CSTAR, Set: 0},
	{Name: "SYSCALL_MASK", Addr: SYSCALL_MASK, Set: 0},
	{Name: "FS_BASE", Addr: FS_BASE, Set: 0},
	{Name: "GS_BASE", Addr: GS_BASE, Set: 0},
	{Name: "KERNEL_GS_BASE", Addr: KERNEL_GS_BASE, Set: 0},
	{Name: "TSC_AUX", Addr: TSC_AUX, Set: 0},
	/* https://readthedocs.org/projects/qemu/downloads/pdf/latest/ p280 */
	{Name: "HV_X64_RESET", Addr: HV_X64_RESET, Set: 0},
	{Name: "HV_X64_TSC_FREQUENCY", Addr: HV_X64_TSC_FREQUENCY, Set: 0},
	{Name: "HV_X64_APIC_FREQUENCY", Addr: HV_X64_APIC_FREQUENCY, Set: 0},
	{Name: "HV_X64_REENLIGHTENMENT_CONTROL", Addr: HV_X64_REENLIGHTENMENT_CONTROL, Set: 0},
	{Name: "HV_X64_TSC_EMULATION_CONTROL  ", Addr: HV_X64_TSC_EMULATION_CONTROL, Set: 0},
	{Name: "HV_X64_TSC_EMULATION_STATUS  ", Addr: HV_X64_TSC_EMULATION_STATUS, Set: 0},
}

var (
	appStyle = gloss.NewStyle().Padding(1, 2)
)

type mmodel struct {
	list list.Model
}

func (m mmodel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m mmodel) View() string {
	return appStyle.Render(m.list.View())
}

func (m mmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
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
	items := make([]list.Item, len(MSRS))
	for i, m := range MSRS {
		items[i] = item{
			title:       fmt.Sprintf("%d", m.Addr),
			description: m.Name,
		}
	}
	d := list.NewDefaultDelegate()
	l := list.New(items, d, 0, 0)
	return mmodel{
		list: l,
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
		o := tea.NewProgram(initialSModel())
		if err := o.Start(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	cpus, _ := msr.AllCPUs()
	msrs := append(MSRS, msr.LockIntel...)
	for _, m := range msrs {
		v, e := msr.MSR(m.Addr).Read(cpus)
		fmt.Printf("%v %v __ %v\n", m.Name, v, e)
	}
}
