package tabs

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Left       key.Binding
	Right      key.Binding
	AddItem    key.Binding
	EditItem   key.Binding
	RemoveItem key.Binding
}

var keys = KeyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left "),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right   "),
	),
	AddItem: key.NewBinding(
		key.WithKeys("z"),
		key.WithHelp("z", "add list"),
	),
	EditItem: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "edit list title"),
	),
	RemoveItem: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "delete list"),
	),
}

func (k KeyMap) Help() [][]key.Binding {
	return [][]key.Binding{{k.Left, k.Right}, {k.AddItem, k.EditItem, k.RemoveItem}}
}
