package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wealthsh/kanban/internal/task"
)

type Form struct {
	focused     task.Status
	title       textinput.Model
	description textarea.Model
}

func NewForm(focused task.Status) *Form {
	form := &Form{
		focused:     focused,
		title:       textinput.New(),
		description: textarea.New(),
	}
	form.title.Focus()
	return form
}

func (m Form) Init() tea.Cmd {
	return nil
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			} else {
				// Finish filling out the form, so save the
				// task and return to the main view.
				Models[FormModel] = m
				return Models[MainModel], m.CreateTask
			}
		}

		var cmd tea.Cmd

		if m.title.Focused() {
			m.title, cmd = m.title.Update(msg)
			return m, cmd
		} else {
			m.description, cmd = m.description.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m Form) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		"Create New Task\n",
		m.title.View(),
		m.description.View(),
	)
}

func (m Form) CreateTask() tea.Msg {
	task := task.New(m.focused, m.title.Value(), m.description.Value())
	return task
}
