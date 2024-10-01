package main

import (
	"flag"
	"go-by-example/static-website-generator/generator"
	plugin_article "go-by-example/static-website-generator/plugin-article"
	pluginfragment "go-by-example/static-website-generator/plugin-fragment"
	"path/filepath"
)

func main() {
	srcFilePathname := flag.String("srcFilePathname", "./static-website-generator", "pathname to src files")
	flag.Parse()

	generator.Build(generator.Config{
		ProcessableExtensions: []string{".html"},
		PreFileProcessors:     []func(fileContent string, config generator.Config) (string, error){plugin_article.PreBuildFile, pluginfragment.PreBuildFile},
		OutputDir:             filepath.Join(*srcFilePathname, "./output"),
		SrcDir:                filepath.Join(*srcFilePathname, "./html"),
	})
}
