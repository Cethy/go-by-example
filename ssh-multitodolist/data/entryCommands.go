package data

import (
	tea "github.com/charmbracelet/bubbletea"
	"ssh-multitodolist/tui/statusBar"
)

type entryMsg struct {
	ListIndex int
	Index     int
}

type EntryCreatedMsg struct {
	entryMsg
	Value string
}

func CreateEntryCmd(listIndex int, value string, r Repository) tea.Cmd {
	list := r.Get(listIndex)
	list.Items = append(list.Items, ListItem{Value: value, Checked: false})
	index := len(list.Items) - 1
	r.Update(listIndex, list)

	return tea.Batch(statusBar.NewStatusCmd("New entry created"), func() tea.Msg {
		return EntryCreatedMsg{entryMsg: entryMsg{ListIndex: listIndex, Index: index}, Value: value}
	})
}

type EntryUpdatedMsg struct {
	entryMsg
	Value string
}

func UpdateEntryCmd(listIndex int, index int, value string, r Repository) tea.Cmd {
	list := r.Get(listIndex)
	list.Items[index].Value = value
	r.Update(listIndex, list)

	return tea.Batch(statusBar.NewStatusCmd("Entry updated"), func() tea.Msg {
		return EntryUpdatedMsg{entryMsg: entryMsg{ListIndex: listIndex, Index: index}, Value: value}
	})
}

type EntryCheckedMsg struct {
	entryMsg
	Checked bool
}

func CheckEntryCmd(listIndex int, index int, checked bool, r Repository) tea.Cmd {
	list := r.Get(listIndex)
	list.Items[index].Checked = checked
	r.Update(listIndex, list)

	return tea.Batch(statusBar.NewStatusCmd("Entry checked"), func() tea.Msg {
		return EntryCheckedMsg{entryMsg: entryMsg{ListIndex: listIndex, Index: index}, Checked: checked}
	})
}

type EntryRemovedMsg struct {
	entryMsg
}

func RemoveEntryCmd(listIndex int, index int, r Repository) tea.Cmd {
	list := r.Get(listIndex)
	list.Items = append(list.Items[:index], list.Items[index+1:]...)
	r.Update(listIndex, list)

	return tea.Batch(statusBar.NewStatusCmd("Entry removed"), func() tea.Msg {
		return EntryRemovedMsg{entryMsg: entryMsg{ListIndex: listIndex, Index: index}}
	})
}

type EntryMovedUpMsg struct {
	entryMsg
}

func MoveEntryUpCmd(listIndex int, index int, r Repository) tea.Cmd {
	if index <= 0 {
		return nil
	}
	list := r.Get(listIndex)

	movingUp := list.Items[index]
	list.Items[index] = list.Items[index-1]
	list.Items[index-1] = movingUp

	r.Update(listIndex, list)

	return tea.Batch(statusBar.NewStatusCmd("Entry moved up"), func() tea.Msg {
		return EntryMovedUpMsg{entryMsg: entryMsg{ListIndex: listIndex, Index: index}}
	})
}

type EntryMovedDownMsg struct {
	entryMsg
}

func MoveEntryDownCmd(listIndex int, index int, r Repository) tea.Cmd {
	list := r.Get(listIndex)

	if index >= len(list.Items)-1 {
		return nil
	}

	movingDown := list.Items[index]
	list.Items[index] = list.Items[index+1]
	list.Items[index+1] = movingDown

	r.Update(listIndex, list)

	return tea.Batch(statusBar.NewStatusCmd("Entry moved down"), func() tea.Msg {
		return EntryMovedDownMsg{entryMsg: entryMsg{ListIndex: listIndex, Index: index}}
	})
}
