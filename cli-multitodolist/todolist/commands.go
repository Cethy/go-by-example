package todolist

import tea "github.com/charmbracelet/bubbletea"

type NewEntryMsg struct {
	Value string
}

func NewNewEntryCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return NewEntryMsg{Value: value}
	}
}

type CancelNewEntryMsg struct{}

func NewCancelNewEntryCmd() tea.Cmd {
	return func() tea.Msg {
		return CancelNewEntryMsg{}
	}
}
