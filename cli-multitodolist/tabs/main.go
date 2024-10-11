package tabs

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	keys      KeyMap
	Tabs      []string
	ActiveTab int

	viewport viewport.Model
}

func New(tabs []string) Model {
	return Model{
		keys:      Keys,
		Tabs:      tabs,
		ActiveTab: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Right):
			m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
			return m, nil
		case key.Matches(msg, m.keys.Left):
			m.ActiveTab = max(m.ActiveTab-1, 0)
			return m, nil
		}
	}

	return m, nil
}

func (m Model) View(width int) string {
	var tabs []string
	for i, t := range m.Tabs {
		if m.ActiveTab == i {
			tabs = append(tabs, activeTab.Render(t))
		} else {
			tabs = append(tabs, tab.Render(t))
		}
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tabs...,
	)
	gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(header)-2)))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, header, gap)
}
