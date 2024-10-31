package room

import (
	"fmt"
	"github.com/charmbracelet/ssh"
	ssh2 "golang.org/x/crypto/ssh"
	"regexp"
	"ssh-multitodolist/app"
	"ssh-multitodolist/data"
)

type Room struct {
	name       string
	Private    bool
	Users      map[string]string
	App        *app.App
	Repository data.Repository
}

func newRoom(name string, private bool, app *app.App, repository data.Repository) *Room {
	return &Room{name, private, make(map[string]string), app, repository}
}

func (r *Room) Join(inviteCode, username string, key ssh.PublicKey) error {
	if r.App.InviteCode == inviteCode {
		r.Users[username] = string(ssh2.MarshalAuthorizedKey(key))
		return nil
	}
	return fmt.Errorf("invalid invite code: %s", r.App.InviteCode)
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

type NoIdentityProvidedError struct{}

func (e NoIdentityProvidedError) Error() string {
	return fmt.Sprint("In order to connect to a private room, you must provide an identity (public key)")
}

func (m *Manager) CreateRoom(roomName string, private bool, userName string, key ssh.PublicKey) (*Room, error) {
	if private && key == nil {
		return nil, NoIdentityProvidedError{}
	}

	a := app.New(roomName, "Welcome to ssh-mutlitodolist! 👋", private)
	r, err := m.repositoryFactory(roomName, a)
	if err != nil {
		return nil, err
	}
	room := newRoom(roomName, private, a, r)

	room.Users[userName] = string(ssh2.MarshalAuthorizedKey(key))

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
