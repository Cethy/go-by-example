package todolist

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Check      key.Binding
	AddItem    key.Binding
	RemoveItem key.Binding

	Enter  key.Binding
	Cancel key.Binding
}

var keys = keyMap{
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
