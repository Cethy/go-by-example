package transformer_fragment

import (
	"go-by-example/static-website-generator/generator"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Config struct {
	FragmentSrcDir string
}

type FragmentTransformer interface {
	generator.Transformer
	GetFragmentContent(fragmentId string) string
}

type fragmentTransformer struct {
	Config    Config
	fragments map[string]string
	generator generator.Generator
}

func (p *fragmentTransformer) GetFragmentContent(fragmentId string) string {
	fragmentPath := filepath.Join(p.Config.FragmentSrcDir, fragmentId+".html")

	if p.fragments[fragmentPath] == "" {
		readFile, err := os.ReadFile(fragmentPath)
		if err != nil {
			//log.Println("WARNING", "cannot find file for fragment", fragmentPath)
			return "{" + fragmentId + "}"
		}

		fileContent, err := p.generator.InvokeTransformers(string(readFile))

		p.fragments[fragmentPath] = fileContent
	}
	return p.fragments[fragmentPath]
}

func (p *fragmentTransformer) getAllFragmentIds(fileContent string) []string {
	pattern := "(\\{[A-z-_]+\\})"
	r, _ := regexp.Compile(pattern)
	return r.FindAllString(fileContent, -1)
}

func (p *fragmentTransformer) Transform(fileContent string) (string, error) {
	// insert fragments
	requiredFragments := p.getAllFragmentIds(fileContent)
	for _, fragmentId := range requiredFragments {
		fragment := p.GetFragmentContent(fragmentId[1 : len(fragmentId)-1])
		fileContent = strings.ReplaceAll(fileContent, fragmentId, fragment)
	}

	return fileContent, nil
}

func NewTransformer(config Config, generator generator.Generator) FragmentTransformer {
	return &fragmentTransformer{
		Config:    config,
		fragments: make(map[string]string),
		generator: generator,
	}
}
