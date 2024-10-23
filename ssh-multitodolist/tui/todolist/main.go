package todolist

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ssh-multitodolist/app"
	"ssh-multitodolist/data"
	"ssh-multitodolist/tui/input"
	"ssh-multitodolist/tui/statusBar"
)

type Model struct {
	state      *app.State
	app        *app.App
	repository *data.Repository
	renderer   *lipgloss.Renderer
	listIndex  int
	Keys       KeyMap
}

func New(state *app.State, application *app.App, repository *data.Repository, renderer *lipgloss.Renderer, listIndex int) Model {
	return Model{
		state:      state,
		app:        application,
		repository: repository,
		renderer:   renderer,
		listIndex:  listIndex,
		Keys:       keys,
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
		m.state.Cursor(msg.Index)
	case CancelCreateEntryMsg:
		m.state.Cursor(m.state.GetPreviousCursor())
		cmds = append(cmds, statusBar.NewStatusCmd("New entry cancelled"))

	case UpdateEntryMsg:
		cmds = append(cmds, data.UpdateEntryCmd(listIndex, m.state.GetCursor(), msg.Value, m.repository))
		m.state.EditCursor(-1)
	case CancelUpdateEntryMsg:
		m.state.EditCursor(-1)
		cmds = append(cmds, statusBar.NewStatusCmd("Updating entry cancelled"))

	case data.EntryRemovedMsg:
		listLen := len(m.repository.Get(listIndex).Items)
		if m.state.GetCursor() >= listLen {
			m.state.Cursor(listLen - 1)
		}

	case MoveItemUpMsg:
		cmds = append(cmds, data.MoveEntryUpCmd(listIndex, m.state.GetCursor(), m.repository))
	case data.EntryMovedUpMsg:
		m.state.Cursor(msg.Index - 1)

	case MoveItemDownMsg:
		cmds = append(cmds, data.MoveEntryDownCmd(listIndex, m.state.GetCursor(), m.repository))
	case data.EntryMovedDownMsg:
		m.state.Cursor(msg.Index + 1)

	case tea.KeyMsg:
		if !isAnyInputActive {
			items := m.repository.Get(listIndex).Items
			switch {
			case key.Matches(msg, m.Keys.Up):
				if m.state.GetCursor() > 0 {
					m.state.Cursor(m.state.GetCursor() - 1)
				}
			case key.Matches(msg, m.Keys.Down):
				if m.state.GetCursor() < len(items)-1 {
					m.state.Cursor(m.state.GetCursor() + 1)
				}
			case key.Matches(msg, m.Keys.Check):
				if len(items) == 0 {
					break
				}
				currentChecked := items[m.state.GetCursor()].Checked
				cmds = append(cmds, data.CheckEntryCmd(listIndex, m.state.GetCursor(), !currentChecked, m.repository))
			case key.Matches(msg, m.Keys.AddItem):
				m.state.PreviousCursor(m.state.GetCursor())
				m.state.Cursor(len(items))
				cmds = append(cmds, input.NewFocusInputCmd("addEntryInput"))
				cmds = append(cmds, statusBar.NewPersistingStatusCmd("Typing new entry"))
			case key.Matches(msg, m.Keys.EditItem):
				items := m.repository.Get(listIndex).Items
				if len(items) == 0 {
					break
				}
				m.state.EditCursor(m.state.GetCursor())
				cmds = append(cmds, input.NewFocusInputValueCmd("editEntryInput", items[m.state.GetCursor()].Value))
				cmds = append(cmds, statusBar.NewPersistingStatusCmd("Editing entry"))
			case key.Matches(msg, m.Keys.RemoveItem):
				if len(items) == 0 {
					break
				}
				cmds = append(cmds, data.RemoveEntryCmd(listIndex, m.state.GetCursor(), m.repository))
			case key.Matches(msg, m.Keys.MoveItemUp):
				cmds = append(cmds, MoveItemUpCmd(m.state.GetCursor()))
			case key.Matches(msg, m.Keys.MoveItemDown):
				cmds = append(cmds, MoveItemDownCmd(m.state.GetCursor()))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) getUpTo2Cursors(index int) [2]string {
	var cursors = [2]string{" ", " "}
	cpt := 0
	for _, s := range m.app.StatesSorted() {
		if s.Username == m.state.Username {
			continue
		}
		if s.GetActiveTab() != m.state.GetActiveTab() {
			continue
		}
		if s.GetCursor() == index {
			cursors[cpt] = m.renderer.NewStyle().
				Foreground(lipgloss.Color(s.Color)).
				Render(">")
			cpt++
		}

		if cpt >= 2 {
			break
		}
	}
	return cursors
}

func (m Model) View(editEntryRender func(t string) string, listIndex int) string {
	items := m.repository.Get(listIndex).Items
	content := ""
	for i, listItem := range items {
		cursors := m.getUpTo2Cursors(i)
		if m.state.GetCursor() == i {
			cursors[1] = ">" // Cursor!
		}
		// Is this choice selected?
		checked := " " // not selected
		if listItem.Checked {
			checked = "x" // selected!
		}

		value := listItem.Value
		if m.state.GetEditCursor() == i {
			value = editEntryRender(value)
		}

		content += fmt.Sprintf(" %s%s [%s] %s", cursors[0], cursors[1], checked, value) + "\n"
	}

	return content
}
