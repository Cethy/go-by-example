package todolist

import tea "github.com/charmbracelet/bubbletea"

type CreateEntryMsg struct {
	Value string
}

func CreateEntryCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return CreateEntryMsg{Value: value}
	}
}

type CancelCreateEntryMsg struct{}

func CancelCreateEntryCmd() tea.Cmd {
	return func() tea.Msg {
		return CancelCreateEntryMsg{}
	}
}

type UpdateEntryMsg struct {
	Value string
}

func UpdateEntryCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return UpdateEntryMsg{Value: value}
	}
}

type CancelUpdateEntryMsg struct{}

func CancelUpdateEntryCmd() tea.Cmd {
	return func() tea.Msg {
		return CancelUpdateEntryMsg{}
	}
}

type MoveItemUpMsg struct {
	cursor int
}

func MoveItemUpCmd(cursor int) tea.Cmd {
	return func() tea.Msg {
		return MoveItemUpMsg{cursor}
	}
}

type MoveItemDownMsg struct {
	cursor int
}

func MoveItemDownCmd(cursor int) tea.Cmd {
	return func() tea.Msg {
		return MoveItemDownMsg{cursor}
	}
}
