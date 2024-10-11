package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go-by-example/cli-multitodolist/data"
	"go-by-example/cli-multitodolist/help"
	"go-by-example/cli-multitodolist/keys"
	"go-by-example/cli-multitodolist/statusBar"
	"go-by-example/cli-multitodolist/tabs"
	"go-by-example/cli-multitodolist/todolist"
	"os"
)

type model struct {
	header    tabs.Model
	todolists []todolist.Model
	keys      keys.KeyMap
	help      help.Model
	statusBar statusBar.Model

	width, height      int
	saveOnQuitCallback func(lists []data.NamedList) error
}

func (m model) Init() tea.Cmd {
	return statusBar.NewStatusCmd(m.statusBar.DefaultValue)
}

func (m model) getActiveTodolist() todolist.Model {
	return m.todolists[m.header.ActiveTab]
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if !m.getActiveTodolist().AddListItemActive {
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

	if !m.getActiveTodolist().AddListItemActive {
		m.header, cmd = m.header.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.statusBar, cmd = m.statusBar.Update(msg)
	cmds = append(cmds, cmd)
	m.todolists[m.header.ActiveTab], cmd = m.getActiveTodolist().Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	header := lipgloss.NewStyle().
		Width(m.width).
		Height(4).
		Align(lipgloss.Left, lipgloss.Top).
		Render(m.header.View(m.width))

	helpKeys := append([][]key.Binding{m.keys.ShortHelp()}, m.header.Keys.Help()...)
	helpKeys = append(helpKeys, m.getActiveTodolist().Keys.Help()...)
	helpView := m.help.View(helpKeys)

	statusBarView := m.statusBar.View(m.width)

	outsideContentHeight := lipgloss.Height(header) + lipgloss.Height(helpView) + lipgloss.Height(statusBarView)

	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-outsideContentHeight).
		Align(lipgloss.Left, lipgloss.Top).
		Render(m.getActiveTodolist().View(m.width, m.height-outsideContentHeight))

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
	}

	helpModel := help.New()

	p := tea.NewProgram(model{
		header:    tabs.New(getLabels(namedLists)),
		help:      helpModel,
		keys:      keys.Keys,
		statusBar: statusBar.New(),
		todolists: todolists,
		saveOnQuitCallback: func(lists []data.NamedList) error {
			return data.WriteData(lists, dbFilepath)
		},
	})

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %value", err)
		os.Exit(1)
	}
}
