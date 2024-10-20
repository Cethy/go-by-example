package tabs

import "github.com/charmbracelet/lipgloss"

func GetStyles(r *lipgloss.Renderer) (lipgloss.Style, lipgloss.Style, lipgloss.Style) {
	highlight := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	// Tabs.

	activeTabBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab := r.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 1)

	activeTab := tab.Border(activeTabBorder, true)

	tabGap := tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	return tab, activeTab, tabGap
}
