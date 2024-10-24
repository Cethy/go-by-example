package chat

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

// /
type ChatMessage struct {
	Message string
	Author  string
	Color   string
}

var chatMessages = []ChatMessage{
	{Message: "foo", Author: "Foo", Color: "12"},
	{Message: "BAR!", Author: "Bar", Color: "128"},
	{Message: "FOO!!", Author: "Foo", Color: "12"},
	{Message: "BAR!!!", Author: "Bar", Color: "128"},
	{Message: "WAZZAAAA 🤪", Author: "Baz", Color: "77"},
	{Message: "...", Author: "Foo", Color: "12"},
	{Message: "...", Author: "Bar", Color: "128"},
	{
		Message: "I am the bread of life. He who comes to me will never go hungry, and he who believes in me will never be thirsty. But as I told you, you have seen me, and yet you do not believe. All that the Father gives me will come to me, and whoever comes to me I will never drive away. For I have come down from heaven not to do my own will but to do the will of him who sent me. And this is the will of him who sent me, that I shall lose none of all those he has given me, but raise them up at the last day. For my Father’s will is that everyone who sees the Son and believes in him may have eternal life, and I will raise them up at the last day.",
		Author:  "Jesus✝",
		Color:   "222",
	},
	{Message: "🙏", Author: "Foo", Color: "12"},
	{Message: "🙏", Author: "Bar", Color: "128"},
	{Message: "🙏", Author: "Baz", Color: "77"},
}

///

type Model struct {
	renderer *lipgloss.Renderer
	Keys     KeyMap
}

func New(r *lipgloss.Renderer) Model {
	return Model{
		renderer: r,
		Keys:     keys,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg, isAnyInputActive bool) (Model, tea.Cmd) {
	return m, nil
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
	for _, msg := range chatMessages {
		conversation = lipgloss.JoinVertical(
			lipgloss.Top,
			conversation,
			m.renderer.NewStyle().
				Foreground(lipgloss.Color(msg.Color)).
				Render(msg.Author)+": "+msg.Message,
		)
	}

	wrapped := wrap.String(wordwrap.String(conversation, width-1), width-1)

	conversation = m.renderer.NewStyle().
		Height(height - lipgloss.Height(textinput) - lipgloss.Height(title) - 1).
		Render(wrapped)

	return containerStyle.Render(lipgloss.JoinVertical(lipgloss.Top, title, conversation, textinput))
}
