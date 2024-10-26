package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"sort"
	"ssh-multitodolist/app/state"
)

type Message struct {
	Message string
	Author  string
	Color   string
}

type user struct {
	Program *tea.Program
	State   *state.State
}

type App struct {
	Users map[string]*user
	chat  []Message
}

func New(welcomeMessage string) *App {
	return &App{
		Users: make(map[string]*user),
		chat:  []Message{{Message: welcomeMessage}},
	}
}

func (a *App) StatesSorted() []*state.State {
	var states []*state.State
	for _, u := range a.Users {
		states = append(states, u.State)
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].Username < states[j].Username
	})

	return states
}

type UserListUpdatedMsg struct{}

func (a *App) AddUser(program *tea.Program, state *state.State) {
	a.Users[state.Username] = &user{
		Program: program,
		State:   state,
	}

	a.Notify(UserListUpdatedMsg{})
}

func (a *App) RemoveUser(username string) {
	delete(a.Users, username)

	a.Notify(UserListUpdatedMsg{})
}

func (a *App) AddChatMessage(m, owner, color string) {
	a.chat = append(a.chat, Message{m, owner, color})
	a.NotifyNewData()
}
func (a *App) GetChatMessages() []Message {
	return a.chat
}

// notify

func (a *App) Notify(msg tea.Msg) {
	for _, u := range a.Users {
		go u.Program.Send(msg)
	}
}

type NewDataMsg struct{}

func (a *App) NotifyNewData() {
	a.Notify(NewDataMsg{})
}

type ListRemovedMsg struct{}

func (a *App) NotifyListRemoved() {
	a.Notify(ListRemovedMsg{})
}

type UserPositionUpdatedMsg struct{}

func (a *App) NotifyUserPositionUpdated() {
	a.Notify(UserPositionUpdatedMsg{})
}

//
