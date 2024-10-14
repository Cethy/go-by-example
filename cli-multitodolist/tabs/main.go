package tabs

import (
	"github.com/charmbracelet/bubbles/key"
	"go-by-example/cli-multitodolist/input"
	"go-by-example/cli-multitodolist/statusBar"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Keys         KeyMap
	Tabs         []string
	EditTab      int
	ActiveTab    int
	focusConfirm bool
}

func New(tabs []string) Model {
	return Model{
		Keys:    keys,
		Tabs:    tabs,
		EditTab: -1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case CreateEntryMsg:
		m.Tabs = append(m.Tabs, msg.Value)
		m.ActiveTab = len(m.Tabs) - 1
	case UpdateEntryMsg:
		m.Tabs[m.ActiveTab] = msg.Value
		m.EditTab = -1
	case CancelUpdateEntryMsg:
		m.EditTab = -1
	case RemoveEntryMsg:
		cmds = append(cmds, statusBar.NewStatusCmd("Confirm deleting this list ? Y/n"))
		m.focusConfirm = true
	case ConfirmRemoveEntryMsg:
		m.Tabs = append(m.Tabs[:msg.Index], m.Tabs[msg.Index+1:]...)
		if m.ActiveTab >= len(m.Tabs) {
			m.ActiveTab = len(m.Tabs) - 1
		}

	case tea.KeyMsg:
		if m.focusConfirm {
			if msg.String() == "y" || msg.String() == "Y" {
				cmds = append(cmds, NewConfirmRemoveEntryCmd(m.ActiveTab))
			}
			m.focusConfirm = false
		} else {
			switch {
			case key.Matches(msg, m.Keys.Right):
				m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
				return m, nil
			case key.Matches(msg, m.Keys.Left):
				m.ActiveTab = max(m.ActiveTab-1, 0)
				return m, nil
			case key.Matches(msg, m.Keys.AddItem):
				cmds = append(cmds, input.NewFocusInputCmd("addListInput"))
				cmds = append(cmds, statusBar.NewPersistingStatusCmd("Creating new list"))
			case key.Matches(msg, m.Keys.EditItem):
				m.EditTab = m.ActiveTab
				cmds = append(cmds, input.NewFocusInputValueCmd("editListInput", m.Tabs[m.ActiveTab]))
				cmds = append(cmds, statusBar.NewPersistingStatusCmd("Editing list title"))
			case key.Matches(msg, m.Keys.RemoveItem):
				cmds = append(cmds, NewRemoveEntryCmd(m.ActiveTab))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View(extraTabs []string, editTabRender func(t string) string, width int) string {
	var tabs []string
	for i, t := range m.Tabs {
		if m.EditTab == i {
			tabs = append(tabs, activeTab.Render(editTabRender(t)))
			continue
		}

		if t == "" {
			t = "[unnamed list]"
		}
		if m.ActiveTab == i {
			tabs = append(tabs, activeTab.Render(t))
		} else {
			tabs = append(tabs, tab.Render(t))
		}
	}
	for _, t := range extraTabs {
		tabs = append(tabs, tab.Render(t))
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tabs...,
	)
	gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(header)-2)))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, header, gap)
}
