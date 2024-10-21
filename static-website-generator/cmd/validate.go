package cmd

import (
	"github.com/spf13/cobra"
	"static-website-generator/validator"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validate articles README metadata and add a random unsplash if no ImgSrc value is set",
	RunE: func(cmd *cobra.Command, args []string) error {
		return validator.Validate()
	},
}
