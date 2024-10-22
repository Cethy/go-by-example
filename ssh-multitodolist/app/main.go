package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"slices"
)

// Shared state of the application

type User struct {
	Program  *tea.Program
	Username string
}

type App struct {
	Users []*User
}

func (a *App) AddUser(program *tea.Program, username string) {
	a.Users = append(a.Users, &User{Program: program, Username: username})

	a.Notify(UserListUpdatedMsg{})
}

func (a *App) RemoveUser(username string) {
	index := slices.IndexFunc(a.Users, func(user *User) bool {
		return user.Username == username
	})
	if index == -1 {
		return
	}
	a.Users = append(a.Users[:index], a.Users[index+1:]...)

	a.Notify(UserListUpdatedMsg{})
}

// notify

func (a *App) Notify(msg tea.Msg) {
	for _, u := range a.Users {
		go u.Program.Send(msg)
	}
}

type UserListUpdatedMsg struct{}

type NewDataMsg struct{}

func (a *App) NotifyNewData() {
	a.Notify(NewDataMsg{})
}
