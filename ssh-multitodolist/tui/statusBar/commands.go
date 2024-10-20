package statusBar

import tea "github.com/charmbracelet/bubbletea"

type StatusMsg struct {
	value   string
	persist bool
}

func NewStatusMsg(value string) StatusMsg {
	return StatusMsg{value: value}
}

func NewStatusCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return StatusMsg{value: value, persist: false}
	}
}
func NewPersistingStatusCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return StatusMsg{value: value, persist: true}
	}
}
