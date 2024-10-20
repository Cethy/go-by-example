package todolist

import (
	"cli-multitodolist/data"
	"cli-multitodolist/input"
	"cli-multitodolist/statusBar"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ListItems      []data.ListItem
	Keys           KeyMap
	Cursor         int // which to-do list item our Cursor is pointing at
	previousCursor int // which to-do list item our Cursor is pointing at (before input is active)
	editCursor     int
}

func New(listItems []data.ListItem) Model {
	return Model{
		ListItems:      listItems,
		Keys:           keys,
		Cursor:         0,
		previousCursor: 0,
		editCursor:     -1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case CreateEntryMsg:
		m.ListItems = append(m.ListItems, data.ListItem{Value: msg.Value, Checked: false})
		m.Cursor = len(m.ListItems) - 1
		cmds = append(cmds, statusBar.NewStatusCmd("New entry added"))
	case CancelCreateEntryMsg:
		m.Cursor = m.previousCursor
		cmds = append(cmds, statusBar.NewStatusCmd("New entry cancelled"))
	case UpdateEntryMsg:
		m.ListItems[m.Cursor] = data.ListItem{Value: msg.Value, Checked: m.ListItems[m.Cursor].Checked}
		m.editCursor = -1
		cmds = append(cmds, statusBar.NewStatusCmd("Entry updated"))
	case CancelUpdateEntryMsg:
		m.editCursor = -1
		cmds = append(cmds, statusBar.NewStatusCmd("Entry update cancelled"))
	case MoveItemUpMsg:
		if msg.cursor > 0 {
			movingUp := m.ListItems[msg.cursor]
			m.ListItems[msg.cursor] = m.ListItems[msg.cursor-1]
			m.ListItems[msg.cursor-1] = movingUp
			m.Cursor = msg.cursor - 1
			cmds = append(cmds, statusBar.NewStatusCmd("Entry moved up"))
		}
	case MoveItemDownMsg:
		if msg.cursor < len(m.ListItems)-1 {
			movingDown := m.ListItems[msg.cursor]
			m.ListItems[msg.cursor] = m.ListItems[msg.cursor+1]
			m.ListItems[msg.cursor+1] = movingDown
			m.Cursor = msg.cursor + 1
			cmds = append(cmds, statusBar.NewStatusCmd("Entry moved down"))
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Up):
			if m.Cursor > 0 {
				m.Cursor--
			}
		case key.Matches(msg, m.Keys.Down):
			if m.Cursor < len(m.ListItems)-1 {
				m.Cursor++
			}
		case key.Matches(msg, m.Keys.Check):
			m.ListItems[m.Cursor].Checked = !m.ListItems[m.Cursor].Checked

			cmds = append(cmds, statusBar.NewStatusCmd("Entry checked"))
		case key.Matches(msg, m.Keys.AddItem):
			cmds = append(cmds, input.NewFocusInputCmd("addEntryInput"))

			m.previousCursor = m.Cursor
			m.Cursor = len(m.ListItems)

			cmds = append(cmds, statusBar.NewPersistingStatusCmd("Typing new entry"))
		case key.Matches(msg, m.Keys.EditItem):
			m.editCursor = m.Cursor
			cmds = append(cmds, input.NewFocusInputValueCmd("editEntryInput", m.ListItems[m.Cursor].Value))
			cmds = append(cmds, statusBar.NewPersistingStatusCmd("Editing entry"))
		case key.Matches(msg, m.Keys.RemoveItem):
			if len(m.ListItems) <= 0 {
				break
			}
			m.ListItems = append(m.ListItems[:m.Cursor], m.ListItems[m.Cursor+1:]...)
			if m.Cursor >= len(m.ListItems)-1 {
				m.Cursor = len(m.ListItems) - 1
			}

			cmds = append(cmds, statusBar.NewStatusCmd("Entry removed"))
		case key.Matches(msg, m.Keys.MoveItemUp):
			cmds = append(cmds, NewMoveItemUpCmd(m.Cursor))
		case key.Matches(msg, m.Keys.MoveItemDown):
			cmds = append(cmds, NewMoveItemDownCmd(m.Cursor))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View(editEntryRender func(t string) string) string {
	content := ""
	for i, listItem := range m.ListItems {
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
