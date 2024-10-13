package todolist

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"go-by-example/cli-multitodolist/data"
	"go-by-example/cli-multitodolist/input"
	"go-by-example/cli-multitodolist/statusBar"
)

type Model struct {
	ListItems      []data.ListItem
	Keys           KeyMap
	Cursor         int // which to-do list item our Cursor is pointing at
	previousCursor int // which to-do list item our Cursor is pointing at (before input is active)
}

func New(listItems []data.ListItem) Model {
	return Model{
		ListItems:      listItems,
		Keys:           keys,
		Cursor:         0,
		previousCursor: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case NewEntryMsg:
		m.ListItems = append(m.ListItems, data.ListItem{Value: msg.Value, Checked: false})
		m.Cursor = len(m.ListItems) - 1
		cmds = append(cmds, statusBar.NewStatusCmd("New entry added"))
	case CancelNewEntryMsg:
		m.Cursor = m.previousCursor
		cmds = append(cmds, statusBar.NewStatusCmd("New entry cancelled"))
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
		case key.Matches(msg, m.Keys.RemoveItem):
			if len(m.ListItems) <= 0 {
				break
			}
			m.ListItems = append(m.ListItems[:m.Cursor], m.ListItems[m.Cursor+1:]...)
			if m.Cursor >= len(m.ListItems)-1 {
				m.Cursor = len(m.ListItems) - 1
			}

			cmds = append(cmds, statusBar.NewStatusCmd("Entry removed"))
		}
	}

	return m, tea.Batch(cmds...)
}

func listItemView(listItem data.ListItem, withCursor bool) string {
	// Is the Cursor pointing at this choice?
	cursor := " " // no Cursor
	if withCursor {
		cursor = ">" // Cursor!
	}

	// Is this choice selected?
	checked := " " // not selected
	if listItem.Checked {
		checked = "x" // selected!
	}

	return fmt.Sprintf("%s [%s] %s", cursor, checked, listItem.Value)
}

func (m Model) View() string {
	content := ""
	for i, listItem := range m.ListItems {
		content += listItemView(listItem, m.Cursor == i) + "\n"
	}

	return content
}
