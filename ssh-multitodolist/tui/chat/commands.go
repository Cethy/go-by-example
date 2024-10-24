package chat

import tea "github.com/charmbracelet/bubbletea"

type CreateMessageMsg struct {
	Message string
}

func CreateMessageCmd(message string) tea.Cmd {
	return func() tea.Msg {
		return CreateMessageMsg{Message: message}
	}
}

type CancelCreateMessageMsg struct{}

func CancelCreateMessageCmd() tea.Cmd {
	return func() tea.Msg {
		return CancelCreateMessageMsg{}
	}
}
