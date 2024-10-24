package textarea

import tea "github.com/charmbracelet/bubbletea"

type FocusInputMsg struct{ id, value string }

func NewFocusInputCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return FocusInputMsg{id, ""}
	}
}
func NewFocusInputValueCmd(id, value string) tea.Cmd {
	return func() tea.Msg {
		return FocusInputMsg{id, value}
	}
}
