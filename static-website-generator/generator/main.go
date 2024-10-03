package generator

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
)

type Config struct {
	ProcessableExtensions []string // = []string{".html"}
	SrcDir                string
	OutputDir             string
}

type Transformer interface {
	Transform(content string) (string, error)
}

type Generator interface {
	Build()
	RegisterTransformer(transformer Transformer)
	InvokeTransformers(fileContent string) (string, error)
	BuildFile(relativeFilePath, fileContent string)
}

type generator struct {
	Config       Config
	transformers []Transformer
}

func (g *generator) copyFile(relativeFilePath string) {
	srcPath := filepath.Join(g.Config.SrcDir, relativeFilePath)
	targetPath := filepath.Join(g.Config.OutputDir, relativeFilePath)

	// simply copy the file
	r, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	w, err := os.Create(targetPath)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		panic(err)
	}
	// https://stackoverflow.com/a/13038961
	w.Sync()

	log.Println("copied", targetPath)
}

func (g *generator) BuildFile(relativeTargetFilePath, fileContent string) {
	targetPath := filepath.Join(g.Config.OutputDir, relativeTargetFilePath)

	err := os.MkdirAll(filepath.Dir(targetPath), 0755) // @todo check perms
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(targetPath, []byte(fileContent), 0644)
	if err != nil {
		panic(err)
	}

	log.Println("built", targetPath)
}

func (g *generator) InvokeTransformers(fileContent string) (string, error) {
	var err error
	for _, transformer := range g.transformers {
		fileContent, err = transformer.Transform(fileContent)
		if err != nil {
			return fileContent, err
		}
	}
	return fileContent, nil
}

func (g *generator) processFile(relativeFilePath string) {
	ext := filepath.Ext(relativeFilePath)
	if !slices.Contains(g.Config.ProcessableExtensions, ext) {
		g.copyFile(relativeFilePath)
		return
	}

	readFile, err := os.ReadFile(filepath.Join(g.Config.SrcDir, relativeFilePath))
	if err != nil {
		panic(err)
	}

	fileContent, err := g.InvokeTransformers(string(readFile))
	if err != nil {
		panic(err)
	}

	g.BuildFile(relativeFilePath, fileContent)
}

func (g *generator) processHtmlDir(relativeDirPath string) {
	dir, err := os.ReadDir(filepath.Join(g.Config.SrcDir, relativeDirPath))
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Join(g.Config.OutputDir, relativeDirPath), 0755) // @todo check perms
	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		if file.IsDir() {
			g.processHtmlDir(filepath.Join(relativeDirPath, file.Name()))
			continue
		}
		g.processFile(filepath.Join(relativeDirPath, file.Name()))
	}
}

func (g *generator) Build() {
	// clean output directory
	err := os.RemoveAll(g.Config.OutputDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	}
	err = os.MkdirAll(g.Config.OutputDir, 0755) // @todo check perms
	if err != nil {
		panic(err)
	}

	g.processHtmlDir(".")
}

func (g *generator) RegisterTransformer(processor Transformer) {
	g.transformers = append(g.transformers, processor)
}

func NewGenerator(config Config) Generator {
	return &generator{Config: config}
}
