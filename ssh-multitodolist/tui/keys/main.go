package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Help           key.Binding
	Quit           key.Binding
	ForceQuit      key.Binding
	ShowInviteCode key.Binding
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
	ShowInviteCode: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "show invite code"),
	),
}

func (k KeyMap) HelpKeys() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.ShowInviteCode}
}
