package statusBar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusMsg struct {
	value   string
	persist bool
}
type DisablePersistMsg struct{}

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

type Model struct {
	value        string
	DefaultValue string
	persist      bool // if value should persist after next action or revert to DefaultValue
}

func New() Model {
	return Model{
		DefaultValue: "What's on your mind today ?",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StatusMsg:
		m.value = msg.value
		m.persist = msg.persist
	case tea.KeyMsg:
		if !m.persist {
			m.value = m.DefaultValue
		}
	}
	return m, nil
}

func (m Model) View(width int) string {
	w := lipgloss.Width
	statusKey := statusStyle.Render("STATUS")
	encoding := encodingStyle.Render("?")
	fishCake := fishCakeStyle.Render("toggle help")
	statusVal := statusText.
		Width(width - w(statusKey) - w(encoding) - w(fishCake)).
		Render(m.value)
	statusBar := statusBarStyle.Width(width).Render(lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		encoding,
		fishCake,
	))

	return statusBar
}
