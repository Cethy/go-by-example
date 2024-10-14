package todolist

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

type MoveItemUpMsg struct {
	cursor int
}

func NewMoveItemUpCmd(cursor int) tea.Cmd {
	return func() tea.Msg {
		return MoveItemUpMsg{cursor}
	}
}

type MoveItemDownMsg struct {
	cursor int
}

func NewMoveItemDownCmd(cursor int) tea.Cmd {
	return func() tea.Msg {
		return MoveItemDownMsg{cursor}
	}
}
