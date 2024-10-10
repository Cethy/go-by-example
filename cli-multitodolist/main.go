package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go-by-example/cli-multitodolist/data"
	"go-by-example/cli-multitodolist/keys"
	"go-by-example/cli-multitodolist/tabs"
	"go-by-example/cli-multitodolist/todolist"
	"os"
)

/*
type newItemUI struct {
	Active bool
	Input  textinput.Model
}

func newNewItemUI() newItemUI {
	// newItem text input
	ti := textinput.New()
	ti.Placeholder = "new item"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Prompt = "  [ ] "

	return newItemUI{
		Active: false,
		Input:  ti,
	}
}

type todoItem struct {
	value   string
	checked bool
}

type todoListUI struct {
	choices []todoItem
	cursor  int // which to-do list item our cursor is pointing at

	keys    keys.keyMap
	help    help.Model
	newItem newItemUI

	width, height      int
	saveOnQuitCallback func(items []todoItem) error
}

func newTodoListUI(choices []todoItem, saveOnQuitCallback func(items []todoItem) error) todoListUI {
	return todoListUI{
		help:    help.New(),
		keys:    keys.keys,
		newItem: newNewItemUI(),

		choices: choices,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		//selected: make(map[int]struct{}),
		saveOnQuitCallback: saveOnQuitCallback,
	}
}

func (m todoListUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.newItem.Active {
		var cmd tea.Cmd
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.Enter):
				m.choices = append(m.choices, todoItem{value: m.newItem.Input.Value(), checked: false})
				m.newItem.Input.SetValue("")
				//m.newItem.Active = false
			case key.Matches(msg, m.keys.Cancel):
				m.newItem.Input.SetValue("")
				m.newItem.Active = false
			}
		}
		m.newItem.Input, cmd = m.newItem.Input.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			// toggle help view
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Quit):
			err := m.saveOnQuitCallback(m.choices)
			if err != nil {
				panic(err)
			}
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Check):
			m.choices[m.cursor].checked = !m.choices[m.cursor].checked
		case key.Matches(msg, m.keys.AddItem):
			m.newItem.Active = true
		case key.Matches(msg, m.keys.RemoveItem):
			m.choices = append(m.choices[:m.cursor], m.choices[m.cursor+1:]...)
		}
	}

	return m, nil
}

func (m todoListUI) View() string {
	header := "What's on your mind today?"

	todolist := ""
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if choice.checked {
			checked = "x" // selected!
		}

		todolist += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.value)
	}

	inputUI := ""
	if m.newItem.Active {
		inputUI = m.newItem.Input.View()
	}

	content := lipgloss.NewStyle().
		//Width(m.width).
		Height(m.height-5).
		Align(lipgloss.Left, lipgloss.Top).
		Render(todolist + inputUI)

	helpView := lipgloss.NewStyle().
		//Width(m.width).
		Height(4).
		Align(lipgloss.Left, lipgloss.Bottom).
		Render(m.help.View(m.keys))

	// Send the UI for rendering
	return lipgloss.JoinVertical(lipgloss.Top, header, content, helpView)
}

func (m todoListUI) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func readData(sourcePath string) []todoItem {
	rawContent, err := os.ReadFile(sourcePath)
	if err != nil {
		panic(err)
	}

	pattern := "\\- \\[(?P<checked> ?x?)\\] (?P<value>[A-z].*)"
	r, _ := regexp.Compile(pattern)
	all := r.FindAllStringSubmatch(string(rawContent), -1)
	var choices []todoItem
	for _, item := range all {
		choices = append(choices, todoItem{
			value:   item[r.SubexpIndex("value")],
			checked: item[r.SubexpIndex("checked")] == "x",
		})
	}

	return choices
}
func writeData(items []todoItem, targetPath string) error {
	content := ""
	for _, item := range items {
		checked := " "
		if item.checked {
			checked = "x"
		}
		content = content + "- [" + checked + "] " + item.value + "\n"
	}

	return os.WriteFile(targetPath, []byte(content), 0644)
}

func main() {
	dbFilepath := "./cli-todolist/TODO.md"

	p := tea.NewProgram(newTodoListUI(
		readData(dbFilepath),
		func(items []todoItem) error {
			return writeData(items, dbFilepath)
		},
	))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %value", err)
		os.Exit(1)
	}
}
*/

type model struct {
	todolists []todolist.Model
	keys      keys.KeyMap
	help      help.Model
	tabs      tabs.Model

	width, height      int
	saveOnQuitCallback func(lists []data.NamedList) error
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) getActiveTodolist() todolist.Model {
	return m.todolists[m.tabs.ActiveTab]
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

			case key.Matches(msg, m.keys.Quit):
				var lists []data.NamedList
				for i, t := range m.todolists {
					lists = append(lists, data.NamedList{
						Name: m.tabs.Tabs[i],
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
		m.tabs, cmd = m.tabs.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.todolists[m.tabs.ActiveTab], cmd = m.getActiveTodolist().Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	header := lipgloss.NewStyle().
		Width(m.width).
		Height(2).
		Align(lipgloss.Left, lipgloss.Bottom).
		Margin(0, 4).
		//Padding(2, 4).
		Render("What's on your mind today ?")

	helpView := lipgloss.NewStyle().
		Width(m.width).
		Height(4).
		Align(lipgloss.Left, lipgloss.Bottom).
		Render(m.help.View(m.keys))

	outsideContentHeight := lipgloss.Height(header) + lipgloss.Height(helpView)

	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-6).
		Align(lipgloss.Left, lipgloss.Top).
		Render(m.tabs.View(func(availableWidth int) string {
			list := m.getActiveTodolist()
			return list.View(availableWidth, m.height-outsideContentHeight-8)
		}))

	return lipgloss.JoinVertical(lipgloss.Top, header, content, helpView)
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

// @todo
// - persist data
// - "add new list" ui
// - update help based on ui
// - edit features
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

	p := tea.NewProgram(model{
		help:      help.New(),
		keys:      keys.Keys,
		tabs:      tabs.New(getLabels(namedLists)),
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
