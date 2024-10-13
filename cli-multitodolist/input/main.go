package input

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Keys               KeyMap
	input              textinput.Model
	Active             bool
	getConfirmInputCmd func(value string) tea.Cmd
	getCancelInputCmd  func() tea.Cmd

	id string
}

func NewInput(Placeholder, Prompt string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = Placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Prompt = Prompt

	return ti
}

func New(id string, getConfirmInputCmd func(value string) tea.Cmd, getCancelInputCmd func() tea.Cmd, input textinput.Model) Model {
	return Model{
		Keys:               keys,
		input:              input,
		Active:             false,
		getConfirmInputCmd: getConfirmInputCmd,
		getCancelInputCmd:  getCancelInputCmd,
		id:                 id,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case FocusInputMsg:
		if m.id == msg.id {
			m.Active = true
		}
	case tea.KeyMsg:
		if m.Active {
			switch {
			case key.Matches(msg, m.Keys.Enter):
				cmds = append(cmds, m.getConfirmInputCmd(m.input.Value()))
				m.Active = false
				m.input.SetValue("")
			case key.Matches(msg, m.Keys.Cancel):
				cmds = append(cmds, m.getCancelInputCmd())
				m.Active = false
				m.input.SetValue("")
			}

			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.input.View()
}
