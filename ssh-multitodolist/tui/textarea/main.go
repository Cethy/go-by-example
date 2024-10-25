package textarea

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ssh-multitodolist/tui/input"
)

// textarea use input commands and keys

type Model struct {
	Keys               input.KeyMap
	input              textarea.Model
	Active             bool
	getConfirmInputCmd func(value string) tea.Cmd
	getCancelInputCmd  func() tea.Cmd

	id string
}

func NewInputProperlyRendered(r *lipgloss.Renderer) textarea.Model {
	ta := textarea.New()
	ta.FocusedStyle, ta.BlurredStyle = defaultStyles(r)
	// set ta.style
	ta.Blur()

	return ta
}

func NewInput(placeholder, prompt string, height, charLimit int, r *lipgloss.Renderer) textarea.Model {
	ta := NewInputProperlyRendered(r)
	ta.Placeholder = placeholder

	ta.CharLimit = charLimit
	//ta.SetWidth(width)
	ta.SetHeight(height)
	ta.Prompt = ta.Prompt + prompt
	ta.ShowLineNumbers = false

	return ta
}

func New(id string, getConfirmInputCmd func(value string) tea.Cmd, getCancelInputCmd func() tea.Cmd, i textarea.Model) Model {
	return Model{
		Keys:               input.Keys,
		input:              i,
		Active:             false,
		getConfirmInputCmd: getConfirmInputCmd,
		getCancelInputCmd:  getCancelInputCmd,
		id:                 id,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case input.FocusInputMsg:
		if m.id == msg.Id {
			m.Active = true
			m.input.SetValue(msg.Value)
			m.input.Focus()
		}
	case tea.KeyMsg:
		if m.Active {
			switch {
			case key.Matches(msg, m.Keys.Enter):
				cmds = append(cmds, m.getConfirmInputCmd(m.input.Value()))
				m.Active = false
				m.input.SetValue("")
				m.input.Blur()
			case key.Matches(msg, m.Keys.Cancel):
				cmds = append(cmds, m.getCancelInputCmd())
				m.Active = false
				m.input.SetValue("")
				m.input.Blur()
			}

			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View(width int) string {
	m.input.SetWidth(width)
	return m.input.View()
}

// DefaultStyles returns the default styles for focused and blurred states for
// the textarea.
func defaultStyles(r *lipgloss.Renderer) (textarea.Style, textarea.Style) {
	focused := textarea.Style{
		Base:             r.NewStyle(),
		CursorLine:       r.NewStyle().Background(lipgloss.AdaptiveColor{Light: "255", Dark: "0"}),
		CursorLineNumber: r.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "240"}),
		EndOfBuffer:      r.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "254", Dark: "0"}),
		LineNumber:       r.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "249", Dark: "7"}),
		Placeholder:      r.NewStyle().Foreground(lipgloss.Color("240")),
		Prompt:           r.NewStyle().Foreground(lipgloss.Color("7")),
		Text:             r.NewStyle(),
	}
	blurred := textarea.Style{
		Base:             r.NewStyle(),
		CursorLine:       r.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "245", Dark: "7"}),
		CursorLineNumber: r.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "249", Dark: "7"}),
		EndOfBuffer:      r.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "254", Dark: "0"}),
		LineNumber:       r.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "249", Dark: "7"}),
		Placeholder:      r.NewStyle().Foreground(lipgloss.Color("240")),
		Prompt:           r.NewStyle().Foreground(lipgloss.Color("7")),
		Text:             r.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "245", Dark: "7"}),
	}

	return focused, blurred
}
