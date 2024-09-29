package main

import (
	"flag"
	"go-by-example/static-website-generator/generator"
	plugin_article "go-by-example/static-website-generator/plugin-article"
	pluginfragment "go-by-example/static-website-generator/plugin-fragment"
)

func main() {
	srcFilePathname := flag.String("srcFilePathname", "./static-website-generator", "pathname to src files")
	flag.Parse()

	generator.Build(*srcFilePathname, generator.Config{
		ProcessableExtensions: []string{".html"},
		PreFileProcessors:     []func(fileContent, srcDir string) (string, error){plugin_article.PreBuildFile, pluginfragment.PreBuildFile},
	})
}
