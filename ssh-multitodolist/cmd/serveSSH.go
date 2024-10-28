package cmd

import (
	"context"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"slices"
	"ssh-multitodolist/app/room"
	"ssh-multitodolist/app/state"
	"ssh-multitodolist/tui/root"
	"strconv"
	"syscall"
	"time"
)

var serveSSHCmd = &cobra.Command{
	Use:   "server",
	Short: "starts ssh multi-user server",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbType, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		var (
			host        = "0.0.0.0"
			port        = "23234"
			roomManager = room.NewManager(getRepositoryFactory(dbType))
		)

		if p, ok := os.LookupEnv("PORT"); ok {
			port = p
		}

		s, err := wish.NewServer(
			wish.WithAddress(net.JoinHostPort(host, port)),
			wish.WithHostKeyPath(".ssh/id_ed25519"),
			// app
			wish.WithMiddleware(
				applicationMiddleware(),
				removeDisconnectedUsersMiddleware,
				usernameAlreadyUsedMiddleware,
				selectRoomMiddleware,
				contextMiddleware(roomManager),
				activeterm.Middleware(),
				logging.Middleware(),
			),
		)
		if err != nil {
			log.Error("Could not start server", "error", err)
		}

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		log.Info("Starting SSH server", "host", host, "port", port)
		go func() {
			if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
				log.Error("Could not start server", "error", err)
				done <- nil
			}
		}()

		<-done
		log.Info("Stopping SSH server")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer func() { cancel() }()
		if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not stop server", "error", err)
		}

		return nil
	},
}

// contextMiddleware adds the room.Manager to the session context.
func contextMiddleware(m *room.Manager) wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(s ssh.Session) {
			ctx := s.Context()
			room.ContextSetManager(ctx, m)

			log.Print(s.Context())
			next(s)
		}
	}
}

// selectRoomMiddleware adds the selected room.Room to the session context.
// it will exit connection if room name isn't valid
// This middleware must be run after the contextMiddleware.
func selectRoomMiddleware(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		ctx := s.Context()
		manager := room.ManagerFromContext(ctx)

		roomName := ""
		if len(s.Command()) > 0 {
			roomName = s.Command()[0]
		}
		roomName, err := room.GetRoomName(roomName)
		if err != nil {
			fmt.Fprintln(s, err)
			s.Exit(1)
		}

		r, err := manager.SelectRoom(roomName)
		if err != nil {
			fmt.Fprintln(s, err)
			s.Exit(1)
		}
		room.ContextSetRoom(ctx, r)

		next(s)
	}
}

// usernameAlreadyUsedMiddleware will exit connection if username is already used in the particular app
// This middleware must be run after the selectRoomMiddleware.
func usernameAlreadyUsedMiddleware(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		r := room.RoomFromContext(s.Context())
		if !r.App.IsUserActive(s.User()) {
			next(s)
		}
		fmt.Fprintln(s, "Username \""+s.User()+"\" is already in use.")
		s.Exit(1)
	}
}

// removeDisconnectedUsersMiddleware handles the removal of user from app when disconnecting
// This middleware must be run after the selectRoomMiddleware.
func removeDisconnectedUsersMiddleware(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		r := room.RoomFromContext(s.Context())
		next(s)
		r.App.RemoveUser(s.User())
	}
}

// applicationMiddleware handles the bubbletea app
// This middleware must be run after the selectRoomMiddleware.
func applicationMiddleware() wish.Middleware {
	programHandler := func(s ssh.Session) *tea.Program {
		ro := room.RoomFromContext(s.Context())

		var (
			r  = bubbletea.MakeRenderer(s) // biggest gotcha working with bubbletea and ssh D:
			st = state.New(s.User(), randomColor(ro.App.GetUsedColors()), ro.App.NotifyUserPositionUpdated)
			m  = root.New(st, ro.App, ro.Repository, r, false)
			p  = tea.NewProgram(m, append(bubbletea.MakeOptions(s), tea.WithAltScreen())...)
		)
		ro.App.AddUser(p, st)

		return p
	}

	return bubbletea.MiddlewareWithProgramHandler(programHandler, termenv.ANSI256)
}

//

func randomColor(alreadyUsedColors []string) string {
	color := strconv.Itoa(rand.Intn(256))
	if slices.Contains(alreadyUsedColors, color) {
		return randomColor(alreadyUsedColors)
	}
	return color
}
