package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

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
	list list.Model
	err  error
}

func New() *Model {
	return &Model{}
}

// initList is called when the application starts up.
// TODO: call this on tea.WindowSizeMsg
func (m *Model) initList(width, height int) {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = "To Do"
	m.list.SetItems([]list.Item{
		Task{status: todo, title: "get milk", description: "get milk from the grocery store"},
		Task{status: todo, title: "clean room", description: "tidy up bedroom on the second floor"},
		Task{status: todo, title: "lunch with friend", description: "get lunch with john doe at 3pm"},
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
		m.initList(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}

func main() {
	m := New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
