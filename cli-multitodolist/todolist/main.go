package todolist

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
	"go-by-example/cli-multitodolist/data"
	"go-by-example/cli-multitodolist/statusBar"
	"strings"
)

type Model struct {
	ListItems         []data.ListItem
	cursor            int // which to-do list item our cursor is pointing at
	Keys              KeyMap
	addListItem       textinput.Model
	AddListItemActive bool
	previousCursor    int // which to-do list item our cursor is pointing at (before input is active)
}

func NewTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "type new entry"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Prompt = "  [ ] "

	return ti
}

func New(listItems []data.ListItem) Model {
	return Model{
		ListItems:         listItems,
		cursor:            0,
		previousCursor:    0,
		Keys:              keys,
		AddListItemActive: false,
		addListItem:       NewTextInput(),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.AddListItemActive {
		var cmd tea.Cmd
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.Keys.Enter):
				m.ListItems = append(m.ListItems, data.ListItem{Value: m.addListItem.Value(), Checked: false})
				m.addListItem.SetValue("")
				m.AddListItemActive = false
				m.cursor = len(m.ListItems) - 1

				cmds = append(cmds, statusBar.NewStatusCmd("New entry added"))
			case key.Matches(msg, m.Keys.Cancel):
				m.addListItem.SetValue("")
				m.AddListItemActive = false
				m.cursor = m.previousCursor

				cmds = append(cmds, statusBar.NewStatusCmd("New entry cancelled"))
			}
		}
		m.addListItem, cmd = m.addListItem.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.Keys.Down):
			if m.cursor < len(m.ListItems)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.Keys.Check):
			m.ListItems[m.cursor].Checked = !m.ListItems[m.cursor].Checked

			cmds = append(cmds, statusBar.NewStatusCmd("Entry checked"))
		case key.Matches(msg, m.Keys.AddItem):
			m.AddListItemActive = true

			m.previousCursor = m.cursor
			m.cursor = len(m.ListItems)

			cmds = append(cmds, statusBar.NewPersistingStatusCmd("Entry Adding"))
		case key.Matches(msg, m.Keys.RemoveItem):
			if len(m.ListItems) <= 0 {
				break
			}
			m.ListItems = append(m.ListItems[:m.cursor], m.ListItems[m.cursor+1:]...)
			if m.cursor >= len(m.ListItems)-1 {
				m.cursor = len(m.ListItems) - 1
			}

			cmds = append(cmds, statusBar.NewStatusCmd("Entry removed"))
		}
	}

	return m, tea.Batch(cmds...)
}

func listItemView(listItem data.ListItem, withCursor bool) string {
	// Is the cursor pointing at this choice?
	cursor := " " // no cursor
	if withCursor {
		cursor = ">" // cursor!
	}

	// Is this choice selected?
	checked := " " // not selected
	if listItem.Checked {
		checked = "x" // selected!
	}

	return fmt.Sprintf("%s [%s] %s", cursor, checked, listItem.Value)
}

func (m Model) View(maxWidth, maxHeight int) string {
	content := ""
	for i, listItem := range m.ListItems {
		content += listItemView(listItem, m.cursor == i) + "\n"
	}

	if m.AddListItemActive {
		content += m.addListItem.View()
	}

	// Send the UI for rendering
	return m.viewport(content, maxWidth, maxHeight)
}

// @todo remember previous slice position (so when cursor goes down and then up, page doesn't scroll before cursor is at the top)
// limit the width & height of the viewable content
func (m Model) viewport(content string, maxWidth int, maxHeight int) string {
	// @todo properly init the app (first render doesn't have the width/height values set)
	if maxHeight < 0 {
		return content
	}

	wrapped := wrap.String(wordwrap.String(content, maxWidth), maxWidth)

	lines := strings.Split(wrapped, "\n")
	if len(lines) < maxHeight {
		return lipgloss.NewStyle().
			Width(maxWidth).
			Height(maxHeight).
			Align(lipgloss.Left, lipgloss.Top).
			Render(strings.Join(lines, "\n"))
	}

	// make sure the cursor is always visible
	if m.cursor+1 >= maxHeight {
		return "↑\n" + strings.Join(lines[m.cursor+3-maxHeight:m.cursor+2-1], "\n") + m.viewportTail(len(lines), maxWidth)
	}
	return strings.Join(lines[0:maxHeight-1], "\n") + m.viewportTail(len(lines), maxWidth)
}

func (m Model) viewportTail(contentLen, maxWidth int) string {
	tail := " "
	if m.cursor+2 < contentLen {
		tail = "↓"
	}
	currentPos := fmt.Sprintf("%d/%d", min(m.cursor+1, contentLen-1), contentLen-1)

	return "\n" + lipgloss.JoinHorizontal(lipgloss.Top, tail, lipgloss.NewStyle().Width(maxWidth-1-len(currentPos)).Render(""), currentPos)
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
