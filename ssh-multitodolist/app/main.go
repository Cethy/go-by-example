package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Shared state of the application

type NewDataMsg struct{}

type App struct {
	Programs []*tea.Program
}

func (b *App) NotifyNewData() {
	for _, p := range b.Programs {
		go p.Send(NewDataMsg{})
	}
}
