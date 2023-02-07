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

// Styling
var (
	unfocusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.HiddenBorder())
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
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

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

type Model struct {
	focused  status
	lists    []list.Model
	err      error
	loaded   bool
	quitting bool
}

func New() *Model {
	return &Model{}
}

func (m *Model) MoveToNext() tea.Msg {
	selected := m.lists[m.focused].SelectedItem()
	task := selected.(Task)
	m.lists[task.status].RemoveItem(m.lists[m.focused].Index())

	// Move task to next list
	task.Next()
	idx := len(m.lists[task.status].Items())
	m.lists[task.status].InsertItem(idx, list.Item(task))

	return nil
}

// Go to next list
func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

// Go to prev list
func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

// initLists is called when the application starts up.
func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height/2)

	// Set this to false if you want to hide the help
	// indicators at the bottom of the terminal
	defaultList.SetShowHelp(false)

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
		Task{status: todo, title: "walk cat", description: "walk the cat at 10:00pm"},
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
			// Set width and height of the columns
			unfocusedStyle.Width(msg.Width / divisor)
			unfocusedStyle.Height(msg.Height - divisor)
			focusedStyle.Width(msg.Width / divisor)
			focusedStyle.Height(msg.Height - divisor)

			// Initialize the lists
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}

	// Detect keystrokes
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "left", "h":
			m.Prev()

		case "right", "l":
			m.Next()

		case "enter":
			m.MoveToNext()
		}
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	// Don't render anything if we're quitting, makes the
	// terminal clean after exiting the application.
	if m.quitting {
		return ""
	}

	if m.loaded {
		todoView := m.lists[todo].View()
		inProgView := m.lists[inProgress].View()
		doneView := m.lists[done].View()

		switch m.focused {
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				unfocusedStyle.Render(inProgView),
				unfocusedStyle.Render(doneView),
			)
		case inProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				unfocusedStyle.Render(todoView),
				focusedStyle.Render(inProgView),
				unfocusedStyle.Render(doneView),
			)
		case done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				unfocusedStyle.Render(todoView),
				unfocusedStyle.Render(inProgView),
				focusedStyle.Render(doneView),
			)
		}
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
