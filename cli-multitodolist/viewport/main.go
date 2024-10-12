package viewport

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
	"strings"
)

type Model struct {
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update() (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View(content string, width, height, focus int) string {
	// @todo properly init the app (first render doesn't have the width/height values set)
	if height < 0 {
		return content
	}

	wrapped := wrap.String(wordwrap.String(content, width), width)

	lines := strings.Split(wrapped, "\n")
	if len(lines) < height {
		return lipgloss.NewStyle().
			Width(width).
			Height(height).
			Align(lipgloss.Left, lipgloss.Top).
			Render(strings.Join(lines, "\n"))
	}

	// make sure the cursor is always visible
	if focus+1 >= height {
		return "↑\n" + strings.Join(lines[focus+3-height:focus+2-1], "\n") + viewportTail(len(lines), width, focus)
	}
	return strings.Join(lines[0:height-1], "\n") + viewportTail(len(lines), width, focus)
}

func viewportTail(contentLen, width, focus int) string {
	tail := " "
	if focus+2 < contentLen {
		tail = "↓"
	}
	currentPos := fmt.Sprintf("%d/%d", min(focus+1, contentLen-1), contentLen-1)

	return "\n" + lipgloss.JoinHorizontal(lipgloss.Top, tail, lipgloss.NewStyle().Width(width-1-len(currentPos)).Render(""), currentPos)
}
