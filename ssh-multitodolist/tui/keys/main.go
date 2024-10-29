package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Help      key.Binding
	Quit      key.Binding
	ForceQuit key.Binding
}

var Keys = KeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q", "quit"),
	),
	ForceQuit: key.NewBinding(
		key.WithKeys("ctrl+c"),
	),
}

func (k KeyMap) HelpKeys() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}
