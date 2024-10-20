package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"ssh-multitodolist/tui/data"
	"ssh-multitodolist/tui/help"
	"ssh-multitodolist/tui/input"
	"ssh-multitodolist/tui/keys"
	"ssh-multitodolist/tui/statusBar"
	"ssh-multitodolist/tui/tabs"
	"ssh-multitodolist/tui/todolist"
	"ssh-multitodolist/tui/viewport"
)

// Model main Model
type Model struct {
	repository     *data.Repository
	renderer       *lipgloss.Renderer
	header         tabs.Model
	todolists      []todolist.Model
	keys           keys.KeyMap
	help           help.Model
	statusBar      statusBar.Model
	viewport       viewport.Model
	addEntryInput  input.Model
	editEntryInput input.Model
	addListInput   input.Model
	editListInput  input.Model

	width, height int
}

func getLabels(tabs []data.NamedList) []string {
	var tabLabels []string
	for _, t := range tabs {
		tabLabels = append(tabLabels, t.Name)
	}
	return tabLabels
}

func New(username string, repository *data.Repository, renderer *lipgloss.Renderer) Model {
	namedLists := repository.List()

	var todolists []todolist.Model
	for _, namedList := range namedLists {
		todolists = append(todolists, todolist.New(namedList.List))
	}

	return Model{
		repository: repository,
		renderer:   renderer,
		header:     tabs.New(getLabels(namedLists), renderer),
		help:       help.New(renderer),
		keys:       keys.Keys,
		statusBar:  statusBar.New(username, renderer),
		todolists:  todolists,
		viewport:   viewport.New(),
		addEntryInput: input.New(
			"addEntryInput",
			todolist.NewCreateEntryCmd,
			todolist.NewCancelCreateEntryCmd,
			input.NewInput("new entry", "  [ ] ", renderer),
		),
		editEntryInput: input.New(
			"editEntryInput",
			todolist.NewUpdateEntryCmd,
			todolist.NewCancelUpdateEntryCmd,
			input.NewInput("edit entry", "", renderer),
		),
		addListInput: input.New(
			"addListInput",
			tabs.NewCreateEntryCmd,
			tabs.NewCancelCreateEntryCmd,
			input.NewInput("new list", "", renderer),
		),
		editListInput: input.New(
			"editListInput",
			tabs.NewUpdateEntryCmd,
			tabs.NewCancelUpdateEntryCmd,
			input.NewInput("edit list", "", renderer),
		),
	}
}

func (m Model) Init() tea.Cmd {
	// @todo handle error
	m.repository.Init()
	return statusBar.NewStatusCmd(m.statusBar.DefaultValue)
}

func (m Model) getActiveTodolist() todolist.Model {
	return m.todolists[m.header.ActiveTab]
}

func (m Model) isAnyInputActive() bool {
	return m.addEntryInput.Active || m.editEntryInput.Active || m.addListInput.Active || m.editListInput.Active
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	log.Printf("(%T) %s\n", msg, msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tabs.CreateEntryMsg:
		m.repository.Create(msg.Value)
		m.todolists = append(m.todolists, todolist.New([]data.ListItem{}))
		cmds = append(cmds, statusBar.NewStatusCmd("New list created"))
	case tabs.CancelCreateEntryMsg:
		cmds = append(cmds, statusBar.NewStatusCmd("New list cancelled"))
	case tabs.ConfirmRemoveEntryMsg:
		m.repository.Delete(msg.Index)
		m.todolists = append(m.todolists[:msg.Index], m.todolists[msg.Index+1:]...)
		cmds = append(cmds, statusBar.NewStatusCmd("list deleted"))
	case tea.KeyMsg:
		if !m.isAnyInputActive() {
			switch {
			case key.Matches(msg, m.keys.Help):
				// toggle help view
				m.help.ShowAll = !m.help.ShowAll
				if m.help.ShowAll {
					cmds = append(cmds, func() tea.Msg {
						return statusBar.NewStatusMsg("🐶 helping")
					})
				}

			case key.Matches(msg, m.keys.Quit):
				for i, t := range m.todolists {
					m.repository.Update(i, data.NamedList{
						Name: m.header.Tabs[i],
						List: t.ListItems,
					})
				}

				err := m.repository.Commit()
				if err != nil {
					panic(err)
				}

				return m, tea.Quit
			}
		}
	}

	if !m.isAnyInputActive() {
		m.header, cmd = m.header.Update(msg)
		cmds = append(cmds, cmd)
		m.todolists[m.header.ActiveTab], cmd = m.getActiveTodolist().Update(msg)
		cmds = append(cmds, cmd)
	}
	m.statusBar, cmd = m.statusBar.Update(msg)
	cmds = append(cmds, cmd)
	m.addEntryInput, cmd = m.addEntryInput.Update(msg)
	cmds = append(cmds, cmd)
	m.editEntryInput, cmd = m.editEntryInput.Update(msg)
	cmds = append(cmds, cmd)
	m.addListInput, cmd = m.addListInput.Update(msg)
	cmds = append(cmds, cmd)
	m.editListInput, cmd = m.editListInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) viewHeader() string {
	newTabView := "+"
	if m.addListInput.Active {
		newTabView = m.addListInput.View()
	}

	return m.renderer.NewStyle().
		Width(m.width).
		Height(4).
		Align(lipgloss.Left, lipgloss.Top).
		Render(m.header.View(
			[]string{newTabView},
			func(t string) string {
				return m.editListInput.View()
			},
			m.width,
		))
}

func (m Model) viewHelp() string {
	helpKeys := [][]key.Binding{m.keys.ShortHelp(), append(m.getActiveTodolist().Keys.HelpDirection(), m.header.Keys.HelpDirection()...)}
	helpKeys = append(helpKeys, m.header.Keys.HelpActions())
	helpKeys = append(helpKeys, m.getActiveTodolist().Keys.HelpActions()...)
	return m.help.View(helpKeys)
}

func (m Model) viewStatusBar() string {
	return m.statusBar.View(m.width)
}

func (m Model) View() string {
	header := m.viewHeader()
	helpView := m.viewHelp()
	statusBarView := m.viewStatusBar()

	outsideContentHeight := lipgloss.Height(header) + lipgloss.Height(helpView) + lipgloss.Height(statusBarView)

	activeTodoList := m.getActiveTodolist()
	addItemInputView := ""
	if m.addEntryInput.Active {
		addItemInputView = m.addEntryInput.View()
	}
	content := m.viewport.View(activeTodoList.View(func(t string) string {
		return m.editEntryInput.View()
	})+addItemInputView, m.width, m.height-outsideContentHeight, activeTodoList.Cursor)

	return lipgloss.JoinVertical(lipgloss.Top, header, content, helpView, statusBarView)
}
