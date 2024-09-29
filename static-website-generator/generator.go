package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
)

var processableExtensions = []string{".html"}

func copyFile(srcPath, targetPath string) {
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

	log.Println("copied", srcPath, "to", targetPath)
}

func buildFile(srcPath, targetPath string) {
	readFile, err := os.ReadFile(srcPath)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(targetPath, readFile, 0644) // @todo check params
	if err != nil {
		panic(err)
	}

	log.Println("built", srcPath, "to", targetPath)
}

func processFile(srcPath, targetPath string) {
	ext := filepath.Ext(srcPath)
	if !slices.Contains(processableExtensions, ext) {
		copyFile(srcPath, targetPath)
		return
	}

	buildFile(srcPath, targetPath)
}

func processDir(dirPath, inputPath, outputPath string) {
	dir, err := os.ReadDir(filepath.Join(inputPath, dirPath))
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Join(outputPath, dirPath), 0755) // @todo check perms
	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		if file.IsDir() {
			processDir(filepath.Join(dirPath, file.Name()), inputPath, outputPath)
			continue
		}
		srcPath := filepath.Join(inputPath, dirPath, file.Name())
		targetPath := filepath.Join(outputPath, dirPath, file.Name())

		processFile(srcPath, targetPath)
	}
}

func build(srcFilePathname string) {
	outputDir := filepath.Join(srcFilePathname, "./output")
	srcDir := filepath.Join(srcFilePathname, "./html")

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

	processDir(".", srcDir, outputDir)
}

func main() {
	srcFilePathname := flag.String("srcFilePathname", "./static-website-generator", "pathname to src files")
	flag.Parse()

	build(*srcFilePathname)
}
