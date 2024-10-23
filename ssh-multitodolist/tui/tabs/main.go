package tabs

import (
	"github.com/charmbracelet/bubbles/key"
	"ssh-multitodolist/app"
	"ssh-multitodolist/data"
	"ssh-multitodolist/tui/input"
	"ssh-multitodolist/tui/statusBar"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	state      *app.State
	app        *app.App
	repository *data.Repository
	renderer   *lipgloss.Renderer
	Keys       KeyMap
}

func New(state *app.State, application *app.App, repository *data.Repository, renderer *lipgloss.Renderer) Model {
	return Model{
		state:      state,
		app:        application,
		repository: repository,
		renderer:   renderer,
		Keys:       keys,
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
		m.state.ActiveTab(msg.Index)
	case CancelCreateListMsg:
		cmds = append(cmds, statusBar.NewStatusCmd("New list cancelled"))

	case UpdateListMsg:
		cmds = append(cmds, data.UpdateListCmd(m.state.GetActiveTab(), msg.Value, m.repository))
		m.state.EditTab(-1)
	case CancelUpdateListMsg:
		m.state.EditTab(-1)
		cmds = append(cmds, statusBar.NewStatusCmd("Updating list cancelled"))

	case ConfirmRemoveListMsg:
		cmds = append(cmds, data.RemoveListCmd(m.state.GetActiveTab(), m.repository))
	case data.ListRemovedMsg, app.UserPositionUpdatedMsg:
		if m.state.GetActiveTab() >= len(m.repository.List())-1 {
			m.state.ActiveTab(len(m.repository.List()) - 1)
		}
	case app.ListRemovedMsg:
		// if somebody else removed a tab before user could confirm
		if m.state.GetRemovingTab() != -1 {
			cmds = append(cmds, statusBar.NewStatusCmd("Removing list cancelled"))
			m.state.RemovingTab(-1)
		}
	case tea.KeyMsg:
		if m.state.GetRemovingTab() >= 0 {
			if msg.String() == "y" || msg.String() == "Y" {
				cmds = append(cmds, ConfirmRemoveListCmd(m.state.GetRemovingTab()))
			} else {
				cmds = append(cmds, statusBar.NewStatusCmd("Removing list cancelled"))
			}
			m.state.RemovingTab(-1)
		} else if !isAnyInputActive {
			switch {
			case key.Matches(msg, m.Keys.Right):
				m.state.ActiveTab(min(m.state.GetActiveTab()+1, len(m.repository.List())-1))
			case key.Matches(msg, m.Keys.Left):
				m.state.ActiveTab(max(m.state.GetActiveTab()-1, 0))
			case key.Matches(msg, m.Keys.AddItem):
				// @todo add state and notify ?
				cmds = append(cmds, input.NewFocusInputCmd("addListInput"))
				cmds = append(cmds, statusBar.NewPersistingStatusCmd("Creating new list"))
			case key.Matches(msg, m.Keys.EditItem):
				m.state.EditTab(m.state.GetActiveTab())
				cmds = append(cmds, input.NewFocusInputValueCmd("editListInput", m.repository.GetName(m.state.GetActiveTab())))
				cmds = append(cmds, statusBar.NewPersistingStatusCmd("Editing list title"))
			case key.Matches(msg, m.Keys.RemoveItem):
				// keep one list at all time
				if len(m.repository.List()) <= 1 {
					break
				}
				cmds = append(cmds, statusBar.NewStatusCmd("Confirm deleting this list ? Y/n"))
				m.state.RemovingTab(m.state.GetActiveTab())
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View(extraTabs []string, editTabRender func(t string) string, width int) string {
	tab, activeTab, tabGap := GetStyles(m.renderer)

	var tabs []string
	for i, t := range m.repository.ListNames() {
		if m.state.GetEditTab() == i {
			tabs = append(tabs, activeTab.Render(editTabRender(t)))
			continue
		}

		if t == "" {
			t = "[unnamed list]"
		}

		for _, s := range m.app.StatesSorted() {
			if s.Username == m.state.Username {
				continue
			}
			if s.GetActiveTab() == i {
				t = m.renderer.NewStyle().
					Foreground(lipgloss.Color(s.Color)).
					Render(">") + t
			}
		}

		if m.state.GetActiveTab() == i {
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
