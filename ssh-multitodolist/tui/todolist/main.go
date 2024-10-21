package todolist

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"ssh-multitodolist/data"
	"ssh-multitodolist/tui/input"
	"ssh-multitodolist/tui/statusBar"
)

type Model struct {
	repository     *data.Repository
	listIndex      int
	Keys           KeyMap
	Cursor         int // which to-do list item our Cursor is pointing at
	previousCursor int // which to-do list item our Cursor is pointing at (before input is active)
	editCursor     int
}

func New(repository *data.Repository, listIndex int) Model {
	return Model{
		repository: repository,
		listIndex:  listIndex,
		Keys:       keys,
		editCursor: -1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg, isAnyInputActive bool, listIndex int) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case CreateEntryMsg:
		cmds = append(cmds, data.CreateEntryCmd(listIndex, msg.Value, m.repository))
	case data.EntryCreatedMsg:
		m.Cursor = msg.Index
	case CancelCreateEntryMsg:
		m.Cursor = m.previousCursor
		cmds = append(cmds, statusBar.NewStatusCmd("New entry cancelled"))

	case UpdateEntryMsg:
		cmds = append(cmds, data.UpdateEntryCmd(listIndex, m.Cursor, msg.Value, m.repository))
		m.editCursor = -1
	case CancelUpdateEntryMsg:
		m.editCursor = -1
		cmds = append(cmds, statusBar.NewStatusCmd("Updating entry cancelled"))

	case data.EntryRemovedMsg:
		listLen := len(m.repository.Get(listIndex).Items)
		if m.Cursor >= listLen {
			m.Cursor = listLen - 1
		}

	case MoveItemUpMsg:
		cmds = append(cmds, data.MoveEntryUpCmd(listIndex, m.Cursor, m.repository))
	case data.EntryMovedUpMsg:
		m.Cursor = msg.Index - 1

	case MoveItemDownMsg:
		cmds = append(cmds, data.MoveEntryDownCmd(listIndex, m.Cursor, m.repository))
	case data.EntryMovedDownMsg:
		m.Cursor = msg.Index + 1

	case tea.KeyMsg:
		if !isAnyInputActive {
			items := m.repository.Get(listIndex).Items
			switch {
			case key.Matches(msg, m.Keys.Up):
				if m.Cursor > 0 {
					m.Cursor--
				}
			case key.Matches(msg, m.Keys.Down):
				if m.Cursor < len(items)-1 {
					m.Cursor++
				}
			case key.Matches(msg, m.Keys.Check):
				if len(items) == 0 {
					break
				}
				currentChecked := items[m.Cursor].Checked
				cmds = append(cmds, data.CheckEntryCmd(listIndex, m.Cursor, !currentChecked, m.repository))
			case key.Matches(msg, m.Keys.AddItem):
				m.previousCursor = m.Cursor
				m.Cursor = len(items)
				cmds = append(cmds, input.NewFocusInputCmd("addEntryInput"))
				cmds = append(cmds, statusBar.NewPersistingStatusCmd("Typing new entry"))
			case key.Matches(msg, m.Keys.EditItem):
				items := m.repository.Get(listIndex).Items
				if len(items) == 0 {
					break
				}
				m.editCursor = m.Cursor
				cmds = append(cmds, input.NewFocusInputValueCmd("editEntryInput", items[m.Cursor].Value))
				cmds = append(cmds, statusBar.NewPersistingStatusCmd("Editing entry"))
			case key.Matches(msg, m.Keys.RemoveItem):
				if len(items) == 0 {
					break
				}
				cmds = append(cmds, data.RemoveEntryCmd(listIndex, m.Cursor, m.repository))
			case key.Matches(msg, m.Keys.MoveItemUp):
				cmds = append(cmds, MoveItemUpCmd(m.Cursor))
			case key.Matches(msg, m.Keys.MoveItemDown):
				cmds = append(cmds, MoveItemDownCmd(m.Cursor))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View(editEntryRender func(t string) string, listIndex int) string {
	items := m.repository.Get(listIndex).Items
	content := ""
	for i, listItem := range items {
		// Is the Cursor pointing at this choice?
		cursor := " " // no Cursor
		if m.Cursor == i {
			cursor = ">" // Cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if listItem.Checked {
			checked = "x" // selected!
		}

		value := listItem.Value
		if m.editCursor == i {
			value = editEntryRender(value)
		}

		content += fmt.Sprintf("%s [%s] %s", cursor, checked, value) + "\n"
	}

	return content
}
