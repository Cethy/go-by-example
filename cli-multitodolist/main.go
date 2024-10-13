package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go-by-example/cli-multitodolist/data"
	"go-by-example/cli-multitodolist/help"
	"go-by-example/cli-multitodolist/input"
	"go-by-example/cli-multitodolist/keys"
	"go-by-example/cli-multitodolist/statusBar"
	"go-by-example/cli-multitodolist/tabs"
	"go-by-example/cli-multitodolist/todolist"
	"go-by-example/cli-multitodolist/viewport"
	"io"
	"log"
	"os"
)

type model struct {
	header        tabs.Model
	todolists     []todolist.Model
	keys          keys.KeyMap
	help          help.Model
	statusBar     statusBar.Model
	viewport      viewport.Model
	addEntryInput input.Model
	addListInput  input.Model

	width, height      int
	saveOnQuitCallback func(lists []data.NamedList) error
}

func (m model) Init() tea.Cmd {
	return statusBar.NewStatusCmd(m.statusBar.DefaultValue)
}

func (m model) getActiveTodolist() todolist.Model {
	return m.todolists[m.header.ActiveTab]
}

func (m model) isAnyInputActive() bool {
	return m.addEntryInput.Active || m.addListInput.Active
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	log.Printf("(%T) %s\n", msg, msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tabs.CreateEntryMsg:
		m.todolists = append(m.todolists, todolist.New([]data.ListItem{}))
		cmds = append(cmds, statusBar.NewStatusCmd("New list created"))
	case tabs.CancelCreateEntryMsg:
		cmds = append(cmds, statusBar.NewStatusCmd("New list cancelled"))
	case tabs.ConfirmRemoveEntryMsg:
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
				var lists []data.NamedList
				for i, t := range m.todolists {
					lists = append(lists, data.NamedList{
						Name: m.header.Tabs[i],
						List: t.ListItems,
					})
				}
				err := m.saveOnQuitCallback(lists)
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
	m.addListInput, cmd = m.addListInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) viewHeader() string {
	newTabView := "+"
	if m.addListInput.Active {
		newTabView = m.addListInput.View()
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Height(4).
		Align(lipgloss.Left, lipgloss.Top).
		Render(m.header.View([]string{newTabView}, m.width))
}

func (m model) viewHelp() string {
	helpKeys := append([][]key.Binding{m.keys.ShortHelp()}, m.header.Keys.Help()...)
	helpKeys = append(helpKeys, m.getActiveTodolist().Keys.Help()...)
	return m.help.View(helpKeys)
}

func (m model) viewStatusBar() string {
	return m.statusBar.View(m.width)
}

func (m model) View() string {
	header := m.viewHeader()
	helpView := m.viewHelp()
	statusBarView := m.viewStatusBar()

	outsideContentHeight := lipgloss.Height(header) + lipgloss.Height(helpView) + lipgloss.Height(statusBarView)

	activeTodoList := m.getActiveTodolist()
	addItemInputView := ""
	if m.addEntryInput.Active {
		addItemInputView = m.addEntryInput.View()
	}
	content := m.viewport.View(activeTodoList.View()+addItemInputView, m.width, m.height-outsideContentHeight, activeTodoList.Cursor)

	return lipgloss.JoinVertical(lipgloss.Top, header, content, helpView, statusBarView)
}

func getLabels(tabs []data.NamedList) []string {
	var tabLabels []string
	for _, t := range tabs {
		label := t.Name
		if label == "" {
			label = "[unnamed list]"
		}
		tabLabels = append(tabLabels, label)
	}
	return tabLabels
}

func main() {
	dbFilepath := "./cli-multitodolist/TODO.md"
	namedLists := data.ReadData(dbFilepath)

	var todolists []todolist.Model
	for _, namedList := range namedLists {
		todolists = append(todolists, todolist.New(namedList.List))
	}

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	} else {
		log.SetOutput(io.Discard)
	}

	helpModel := help.New()

	p := tea.NewProgram(model{
		header:    tabs.New(getLabels(namedLists)),
		help:      helpModel,
		keys:      keys.Keys,
		statusBar: statusBar.New(),
		todolists: todolists,
		viewport:  viewport.New(),
		addEntryInput: input.New(
			"addEntryInput",
			todolist.NewNewEntryCmd,
			todolist.NewCancelNewEntryCmd,
			input.NewInput("new entry", "  [ ] "),
		),
		addListInput: input.New(
			"addListInput",
			tabs.NewCreateEntryCmd,
			tabs.NewCancelCreateEntryCmd,
			input.NewInput("new list", ""),
		),
		saveOnQuitCallback: func(lists []data.NamedList) error {
			return data.WriteData(lists, dbFilepath)
		},
	})

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %value", err)
		os.Exit(1)
	}
}
