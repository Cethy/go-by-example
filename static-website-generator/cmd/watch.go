package cmd

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch for file change & regenerate website",
	RunE: func(cmd *cobra.Command, args []string) error {
		basePathname, err := cmd.Flags().GetString("basePathname")
		if err != nil {
			return err
		}

		basePublicPath, err := cmd.Flags().GetString("basePublicPath")
		if err != nil {
			return err
		}

		watchFilePaths := append(
			getAllSubDir(filepath.Join(basePathname, "./html")),
			getAllSubDir(filepath.Join(basePathname, "./fragments"))...,
		)

		return watch(func() {
			runGenerator(basePathname, basePublicPath)
		}, watchFilePaths...)
	},
}

func watch(callback func(), paths ...string) error {
	if len(paths) < 1 {
		return errors.New("must specify at least one path to watch")
	}

	// Create a new watcher.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	// Start listening for events.
	go watchLoop(w, callback)

	// Add all paths from the commandline.
	for _, p := range paths {
		err = w.Add(p)
		if err != nil {
			return err
		}
	}

	log.Println("ready; press ^C to exit")
	<-make(chan struct{}) // Block forever

	return nil
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
			log.Printf("ERROR: %s", err)
		// Read from Events.
		case e, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			for _, watchedOp := range watchedOps {
				if e.Op.Has(watchedOp) && e.Name[len(e.Name)-1] != '~' {
					log.Printf("op:%s,file:%s", e.Op, e.Name)
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
			subDirs = append(subDirs, getAllSubDir(path)...)
		}
	}

	return subDirs
}
