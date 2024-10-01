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

func BuildFile(targetPath, fileContent string) {
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

func PreProcess(fileContent string) (string, error) {
	var err error
	for _, preProcessor := range config.PreFileProcessors {
		fileContent, err = preProcessor(fileContent, config)
		if err != nil {
			return fileContent, err
		}
	}
	return fileContent, nil
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

	fileContent, err := PreProcess(string(readFile))
	if err != nil {
		panic(err)
	}

	BuildFile(targetPath, fileContent)
}

func processDir(dirPath string) {
	dir, err := os.ReadDir(filepath.Join(config.SrcDir, dirPath))
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Join(config.OutputDir, dirPath), 0755) // @todo check perms
	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		if file.IsDir() {
			processDir(filepath.Join(dirPath, file.Name()))
			continue
		}
		srcPath := filepath.Join(config.SrcDir, dirPath, file.Name())
		targetPath := filepath.Join(config.OutputDir, dirPath, file.Name())

		processFile(srcPath, targetPath)
	}
}

type Config struct {
	ProcessableExtensions []string // = []string{".html"}
	PreFileProcessors     []func(fileContent string, config Config) (string, error)
	SrcDir                string
	OutputDir             string
}

var config Config

func Build(newConfig Config) {
	config = newConfig

	// clean output directory
	err := os.RemoveAll(config.OutputDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	}
	err = os.MkdirAll(config.OutputDir, 0755) // @todo check perms
	if err != nil {
		panic(err)
	}

	processDir(".")
}
