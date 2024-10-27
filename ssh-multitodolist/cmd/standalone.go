package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"ssh-multitodolist/app"
	"ssh-multitodolist/app/room"
	"ssh-multitodolist/app/state"
	"ssh-multitodolist/data/file"
	"ssh-multitodolist/tui/root"
)

var standAloneCmd = &cobra.Command{
	Use:   "standalone [room]",
	Short: "starts as a standalone app",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(os.Getenv("DEBUG")) > 0 {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				fmt.Println("fatal:", err)
				os.Exit(1)
			}
			defer f.Close()
		} else {
			log.SetOutput(io.Discard)
		}

		roomName := ""
		if len(args) > 0 {
			roomName = args[0]
		}
		roomName, err := room.GetRoomName(roomName)
		if err != nil {
			return err
		}

		repository := file.New(roomName, func() {}, func() {})
		err = repository.Init()
		if err != nil {
			return err
		}

		var (
			application = app.New("")
			r           = lipgloss.DefaultRenderer()
			st          = state.New("", "255", application.NotifyUserPositionUpdated)
			m           = root.New(st, application, repository, r, true)
			p           = tea.NewProgram(m, tea.WithAltScreen())
		)

		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %value", err)
			os.Exit(1)
		}
		return nil
	},
}
