package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"ssh-multitodolist/app"
	"ssh-multitodolist/data"
	"ssh-multitodolist/data/storage"
)

var rootCmd = &cobra.Command{
	Use:   "multi-todolist",
	Short: "SSH-enabled multi-user todolist",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	serveSSHCmd.Flags().String("db", "file", "db type, file or redis")
	rootCmd.AddCommand(serveSSHCmd)

	standAloneCmd.Flags().String("db", "file", "db type, file or redis")
	rootCmd.AddCommand(standAloneCmd)
}

func getRepositoryFactory(storageType string) func(roomName string, app *app.App) (data.Repository, error) {
	switch storageType {
	case "file":
		return func(roomName string, app *app.App) (data.Repository, error) {
			s := storage.NewFileStorage(roomName)
			return data.New(s, app.NotifyNewData, app.NotifyListRemoved)
		}
	case "redis":
		addr, ok := os.LookupEnv("REDIS_ADDR")
		if !ok {
			addr = "localhost:6379"
		}
		password, ok := os.LookupEnv("REDIS_PASSWORD")
		if !ok {
			password = ""
		}
		rdb := storage.NewRedisClient(addr, password)

		return func(roomName string, app *app.App) (data.Repository, error) {
			s := storage.NewRedisStorage(rdb, roomName)
			return data.New(s, app.NotifyNewData, app.NotifyListRemoved)
		}
	}
	return nil
}
