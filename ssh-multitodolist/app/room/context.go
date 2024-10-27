package room

import (
	"github.com/charmbracelet/ssh"
)

// contextKeyManager is the room Manager context key.
var contextKeyManager = &struct{ string }{"roomManager"}

// ManagerFromContext returns the Manager from the given context.
func ManagerFromContext(ctx ssh.Context) *Manager {
	if m, ok := ctx.Value(contextKeyManager).(*Manager); ok {
		return m
	}

	return nil
}

// ContextSetManager sets the given Manager to context.
func ContextSetManager(ctx ssh.Context, m *Manager) {
	ctx.SetValue(contextKeyManager, m)
	//return context.Set(ctx, contextKeyManager, m)
}

// ContextKeyRoom is the Room context key.
var ContextKeyRoom = &struct{ string }{"room"}

// RoomFromContext returns the Room from the given context.
func RoomFromContext(ctx ssh.Context) *Room {
	if r, ok := ctx.Value(ContextKeyRoom).(*Room); ok {
		return r
	}

	return nil
}

// ContextSetRoom sets the given Room to context.
func ContextSetRoom(ctx ssh.Context, r *Room) {
	ctx.SetValue(ContextKeyRoom, r)
	//return context.WithValue(ctx, ContextKeyRoom, r)
}
