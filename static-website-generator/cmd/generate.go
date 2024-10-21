package cmd

import (
	"github.com/spf13/cobra"
	"path/filepath"
	"static-website-generator/generator"
	transformerarticle "static-website-generator/generator/transformer-article"
	transformerbasepublicpath "static-website-generator/generator/transformer-basePublicPath"
	transformerFragment "static-website-generator/generator/transformer-fragment"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate static website",
	RunE: func(cmd *cobra.Command, args []string) error {
		basePathname, err := cmd.Flags().GetString("basePathname")
		if err != nil {
			return err
		}
		basePublicPath, err := cmd.Flags().GetString("basePublicPath")
		if err != nil {
			return err
		}

		runGenerator(basePathname, basePublicPath)

		return nil
	},
}

func runGenerator(basePathname, basePublicPath string) {

	g := generator.NewGenerator(generator.Config{
		ProcessableExtensions: []string{".html"},
		OutputDir:             filepath.Join(basePathname, "./output"),
		SrcDir:                filepath.Join(basePathname, "./html"),
	})

	basePublicPathTransformer := transformerbasepublicpath.NewTransformer(transformerbasepublicpath.Config{
		BasePublicPath: basePublicPath,
	})

	fragmentTransformer := transformerFragment.NewTransformer(transformerFragment.Config{
		FragmentSrcDir: filepath.Join(basePathname, "./fragments"),
	}, g)

	articleTransformer := transformerarticle.NewTransformer(
		transformerarticle.GetDefaultConfig(),
		g,
		fragmentTransformer,
	)

	// order matter
	g.RegisterTransformer(articleTransformer)
	g.RegisterTransformer(basePublicPathTransformer)
	g.RegisterTransformer(fragmentTransformer)

	g.Build()
}
