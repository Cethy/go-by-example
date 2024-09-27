package main

import (
	"fmt"
	"net/http"
)

const port = 8001

func main() {
	//http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)

		fmt.Fprintf(w, "Welcome to my website!")
	})

	publicDir := http.FileServer(http.Dir("public/"))
	http.Handle("/public/", http.StripPrefix("/public/", publicDir))
	http.Handle("/favicon.ico", publicDir)

	http.ListenAndServe(":"+fmt.Sprint(port), nil)
}
