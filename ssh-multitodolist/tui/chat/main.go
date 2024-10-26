package chat

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
	"ssh-multitodolist/app"
	"ssh-multitodolist/app/state"
	"ssh-multitodolist/tui/input"
	"ssh-multitodolist/tui/statusBar"
	"strings"
)

type Model struct {
	state    *state.State
	app      *app.App
	renderer *lipgloss.Renderer
	Keys     KeyMap
}

func New(s *state.State, a *app.App, r *lipgloss.Renderer) Model {
	return Model{
		state:    s,
		app:      a,
		renderer: r,
		Keys:     keys,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg, isAnyInputActive bool) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case AddMessageMsg:
		m.app.AddChatMessage(msg.Message, m.state.Username, m.state.Color)
		cmds = append(cmds, statusBar.NewStatusCmd("Message sent"))

	case CancelAddMessageMsg:
		cmds = append(cmds, statusBar.NewStatusCmd("Message cancelled"))
	case tea.KeyMsg:
		if !isAnyInputActive {
			switch {
			case key.Matches(msg, m.Keys.AddMessage):
				cmds = append(cmds, input.NewFocusInputCmd("newChatMessageInput"))
				cmds = append(cmds, statusBar.NewPersistingStatusCmd("Typing message"))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View(newChatMessageRender func() string, width, height int) string {
	containerStyle := m.renderer.NewStyle().
		Height(height).
		Width(width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		BorderLeft(true)
	title := m.renderer.NewStyle().
		Width(width).
		Underline(true).
		Align(lipgloss.Center, lipgloss.Center).
		Render("\nChat")

	textinput := newChatMessageRender()

	conversation := ""
	for _, msg := range m.app.GetChatMessages() {
		author := ""
		if msg.Author != "" {
			author += m.renderer.NewStyle().
				Foreground(lipgloss.Color(msg.Color)).
				Render(msg.Author) + ": "
		}
		conversation = lipgloss.JoinVertical(
			lipgloss.Top,
			conversation,
			author+msg.Message,
		)
	}

	conversationHeightAvailable := max(height-lipgloss.Height(textinput)-lipgloss.Height(title)-3, 0)

	wrapped := wrap.String(wordwrap.String(conversation, width-1), width-1)
	split := strings.Split(wrapped, "\n")

	rendered := split[max(len(split)-conversationHeightAvailable, 0):]

	conversation = m.renderer.NewStyle().
		Height(height - lipgloss.Height(textinput) - lipgloss.Height(title) - 1).
		Render("\n" + strings.Join(rendered, "\n"))

	return containerStyle.Render(lipgloss.JoinVertical(lipgloss.Top, title, conversation, textinput))
}
