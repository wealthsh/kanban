package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wealthsh/kanban/internal/task"
	"github.com/wealthsh/kanban/internal/ui"
)

func main() {
	ui.Models = []tea.Model{ui.New(), ui.NewForm(task.Todo)}
	m := ui.Models[ui.MainModel]
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
