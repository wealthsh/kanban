package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wealthsh/kanban/internal/task"
)

const divisor = 4

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

type Model struct {
	focused  task.Status
	lists    []list.Model
	err      error
	loaded   bool
	quitting bool
}

func New() *Model {
	return &Model{}
}

func (m *Model) Next() {
	if m.focused == task.Done {
		m.focused = task.Todo
	} else {
		m.focused++
	}
}

func (m *Model) Prev() {
	if m.focused == task.Todo {
		m.focused = task.Done
	} else {
		m.focused--
	}
}

func (m *Model) MoveToNext() tea.Msg {
	selected := m.lists[m.focused].SelectedItem()
	task := selected.(task.Task)
	m.lists[task.Status()].RemoveItem(m.lists[m.focused].Index())

	// Move task to next list
	task.Next()
	idx := len(m.lists[task.Status()].Items())
	m.lists[task.Status()].InsertItem(idx, list.Item(task))

	return nil
}

// initLists is called when the application starts up.
func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height/2)

	// Set this to false if you want to hide the help
	// indicators at the bottom of the terminal
	defaultList.SetShowHelp(false)

	m.lists = []list.Model{defaultList, defaultList, defaultList}

	m.lists[task.Todo].Title = "To Do"
	m.lists[task.InProgress].Title = "In Progress"
	m.lists[task.Done].Title = "Done"

	m.lists[task.Todo].SetItems([]list.Item{
		task.New(task.Todo, "get milk", "get milk from the grocery store"),
		task.New(task.Todo, "clean room", "tidy up bedroom on the second floor"),
		task.New(task.Todo, "lunch with friend", "get lunch with john doe at 3pm"),
	})
	m.lists[task.InProgress].SetItems([]list.Item{
		task.New(task.InProgress, "walk dog", "walk the dog at 8:30pm"),
		task.New(task.InProgress, "interview", "interview the cat at 10:00pm"),
	})
	m.lists[task.Done].SetItems([]list.Item{
		task.New(task.Done, "buy groceries", "buy groceries at the grocery store"),
		task.New(task.Done, "buy new gloves", "buy new gloves for winter"),
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			// Set width and height of the columns
			unfocusedStyle.Width(msg.Width / divisor)
			unfocusedStyle.Height(msg.Height - divisor)
			focusedStyle.Width(msg.Width / divisor)
			focusedStyle.Height(msg.Height - divisor)

			// Initialize the lists with mock data
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}

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
		case "n":
			Models[MainModel] = m
			Models[FormModel] = NewForm(m.focused)
			return Models[FormModel].Update(nil)
		case "d":
			// TODO: delete
		}

	case task.Task:
		task := msg
		return m, m.lists[task.Status()].InsertItem(len(m.lists[task.Status()].Items()), task)
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)

	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if m.loaded {
		todoView := m.lists[task.Todo].View()
		inProgView := m.lists[task.InProgress].View()
		doneView := m.lists[task.Done].View()

		switch m.focused {
		case task.InProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				unfocusedStyle.Render(todoView),
				focusedStyle.Render(inProgView),
				unfocusedStyle.Render(doneView),
			)
		case task.Done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				unfocusedStyle.Render(todoView),
				unfocusedStyle.Render(inProgView),
				focusedStyle.Render(doneView),
			)
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				unfocusedStyle.Render(inProgView),
				unfocusedStyle.Render(doneView),
			)
		}
	}

	return "Loading..."
}
