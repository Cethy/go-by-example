package main

import (
	"context"
	"errors"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/muesli/termenv"
	"net"
	"os"
	"os/signal"
	"ssh-multitodolist/app"
	"ssh-multitodolist/data"
	"ssh-multitodolist/tui/root"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

var (
	host = "0.0.0.0"
	port = "23234"
)

func main() {
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	a := &app.App{}

	r := data.New("./TODO.md", a.NotifyNewData)
	err := r.Init()
	if err != nil {
		panic(err)
	}

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			teaMiddleware(r, a),
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
}

func teaMiddleware(r *data.Repository, app *app.App) wish.Middleware {
	teaHandler := func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		_, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil, nil
		}

		// biggest gotcha working with bubbletea and ssh D:
		renderer := bubbletea.MakeRenderer(s)

		m := root.New(s.User(), r, renderer)

		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}
	programHandler := func(s ssh.Session) *tea.Program {
		m, opts := teaHandler(s)
		if m == nil {
			return nil
		}

		p := tea.NewProgram(m, append(bubbletea.MakeOptions(s), opts...)...)
		app.Programs = append(app.Programs, p)

		return p
	}

	return bubbletea.MiddlewareWithProgramHandler(programHandler, termenv.ANSI256)
}
