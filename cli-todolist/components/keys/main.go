package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Help       key.Binding
	Quit       key.Binding
	Check      key.Binding
	AddItem    key.Binding
	RemoveItem key.Binding

	Enter  key.Binding
	Cancel key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.AddItem, k.RemoveItem, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.AddItem, k.RemoveItem, k.Check},
		{k.Up, k.Down, k.Help, k.Quit},
	}
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "save & quit"),
	),
	Check: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("space", "check item "),
	),
	AddItem: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add item"),
	),
	RemoveItem: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "remove item "),
	),

	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "valid new item"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel new item"),
	),
}
