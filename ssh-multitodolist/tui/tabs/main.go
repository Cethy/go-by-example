package tabs

import (
	"github.com/charmbracelet/bubbles/key"
	"ssh-multitodolist/data"
	"ssh-multitodolist/tui/input"
	"ssh-multitodolist/tui/statusBar"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	repository   *data.Repository
	renderer     *lipgloss.Renderer
	Keys         KeyMap
	EditTab      int
	ActiveTab    int
	focusConfirm bool
}

func New(repository *data.Repository, renderer *lipgloss.Renderer) Model {
	return Model{
		repository: repository,
		renderer:   renderer,
		Keys:       keys,
		EditTab:    -1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg, isAnyInputActive bool) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case CreateListMsg:
		cmds = append(cmds, data.CreateListCmd(msg.Value, m.repository))
	case data.ListCreatedMsg:
		m.ActiveTab = msg.Index
	case CancelCreateListMsg:
		cmds = append(cmds, statusBar.NewStatusCmd("New list cancelled"))

	case UpdateListMsg:
		cmds = append(cmds, data.UpdateListCmd(m.ActiveTab, msg.Value, m.repository))
		m.EditTab = -1
	case CancelUpdateListMsg:
		m.EditTab = -1
		cmds = append(cmds, statusBar.NewStatusCmd("Updating list cancelled"))

	case ConfirmRemoveListMsg:
		cmds = append(cmds, data.RemoveListCmd(m.ActiveTab, m.repository))
	case data.ListRemovedMsg:
		if m.ActiveTab >= len(m.repository.List())-1 {
			m.ActiveTab = len(m.repository.List()) - 1
		}
	case tea.KeyMsg:
		if m.focusConfirm {
			if msg.String() == "y" || msg.String() == "Y" {
				cmds = append(cmds, ConfirmRemoveListCmd(m.ActiveTab))
			} else {
				cmds = append(cmds, statusBar.NewStatusCmd("Removing list cancelled"))
			}
			m.focusConfirm = false
		} else if !isAnyInputActive {
			switch {
			case key.Matches(msg, m.Keys.Right):
				m.ActiveTab = min(m.ActiveTab+1, len(m.repository.List())-1)
			case key.Matches(msg, m.Keys.Left):
				m.ActiveTab = max(m.ActiveTab-1, 0)
			case key.Matches(msg, m.Keys.AddItem):
				cmds = append(cmds, input.NewFocusInputCmd("addListInput"))
				cmds = append(cmds, statusBar.NewPersistingStatusCmd("Creating new list"))
			case key.Matches(msg, m.Keys.EditItem):
				m.EditTab = m.ActiveTab
				cmds = append(cmds, input.NewFocusInputValueCmd("editListInput", m.repository.GetName(m.ActiveTab)))
				cmds = append(cmds, statusBar.NewPersistingStatusCmd("Editing list title"))
			case key.Matches(msg, m.Keys.RemoveItem):
				cmds = append(cmds, statusBar.NewStatusCmd("Confirm deleting this list ? Y/n"))
				m.focusConfirm = true
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View(extraTabs []string, editTabRender func(t string) string, width int) string {
	tab, activeTab, tabGap := GetStyles(m.renderer)

	var tabs []string
	for i, t := range m.repository.ListNames() {
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
