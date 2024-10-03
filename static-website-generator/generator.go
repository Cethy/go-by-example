package main

import (
	"flag"
	"go-by-example/static-website-generator/generator"
	transformerarticle "go-by-example/static-website-generator/generator/transformer-article"
	transformerFragment "go-by-example/static-website-generator/generator/transformer-fragment"
	"path/filepath"
)

func main() {
	basePathname := flag.String("basePathname", "./static-website-generator", "pathname to src files")
	flag.Parse()

	g := generator.NewGenerator(generator.Config{
		ProcessableExtensions: []string{".html"},
		OutputDir:             filepath.Join(*basePathname, "./output"),
		SrcDir:                filepath.Join(*basePathname, "./html"),
	})

	fragmentTransformer := transformerFragment.NewTransformer(transformerFragment.Config{
		FragmentSrcDir: filepath.Join(*basePathname, "./fragments"),
	})

	articleTransformer := transformerarticle.NewTransformer(
		transformerarticle.GetDefaultConfig(),
		g,
		fragmentTransformer,
	)

	// order matter
	g.RegisterTransformer(articleTransformer)
	g.RegisterTransformer(fragmentTransformer)

	g.Build()
}
