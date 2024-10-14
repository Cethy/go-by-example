package todolist

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Check      key.Binding
	AddItem    key.Binding
	EditItem   key.Binding
	RemoveItem key.Binding

	Enter  key.Binding
	Cancel key.Binding
}

var keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Check: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("space", "check entry "),
	),
	AddItem: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add new entry "),
	),
	EditItem: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "edit entry"),
	),
	RemoveItem: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "remove entry "),
	),

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
	return [][]key.Binding{
		{k.Up, k.Down, k.Check, k.AddItem, k.EditItem, k.RemoveItem},
		{k.Enter, k.Cancel},
	}
}
