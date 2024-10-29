package chat

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	AddMessage key.Binding
	Enter      key.Binding
	Cancel     key.Binding
}

var keys = KeyMap{
	AddMessage: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "add new message "),
	),

	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "valid new message"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel new message"),
	),
}

func (k KeyMap) HelpActions() []key.Binding {
	return []key.Binding{
		k.AddMessage,
		k.Enter, k.Cancel,
	}
}
