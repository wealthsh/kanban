package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 4

const (
	todo status = iota
	inProgress
	done
)

type Task struct {
	status      status
	title       string
	description string
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

type Model struct {
	focused status
	lists   []list.Model
	err     error
	loaded  bool
}

func New() *Model {
	return &Model{}
}

// initLists is called when the application starts up.
func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height)

	// Set this to false if you want to hide the help
	// indicators at the bottom of the terminal
	defaultList.SetShowHelp(true)

	m.lists = []list.Model{defaultList, defaultList, defaultList}

	m.lists[todo].Title = "To Do"
	m.lists[inProgress].Title = "In Progress"
	m.lists[done].Title = "Done"

	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "get milk", description: "get milk from the grocery store"},
		Task{status: todo, title: "clean room", description: "tidy up bedroom on the second floor"},
		Task{status: todo, title: "lunch with friend", description: "get lunch with john doe at 3pm"},
	})
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: todo, title: "walk dog", description: "walk the dog at 8:30pm"},
	})
	m.lists[done].SetItems([]list.Item{
		Task{status: todo, title: "shopping", description: "buy new gloves for winter"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Depending on the type of message we receive,
	// we'll update the model and/or dispatch a command.
	switch msg := msg.(type) {

	// The message we get on program startup that
	// gives us terminal dimensions.
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loaded {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.lists[todo].View(),
			m.lists[inProgress].View(),
			m.lists[done].View(),
		)
	}

	return "loading..."
}

func main() {
	m := New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
