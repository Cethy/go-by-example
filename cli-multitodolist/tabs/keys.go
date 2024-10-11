package tabs

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Left  key.Binding
	Right key.Binding
}

var Keys = KeyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left "),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right   "),
	),
}

func (k KeyMap) Help() [][]key.Binding {
	return [][]key.Binding{{k.Left, k.Right}}
}
