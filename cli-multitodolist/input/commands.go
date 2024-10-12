package input

import tea "github.com/charmbracelet/bubbletea"

type FocusInputMsg struct{ id string }

func NewFocusInputCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return FocusInputMsg{id: id}
	}
}
