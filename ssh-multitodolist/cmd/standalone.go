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
	"ssh-multitodolist/data"
	"ssh-multitodolist/tui/root"
)

var standAloneCmd = &cobra.Command{
	Use:   "standalone",
	Short: "starts as a standalone app",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			application = app.New("Welcome to the chat! 👋")
			repository  = data.New("./TODO.md", func() {}, func() {})
			state       = application.NewState("")
			renderer    = lipgloss.DefaultRenderer()
		)

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

		err := repository.Init()
		if err != nil {
			return err
		}

		model := root.New(state, application, repository, renderer, true)

		p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %value", err)
			os.Exit(1)
		}
		return nil
	},
}
