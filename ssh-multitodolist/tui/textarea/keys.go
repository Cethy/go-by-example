package textarea

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Enter  key.Binding
	Cancel key.Binding
}

var keys = KeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "valid new entry"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel new entry"),
	),
}

func (k KeyMap) Help() [][]key.Binding {
	return [][]key.Binding{{k.Enter, k.Cancel}}
}
