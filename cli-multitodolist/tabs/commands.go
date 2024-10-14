package tabs

import tea "github.com/charmbracelet/bubbletea"

type CreateEntryMsg struct {
	Value string
}

func NewCreateEntryCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return CreateEntryMsg{Value: value}
	}
}

type CancelCreateEntryMsg struct{}

func NewCancelCreateEntryCmd() tea.Cmd {
	return func() tea.Msg {
		return CancelCreateEntryMsg{}
	}
}

type UpdateEntryMsg struct {
	Value string
}

func NewUpdateEntryCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return UpdateEntryMsg{Value: value}
	}
}

type CancelUpdateEntryMsg struct{}

func NewCancelUpdateEntryCmd() tea.Cmd {
	return func() tea.Msg {
		return CancelUpdateEntryMsg{}
	}
}

type RemoveEntryMsg struct {
	Index int
}

func NewRemoveEntryCmd(index int) tea.Cmd {
	return func() tea.Msg {
		return RemoveEntryMsg{index}
	}
}

type ConfirmRemoveEntryMsg struct {
	Index int
}

func NewConfirmRemoveEntryCmd(index int) tea.Cmd {
	return func() tea.Msg {
		return ConfirmRemoveEntryMsg{index}
	}
}
