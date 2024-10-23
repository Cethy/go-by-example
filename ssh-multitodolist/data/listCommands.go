package data

import (
	tea "github.com/charmbracelet/bubbletea"
	"ssh-multitodolist/tui/statusBar"
)

type ListMsg struct {
	Index int
}

type ListCreatedMsg struct {
	ListMsg
	Name string
}

func CreateListCmd(name string, r *Repository) tea.Cmd {
	index := r.Create(name)
	return tea.Batch(statusBar.NewStatusCmd("New list created"), func() tea.Msg {
		return ListCreatedMsg{ListMsg: ListMsg{Index: index}, Name: name}
	})
}

type ListUpdatedMsg struct {
	ListMsg
	Name string
}

func UpdateListCmd(index int, name string, r *Repository) tea.Cmd {
	r.UpdateName(index, name)
	return func() tea.Msg {
		return ListUpdatedMsg{ListMsg: ListMsg{Index: index}, Name: name}
	}
}

type ListRemovedMsg struct {
	ListMsg
}

func RemoveListCmd(index int, r *Repository) tea.Cmd {
	// keep one list at all time
	if len(r.List()) <= 1 {
		return nil
	}

	r.Delete(index)
	return func() tea.Msg {
		return ListRemovedMsg{ListMsg: ListMsg{Index: index}}
	}
}
