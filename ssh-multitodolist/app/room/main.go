package room

import (
	"fmt"
	"log"
	"regexp"
	"ssh-multitodolist/app"
	"ssh-multitodolist/data"
)

type Room struct {
	name       string
	App        *app.App
	Repository data.Repository
}

func newRoom(name string, app *app.App, repository data.Repository) *Room {
	return &Room{name, app, repository}
}

type Manager struct {
	rooms             []*Room
	repositoryFactory func(roomName string, app *app.App) data.Repository
}

func NewManager(rf func(roomName string, app *app.App) data.Repository) *Manager {
	return &Manager{make([]*Room, 0), rf}
}

func (m *Manager) SelectRoom(roomName string) (*Room, error) {
	log.Println(m)
	if len(m.rooms) > 0 {
		for _, room := range m.rooms {
			if room.name == roomName {
				return room, nil
			}
		}
	}

	a := app.New("Welcome to ssh-mutlitodolist! 👋")
	r := m.repositoryFactory(roomName, a)
	err := r.Init()
	if err != nil {
		return nil, err
	}
	room := newRoom(roomName, a, r)

	m.rooms = append(m.rooms, room)
	return room, nil
}

//

func GetRoomName(roomName string) (string, error) {
	if len(roomName) == 0 {
		return "TODO", nil
	}
	if !isValidRoomName(roomName) {
		return "", fmt.Errorf("invalid ROOM name: %s", roomName)
	}
	return roomName, nil
}

func isValidRoomName(roomName string) bool {
	// Don't allow (.) in order to prevent the creation of files in parent directories
	re := regexp.MustCompile(`^[a-zA-Z0-9\-_/]+$`)
	return re.MatchString(roomName)
}
