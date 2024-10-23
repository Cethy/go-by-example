package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"math/rand"
	"slices"
	"sort"
	"strconv"
)

// Shared state of the application

type State struct {
	Username    string
	Color       string
	editTab     int
	activeTab   int
	removingTab int
	notify      func()
}

func (s *State) EditTab(v int) {
	oldV := s.editTab
	s.editTab = v

	if oldV != s.editTab {
		s.notify()
	}
}
func (s *State) GetEditTab() int {
	return s.editTab
}
func (s *State) ActiveTab(v int) {
	oldV := s.activeTab
	s.activeTab = v

	if oldV != s.activeTab {
		s.notify()
	}
}
func (s *State) GetActiveTab() int {
	return s.activeTab
}
func (s *State) RemovingTab(v int) {
	oldV := s.removingTab
	s.removingTab = v

	if oldV != s.removingTab {
		s.notify()
	}
}
func (s *State) GetRemovingTab() int {
	return s.removingTab
}

type User struct {
	Program *tea.Program
	State   *State
}

type App struct {
	Users map[string]*User
}

func New() *App {
	return &App{
		Users: make(map[string]*User),
	}
}

func (a *App) NewState(username string) *State {
	return &State{
		Username:    username,
		Color:       randomColor([]string{}),
		editTab:     -1,
		removingTab: -1,
		notify:      a.NotifyUserPositionUpdated,
	}
}

func (a *App) StatesSorted() []*State {
	var states []*State
	for _, u := range a.Users {
		states = append(states, u.State)
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].Username < states[j].Username
	})

	return states
}

type UserListUpdatedMsg struct{}

func (a *App) AddUser(program *tea.Program, state *State) {
	a.Users[state.Username] = &User{
		Program: program,
		State:   state,
	}

	a.Notify(UserListUpdatedMsg{})
}

func (a *App) RemoveUser(username string) {
	delete(a.Users, username)

	a.Notify(UserListUpdatedMsg{})
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

func randomColor(alreadyUsedColors []string) string {
	color := strconv.Itoa(rand.Intn(256))
	if slices.Contains(alreadyUsedColors, color) {
		return randomColor(alreadyUsedColors)
	}
	return color
}
