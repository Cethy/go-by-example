package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("p", 8001, "Port number")
	flag.Parse()

	//http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)

		fmt.Fprintf(w, "Welcome to my website!")
	})

	publicDir := http.FileServer(http.Dir("public/"))
	http.Handle("/public/", http.StripPrefix("/public/", publicDir))
	http.Handle("/favicon.ico", publicDir)

	fmt.Println("Server listening on port:", *port)
	err := http.ListenAndServe(":"+fmt.Sprint(*port), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
