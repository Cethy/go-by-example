package main

import (
	"flag"
	"path/filepath"
	"static-website-generator/generator"
	transformerarticle "static-website-generator/generator/transformer-article"
	transformerbasepublicpath "static-website-generator/generator/transformer-basePublicPath"
	transformerFragment "static-website-generator/generator/transformer-fragment"
)

func main() {
	basePathname := flag.String("basePathname", "./", "pathname to src files")
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
	g.RegisterTransformer(articleTransformer)
	g.RegisterTransformer(basePublicPathTransformer)
	g.RegisterTransformer(fragmentTransformer)

	g.Build()
}
