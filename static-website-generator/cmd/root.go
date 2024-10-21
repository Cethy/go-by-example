package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "static-website-generator",
	Short: "Opinionated static site generator",
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
	rootCmd.AddCommand(validateCmd)

	generateCmd.Flags().String("basePathname", "./", "pathname to src files")
	generateCmd.Flags().String("basePublicPath", "/", "base public path (eg: https://cethy.github.io/go-by-example/)")
	rootCmd.AddCommand(generateCmd)

	watchCmd.Flags().String("basePathname", "./", "pathname to src files")
	watchCmd.Flags().String("basePublicPath", "/", "base public path (eg: https://cethy.github.io/go-by-example/)")
	rootCmd.AddCommand(watchCmd)

	serveCmd.Flags().IntP("port", "p", 8007, "port to serve on")
	serveCmd.Flags().String("srcFilePathname", "./", "pathname to src files")
	rootCmd.AddCommand(serveCmd)
}
