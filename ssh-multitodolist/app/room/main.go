package room

import (
	"fmt"
	"regexp"
	"ssh-multitodolist/app"
	"ssh-multitodolist/data"
)

type Room struct {
	name       string
	Private    bool
	App        *app.App
	Repository data.Repository
}

func newRoom(name string, private bool, app *app.App, repository data.Repository) *Room {
	return &Room{name, private, app, repository}
}

type Manager struct {
	rooms             []*Room
	repositoryFactory func(roomName string, app *app.App) (data.Repository, error)
}

func NewManager(rf func(roomName string, app *app.App) (data.Repository, error)) *Manager {
	return &Manager{make([]*Room, 0), rf}
}

type NotFoundError struct {
	roomName string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("room \"%s\" not found", e.roomName)
}

type ConflictingPrivateError struct{}

func (e ConflictingPrivateError) Error() string {
	return fmt.Sprint("cannot access this room with this privacy level")
}

func (m *Manager) SelectRoom(roomName string, private bool) (*Room, error) {
	if len(m.rooms) > 0 {
		for _, room := range m.rooms {
			if room.name == roomName {
				if private != room.Private {
					return nil, ConflictingPrivateError{}
				}
				return room, nil
			}
		}
	}
	return nil, NotFoundError{roomName}
}

func (m *Manager) CreateRoom(roomName string, private bool) (*Room, error) {
	a := app.New(roomName, "Welcome to ssh-mutlitodolist! 👋", private)
	r, err := m.repositoryFactory(roomName, a)
	if err != nil {
		return nil, err
	}
	room := newRoom(roomName, private, a, r)

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
