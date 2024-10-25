package statusBar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	renderer     *lipgloss.Renderer
	value        string
	DefaultValue string
	persist      bool // if value should persist after next action or revert to DefaultValue
}

func New(username string, renderer *lipgloss.Renderer) Model {
	d := "What's on your mind today ?"
	if username != "" {
		d = "Hi " + username + ", what's on your mind today ?"
	}

	return Model{
		renderer:     renderer,
		DefaultValue: d,
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
	statusStyle, encodingStyle, fishCakeStyle, statusText, statusBarStyle := GetStyles(m.renderer)
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
