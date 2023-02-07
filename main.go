package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wealthsh/kanban/internal/ui"
)

func main() {
	m := ui.New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
