package help

import "github.com/charmbracelet/lipgloss"

func GetStyles(r *lipgloss.Renderer) (lipgloss.Style, lipgloss.Style) {
	keyStyle := r.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#909090",
		Dark:  "#626262",
	})

	descStyle := r.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#B2B2B2",
		Dark:  "#4A4A4A",
	})

	return keyStyle, descStyle
}
