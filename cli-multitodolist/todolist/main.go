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
	"strings"
)

type Model struct {
	ListItems         []data.ListItem
	cursor            int // which to-do list item our cursor is pointing at
	keys              keyMap
	addListItem       textinput.Model
	AddListItemActive bool

	//saveOnQuitCallback func(items []ListItem) error
}

func NewTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "list item"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Prompt = "  [ ] "

	return ti
}

func New(listItems []data.ListItem /*, saveOnQuitCallback func(items []ListItem) error*/) Model {
	return Model{
		ListItems:         listItems,
		cursor:            0,
		keys:              keys,
		AddListItemActive: false,
		addListItem:       NewTextInput(),
		//saveOnQuitCallback: saveOnQuitCallback,
	}
}

func (m Model) Update(originalMsg tea.Msg) (Model, tea.Cmd) {
	if m.AddListItemActive {
		var cmd tea.Cmd
		switch msg := originalMsg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.Enter):
				m.ListItems = append(m.ListItems, data.ListItem{Value: m.addListItem.Value(), Checked: false})
				m.addListItem.SetValue("")
				m.AddListItemActive = false
				m.cursor = len(m.ListItems) - 1
			case key.Matches(msg, m.keys.Cancel):
				m.addListItem.SetValue("")
				m.AddListItemActive = false
			}
		}
		m.addListItem, cmd = m.addListItem.Update(originalMsg)
		return m, cmd
	}
	switch msg := originalMsg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.ListItems)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Check):
			m.ListItems[m.cursor].Checked = !m.ListItems[m.cursor].Checked
		case key.Matches(msg, m.keys.AddItem):
			m.AddListItemActive = true
			m.cursor = len(m.ListItems)
		case key.Matches(msg, m.keys.RemoveItem):
			if len(m.ListItems) <= 0 {
				break
			}
			m.ListItems = append(m.ListItems[:m.cursor], m.ListItems[m.cursor+1:]...)
			if m.cursor >= len(m.ListItems)-1 {
				m.cursor = len(m.ListItems) - 1
			}
		}
	}

	return m, nil
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
