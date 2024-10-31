package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tjarratt/babble"
	"sort"
	"ssh-multitodolist/app/state"
)

type message struct {
	Message string
	Author  string
	Color   string
}

type user struct {
	Program *tea.Program
	State   *state.State
}

type App struct {
	RoomName   string
	IsPrivate  bool
	InviteCode string
	users      map[string]*user
	chat       []message
}

func New(roomName, welcomeMessage string, isPrivate bool) *App {
	babbler := babble.NewBabbler()
	babbler.Count = 3

	return &App{
		RoomName:   roomName,
		IsPrivate:  isPrivate,
		InviteCode: babbler.Babble(),
		users:      make(map[string]*user),
		chat:       []message{{Message: welcomeMessage}},
	}
}

func (a *App) IsUserActive(username string) bool {
	return a.users[username] != nil
}

func (a *App) StatesSorted() []*state.State {
	var states []*state.State
	for _, u := range a.users {
		states = append(states, u.State)
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].Username < states[j].Username
	})

	return states
}

func (a *App) GetUsedColors() []string {
	colors := make([]string, 0)
	for _, u := range a.users {
		colors = append(colors, u.State.Color)
	}
	return colors
}

type UserListUpdatedMsg struct{}

func (a *App) AddUser(program *tea.Program, state *state.State) {
	a.users[state.Username] = &user{
		Program: program,
		State:   state,
	}

	a.Notify(UserListUpdatedMsg{})
}

func (a *App) RemoveUser(username string) {
	delete(a.users, username)

	a.Notify(UserListUpdatedMsg{})
}

func (a *App) AddChatMessage(m, owner, color string) {
	a.chat = append(a.chat, message{m, owner, color})
	a.NotifyNewData()
}
func (a *App) GetChatMessages() []message {
	return a.chat
}

// notify

func (a *App) Notify(msg tea.Msg) {
	for _, u := range a.users {
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
