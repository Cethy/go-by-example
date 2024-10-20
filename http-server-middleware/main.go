package main

import (
	"flag"
	"fmt"
	httpmiddleware "http-server-middleware/http-middleware"
	"log"
	"net/http"
)

func customMiddleware() httpmiddleware.Middleware {
	return httpmiddleware.CreateNewMiddleware(func(w http.ResponseWriter, r *http.Request) {
		defer func() { log.Println("Hi! from customMiddleware") }()
	})
}

func Hello(w http.ResponseWriter, _ *http.Request) (status int, err error) {
	_, err = fmt.Fprintf(w, "Welcome to my website!")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 400, err
	}
	return 200, nil
}

func main() {
	port := flag.Int("p", 8002, "Port number")
	flag.Parse()

	http.HandleFunc("/", httpmiddleware.Chain(Hello, customMiddleware(), httpmiddleware.Method("GET"), httpmiddleware.Logging(), httpmiddleware.PostLogging()))

	publicDir := http.FileServer(http.Dir("public/"))
	http.HandleFunc("/public/", httpmiddleware.ChainOG(http.StripPrefix("/public/", publicDir).ServeHTTP, httpmiddleware.Logging()))
	http.HandleFunc("/favicon.ico", httpmiddleware.ChainOG(publicDir.ServeHTTP, httpmiddleware.Logging()))

	fmt.Println("Server listening on port:", *port)
	err := http.ListenAndServe(":"+fmt.Sprint(*port), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
