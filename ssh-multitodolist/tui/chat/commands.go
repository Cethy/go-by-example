package chat

import tea "github.com/charmbracelet/bubbletea"

type AddMessageMsg struct {
	Message string
}

func AddMessageCmd(message string) tea.Cmd {
	return func() tea.Msg {
		return AddMessageMsg{Message: message}
	}
}

type CancelAddMessageMsg struct{}

func CancelAddMessageCmd() tea.Cmd {
	return func() tea.Msg {
		return CancelAddMessageMsg{}
	}
}
