package tabs

import (
	"fmt"
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

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(1, 1).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (m Model) View(getContent func(availableWidth int) string) string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.ActiveTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(fmt.Sprintf("%-10s", t)))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	availableWidth := lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize() + windowStyle.GetHorizontalPadding()
	doc.WriteString(windowStyle.Width(availableWidth).Render(getContent(availableWidth - 2)))
	return docStyle.Render(doc.String())
}

func New(tabs []string) Model {
	return Model{
		Tabs:      tabs,
		ActiveTab: 0,
		keys:      Keys,
	}
}
