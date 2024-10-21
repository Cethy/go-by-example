package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"path/filepath"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve generated website",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
		srcFilePathname, err := cmd.Flags().GetString("srcFilePathname")
		if err != nil {
			return err
		}

		http.Handle("/", http.FileServer(http.Dir(filepath.Join(srcFilePathname, "./output"))))

		fmt.Println("Server listening on port:", port)
		err = http.ListenAndServe(":"+fmt.Sprint(port), nil)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}
