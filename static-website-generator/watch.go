package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"go-by-example/static-website-generator/generator"
	transformerarticle "go-by-example/static-website-generator/generator/transformer-article"
	transformerbasepublicpath "go-by-example/static-website-generator/generator/transformer-basePublicPath"
	transformerFragment "go-by-example/static-website-generator/generator/transformer-fragment"
	"os"
	"path/filepath"
	"time"
)

func exit(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, filepath.Base(os.Args[0])+": "+format+"\n", a...)
	os.Exit(1)
}

// Print line prefixed with the time (a bit shorter than log.Print; we don't
// really need the date and ms is useful here).
func printTime(s string, args ...interface{}) {
	fmt.Printf(time.Now().Format("15:04:05.0000")+" "+s+"\n", args...)
}

func watch(callback func(), paths ...string) {
	if len(paths) < 1 {
		exit("must specify at least one path to watch")
	}

	// Create a new watcher.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		exit("creating a new watcher: %s", err)
	}
	defer w.Close()

	// Start listening for events.
	go watchLoop(w, callback)

	// Add all paths from the commandline.
	for _, p := range paths {
		err = w.Add(p)
		if err != nil {
			exit("%q: %s", p, err)
		}
	}

	printTime("ready; press ^C to exit")
	<-make(chan struct{}) // Block forever
}

func watchLoop(w *fsnotify.Watcher, callback func()) {
	watchedOps := []fsnotify.Op{fsnotify.Write, fsnotify.Create, fsnotify.Rename, fsnotify.Remove}
	for {
		select {
		// Read from Errors.
		case err, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			printTime("ERROR: %s", err)
		// Read from Events.
		case e, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			for _, watchedOp := range watchedOps {
				if e.Op.Has(watchedOp) && e.Name[len(e.Name)-1] != '~' {
					printTime("%s", e.Op)
					callback()
					break
				}
			}
		}
	}
}

func getAllSubDir(srcDir string) []string {
	dir, err := os.ReadDir(srcDir)
	if err != nil {
		panic(err)
	}

	subDirs := []string{srcDir}

	for _, file := range dir {
		if file.IsDir() {
			path := filepath.Join(srcDir, file.Name())
			subDirs = append(subDirs, path)
			subDirs = append(subDirs, getAllSubDir(path)...)
		}
	}

	return subDirs
}

func RunGenerator(basePathname, basePublicPath string) {
	g := generator.NewGenerator(generator.Config{
		ProcessableExtensions: []string{".html"},
		OutputDir:             filepath.Join(basePathname, "./output"),
		SrcDir:                filepath.Join(basePathname, "./html"),
	})

	basePublicPathTransformer := transformerbasepublicpath.NewTransformer(transformerbasepublicpath.Config{
		BasePublicPath: basePublicPath,
	})

	fragmentTransformer := transformerFragment.NewTransformer(transformerFragment.Config{
		FragmentSrcDir: filepath.Join(basePathname, "./fragments"),
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

func main() {
	basePathname := flag.String("basePathname", "./static-website-generator", "pathname to src files")
	basePublicPath := flag.String("basePublicPath", "/", "base public path (eg: https://cethy.github.io/go-by-example/)")
	flag.Parse()

	watchFilePaths := append(
		getAllSubDir(filepath.Join(*basePathname, "./html")),
		getAllSubDir(filepath.Join(*basePathname, "./fragments"))...,
	)

	watch(func() {
		RunGenerator(*basePathname, *basePublicPath)
	}, watchFilePaths...)
}
