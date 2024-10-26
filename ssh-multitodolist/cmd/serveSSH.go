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
	"ssh-multitodolist/app"
	"ssh-multitodolist/app/state"
	"ssh-multitodolist/data"
	"ssh-multitodolist/tui/root"
	"strconv"
	"syscall"
	"time"
)

var serveSSHCmd = &cobra.Command{
	Use:   "server",
	Short: "starts ssh multi-user server",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			host = "0.0.0.0"
			port = "23234"
		)

		if os.Getenv("PORT") != "" {
			port = os.Getenv("PORT")
		}

		a := app.New("Welcome to the chat! 👋")

		r := data.New("./TODO.md", a.NotifyNewData, a.NotifyListRemoved)
		err := r.Init()
		if err != nil {
			panic(err)
		}

		s, err := wish.NewServer(
			wish.WithAddress(net.JoinHostPort(host, port)),
			wish.WithHostKeyPath(".ssh/id_ed25519"),

			// app
			wish.WithMiddleware(
				applicationMiddleware(r, a),
				removeDisconnectedUsersMiddleware(a),
				usernameAlreadyUsedMiddleware(a),
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

// will exit connections if username is already used on some other connection
func usernameAlreadyUsedMiddleware(a *app.App) wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(s ssh.Session) {
			if a.Users[s.User()] == nil {
				next(s)
			}
			fmt.Fprintln(s, "Username \""+s.User()+"\" is already used.")
			s.Exit(1)
		}
	}
}

func removeDisconnectedUsersMiddleware(a *app.App) wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			next(sess)
			a.RemoveUser(sess.User())
		}
	}
}

func applicationMiddleware(repository *data.Repository, application *app.App) wish.Middleware {
	programHandler := func(s ssh.Session) *tea.Program {
		_, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil
		}

		var (
			r  = bubbletea.MakeRenderer(s) // biggest gotcha working with bubbletea and ssh D:
			st = state.New(s.User(), randomColor([]string{}), application.NotifyUserPositionUpdated)
			m  = root.New(st, application, repository, r, false)
			p  = tea.NewProgram(m, append(bubbletea.MakeOptions(s), tea.WithAltScreen())...)
		)
		application.AddUser(p, st)

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
