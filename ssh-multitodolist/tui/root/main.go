package root

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"ssh-multitodolist/app"
	"ssh-multitodolist/app/state"
	"ssh-multitodolist/data"
	"ssh-multitodolist/tui/chat"
	"ssh-multitodolist/tui/help"
	"ssh-multitodolist/tui/input"
	"ssh-multitodolist/tui/keys"
	"ssh-multitodolist/tui/statusBar"
	"ssh-multitodolist/tui/tabs"
	"ssh-multitodolist/tui/textarea"
	"ssh-multitodolist/tui/todolist"
	"ssh-multitodolist/tui/viewport"
)

// Model main Model
type Model struct {
	state               *state.State
	app                 *app.App
	renderer            *lipgloss.Renderer
	tabs                tabs.Model
	todolist            todolist.Model
	keys                keys.KeyMap
	help                help.Model
	statusBar           statusBar.Model
	viewport            viewport.Model
	chat                chat.Model
	addEntryInput       input.Model
	editEntryInput      input.Model
	addListInput        input.Model
	editListInput       input.Model
	newChatMessageInput textarea.Model
	standalone          bool
	height              int
	chatWidth           int
	mainPanelWidth      int
}

func New(state *state.State, application *app.App, repository data.Repository, renderer *lipgloss.Renderer, standalone bool) Model {
	chatWidth := 40
	if standalone {
		chatWidth = 0
	}

	return Model{
		state:     state,
		app:       application,
		renderer:  renderer,
		tabs:      tabs.New(state, application, repository, renderer),
		help:      help.New(renderer),
		keys:      keys.Keys,
		statusBar: statusBar.New(state.Username, renderer),
		todolist:  todolist.New(state, application, repository, renderer, 0),
		viewport:  viewport.New(),
		chat:      chat.New(state, application, renderer),
		addEntryInput: input.New(
			"addEntryInput",
			todolist.CreateEntryCmd,
			todolist.CancelCreateEntryCmd,
			input.NewInput("new entry", "    [ ] ", renderer),
		),
		editEntryInput: input.New(
			"editEntryInput",
			todolist.UpdateEntryCmd,
			todolist.CancelUpdateEntryCmd,
			input.NewInput("edit entry", "", renderer),
		),
		addListInput: input.New(
			"addListInput",
			tabs.CreateListCmd,
			tabs.CancelCreateListCmd,
			input.NewInput("new list", "", renderer),
		),
		editListInput: input.New(
			"editListInput",
			tabs.UpdateListCmd,
			tabs.CancelUpdateListCmd,
			input.NewInput("edit list", "", renderer),
		),
		newChatMessageInput: textarea.New(
			"newChatMessageInput",
			chat.AddMessageCmd,
			chat.CancelAddMessageCmd,
			textarea.NewInput("Send a message...", "", 6, 280, renderer),
		),
		chatWidth:  chatWidth,
		standalone: standalone,
	}
}

func (m Model) Init() tea.Cmd {
	return statusBar.NewStatusCmd(m.statusBar.DefaultValue)
}

func (m Model) isAnyInputActive() bool {
	return m.addEntryInput.Active ||
		m.editEntryInput.Active ||
		m.addListInput.Active ||
		m.editListInput.Active ||
		m.newChatMessageInput.Active
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	log.Printf("(%T) %s\n", msg, msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.mainPanelWidth = max(msg.Width-m.chatWidth, 0)
		m.height = msg.Height
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.ForceQuit) {
			return m, tea.Quit
		} else if !m.isAnyInputActive() {
			switch {
			case key.Matches(msg, m.keys.ShowInviteCode):
				cmds = append(cmds, statusBar.NewStatusCmd("Invite code: "+m.app.InviteCode))

			case key.Matches(msg, m.keys.Help):
				// toggle help view
				m.help.ShowAll = !m.help.ShowAll
				if m.help.ShowAll {
					cmds = append(cmds, func() tea.Msg {
						return statusBar.NewStatusMsg("🐶 helping")
					})
				}
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			}
		}
	}

	m.tabs, cmd = m.tabs.Update(msg, m.isAnyInputActive())
	cmds = append(cmds, cmd)
	m.todolist, cmd = m.todolist.Update(msg, m.isAnyInputActive(), m.state.GetActiveTab())
	cmds = append(cmds, cmd)
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
	if !m.standalone {
		m.chat, cmd = m.chat.Update(msg, m.isAnyInputActive())
		cmds = append(cmds, cmd)
		m.newChatMessageInput, cmd = m.newChatMessageInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) viewHeader() string {
	newTabView := "+"
	if m.addListInput.Active {
		newTabView = m.addListInput.View()
	}

	return m.renderer.NewStyle().
		Width(m.mainPanelWidth).
		Height(4).
		Align(lipgloss.Left, lipgloss.Top).
		Render(m.tabs.View(
			[]string{newTabView},
			m.app.RoomName,
			m.app.IsPrivate,
			func(t string) string {
				return m.editListInput.View()
			},
			m.mainPanelWidth,
		))
}

func (m Model) viewHelp() string {
	helpKeys := []help.KeysForView{
		{"global", m.keys.HelpKeys()},
		{"navigation", append(m.todolist.Keys.HelpDirection(), m.tabs.Keys.HelpDirection()...)},
		{"actions", m.tabs.Keys.HelpActions()},
		{"more actions", m.todolist.Keys.HelpActions()},
	}
	if !m.standalone {
		helpKeys = append(helpKeys, help.KeysForView{Title: "chat", Keys: m.chat.Keys.HelpActions()})
	}
	return m.help.View(helpKeys)
}

func (m Model) viewStatusBar() string {
	return m.statusBar.View(m.mainPanelWidth + m.chatWidth)
}

func (m Model) viewChat(height int) string {
	return m.chat.View(func() string {
		return m.newChatMessageInput.View(m.chatWidth)
	}, m.chatWidth, height)
}

func (m Model) viewConnectedUsers() string {
	if m.standalone {
		return ""
	}
	connectedUsers := "Connected users: " + m.state.Username
	for _, s := range m.app.StatesSorted() {
		if s.Username == m.state.Username {
			continue
		}
		connectedUsers += ", " + m.renderer.NewStyle().
			Foreground(lipgloss.Color(s.Color)).
			Render(s.Username)
	}
	return connectedUsers
}

func (m Model) viewContent(outsideContentHeight int) string {
	addItemInputView := ""
	if m.addEntryInput.Active {
		addItemInputView = m.addEntryInput.View()
	}

	return m.viewport.View(
		m.todolist.View(
			func(t string) string {
				return m.editEntryInput.View()
			},
			m.state.GetActiveTab(),
		)+addItemInputView,
		m.mainPanelWidth,
		m.height-outsideContentHeight,
		m.state.GetCursor(),
	)
}

func (m Model) View() string {
	var (
		tabsView           = m.viewHeader()
		helpView           = m.viewHelp()
		statusBarView      = m.viewStatusBar()
		chatView           = m.viewChat(m.height - lipgloss.Height(statusBarView))
		connectedUsersView = m.viewConnectedUsers()
	)

	outsideContentHeight := lipgloss.Height(tabsView) +
		lipgloss.Height(helpView) +
		lipgloss.Height(statusBarView) +
		lipgloss.Height(connectedUsersView)

	content := m.viewContent(outsideContentHeight)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.JoinVertical(lipgloss.Top, tabsView, content, helpView, connectedUsersView),
			chatView,
		),
		statusBarView,
	)
}
