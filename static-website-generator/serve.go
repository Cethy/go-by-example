package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	port := flag.Int("p", 8007, "Port number")
	srcFilePathname := flag.String("srcFilePathname", "./", "pathname to src files")

	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(filepath.Join(*srcFilePathname, "./output"))))

	fmt.Println("Server listening on port:", *port)
	err := http.ListenAndServe(":"+fmt.Sprint(*port), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
