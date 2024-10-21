package tabs

import tea "github.com/charmbracelet/bubbletea"

type ChangeListMsg struct {
	Index int
}

type CreateListMsg struct {
	Value string
}

func CreateListCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return CreateListMsg{Value: value}
	}
}

type CancelCreateListMsg struct{}

func CancelCreateListCmd() tea.Cmd {
	return func() tea.Msg {
		return CancelCreateListMsg{}
	}
}

type UpdateListMsg struct {
	Value string
}

func UpdateListCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return UpdateListMsg{Value: value}
	}
}

type CancelUpdateListMsg struct{}

func CancelUpdateListCmd() tea.Cmd {
	return func() tea.Msg {
		return CancelUpdateListMsg{}
	}
}

type ConfirmRemoveListMsg struct {
	Index int
}

func ConfirmRemoveListCmd(index int) tea.Cmd {
	return func() tea.Msg {
		return ConfirmRemoveListMsg{index}
	}
}
