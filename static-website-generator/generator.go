package main

import (
	"flag"
	"go-by-example/static-website-generator/generator"
	transformerarticle "go-by-example/static-website-generator/generator/transformer-article"
	transformerbasepublicpath "go-by-example/static-website-generator/generator/transformer-basePublicPath"
	transformerFragment "go-by-example/static-website-generator/generator/transformer-fragment"
	"path/filepath"
)

func main() {
	basePathname := flag.String("basePathname", "./static-website-generator", "pathname to src files")
	basePublicPath := flag.String("basePublicPath", "/", "base public path (eg: https://cethy.github.io/go-by-example/)")
	flag.Parse()

	g := generator.NewGenerator(generator.Config{
		ProcessableExtensions: []string{".html"},
		OutputDir:             filepath.Join(*basePathname, "./output"),
		SrcDir:                filepath.Join(*basePathname, "./html"),
	})

	basePublicPathTransformer := transformerbasepublicpath.NewTransformer(transformerbasepublicpath.Config{
		BasePublicPath: *basePublicPath,
	})

	fragmentTransformer := transformerFragment.NewTransformer(transformerFragment.Config{
		FragmentSrcDir: filepath.Join(*basePathname, "./fragments"),
	}, g)

	articleTransformer := transformerarticle.NewTransformer(
		transformerarticle.GetDefaultConfig(),
		g,
		fragmentTransformer,
	)

	// order matter
	g.RegisterTransformer(basePublicPathTransformer)
	g.RegisterTransformer(articleTransformer)
	g.RegisterTransformer(fragmentTransformer)

	g.Build()
}
