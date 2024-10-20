package help

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type Model struct {
	renderer *lipgloss.Renderer
	ShowAll  bool
}

func New(renderer *lipgloss.Renderer) Model {
	return Model{
		renderer: renderer,
		ShowAll:  false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update() (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View(helpKeys [][]key.Binding) string {
	keyStyle, descStyle := GetStyles(m.renderer)

	if m.ShowAll {
		var cols []string
		for _, group := range helpKeys {
			var (
				keys         []string
				descriptions []string
			)
			for _, kb := range group {
				if !kb.Enabled() {
					continue
				}
				keys = append(keys, kb.Help().Key)
				descriptions = append(descriptions, kb.Help().Desc)
			}

			cols = append(cols, lipgloss.JoinHorizontal(lipgloss.Top,
				keyStyle.Render(strings.Join(keys, "\n")),
				keyStyle.Render(" "),
				descStyle.Render(strings.Join(descriptions, "\n")),
				keyStyle.Render("  "),
			))
		}
		return "\n" + lipgloss.JoinHorizontal(lipgloss.Top, cols...) + "\n"
	}
	return ""
}
