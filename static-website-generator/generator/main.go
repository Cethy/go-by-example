package generator

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
)

func CopyFile(srcPath, targetPath string) {
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

func buildFile(targetPath, fileContent string) {
	err := os.WriteFile(targetPath, []byte(fileContent), 0644)
	if err != nil {
		panic(err)
	}

	log.Println("built", targetPath)
}

func processFile(srcPath, targetPath string) {
	ext := filepath.Ext(srcPath)
	if !slices.Contains(config.ProcessableExtensions, ext) {
		CopyFile(srcPath, targetPath)
		return
	}

	readFile, err := os.ReadFile(srcPath)
	if err != nil {
		panic(err)
	}

	fileContent := string(readFile)
	for _, preProcessor := range config.PreFileProcessors {
		fileContent, err = preProcessor(fileContent, srcDir)
		if err != nil {
			return
		}
	}

	buildFile(targetPath, fileContent)
}

func processDir(dirPath string) {
	dir, err := os.ReadDir(filepath.Join(srcDir, dirPath))
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Join(outputDir, dirPath), 0755) // @todo check perms
	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		if file.IsDir() {
			processDir(filepath.Join(dirPath, file.Name()))
			continue
		}
		srcPath := filepath.Join(srcDir, dirPath, file.Name())
		targetPath := filepath.Join(outputDir, dirPath, file.Name())

		processFile(srcPath, targetPath)
	}
}

type Config struct {
	ProcessableExtensions []string // = []string{".html"}
	PreFileProcessors     []func(fileContent, srcDir string) (string, error)
}

var outputDir string
var srcDir string
var config Config

func Build(srcFilePathname string, newConfig Config) {
	config = newConfig
	outputDir = filepath.Join(srcFilePathname, "./output")
	srcDir = filepath.Join(srcFilePathname, "./html")

	// clean output directory
	err := os.RemoveAll(outputDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	}
	err = os.MkdirAll(outputDir, 0755) // @todo check perms
	if err != nil {
		panic(err)
	}

	processDir(".")
}
