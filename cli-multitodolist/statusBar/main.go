package statusBar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View(status string, width int) string {
	w := lipgloss.Width
	statusKey := statusStyle.Render("STATUS")
	encoding := encodingStyle.Render("?")
	fishCake := fishCakeStyle.Render("toggle help")
	statusVal := statusText.
		Width(width - w(statusKey) - w(encoding) - w(fishCake)).
		Render(status)
	statusBar := statusBarStyle.Width(width).Render(lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		encoding,
		fishCake,
	))

	return statusBar
}
