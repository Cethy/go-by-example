package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
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

func removeExt(filename string) string {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

func isFragment(filename string) bool {
	runeSlice := []rune(filename)
	return runeSlice[0] == '{' && runeSlice[len(runeSlice)-1] == '}'
}

var fragments map[string]string = make(map[string]string)

func getFragment(fragmentId string) string {
	if fragments[fragmentId] == "" {
		readFile, err := os.ReadFile(fragmentId)
		if err != nil {
			panic(err)
		}

		fragments[fragmentId] = string(readFile)
	}
	return fragments[fragmentId]
}

func buildFile(srcPath, targetPath string) {
	// ignore if it's a fragment
	filename := removeExt(srcPath)
	if isFragment(filename) {
		return
	}

	readFile, err := os.ReadFile(srcPath)
	if err != nil {
		panic(err)
	}
	fileContent := string(readFile)

	// insert fragments
	pattern := "(\\{[A-z]+\\})"
	r, _ := regexp.Compile(pattern)
	requiredFragments := r.FindAllString(fileContent, -1)
	for _, fragmentId := range requiredFragments {
		fragment := getFragment(filepath.Join(filepath.Dir(srcPath), fragmentId+".html"))
		fileContent = strings.ReplaceAll(fileContent, fragmentId, fragment)
	}

	err = os.WriteFile(targetPath, []byte(fileContent), 0644)
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

var outputDir string
var srcDir string

func build(srcFilePathname string) {
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

func main() {
	srcFilePathname := flag.String("srcFilePathname", "./static-website-generator", "pathname to src files")
	flag.Parse()

	build(*srcFilePathname)
}
