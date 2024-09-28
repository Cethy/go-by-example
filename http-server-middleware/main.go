package main

import (
	"flag"
	"fmt"
	http_middleware "go-by-example/libs/http-middleware"
	"log"
	"net/http"
)

func middlewares(f http.HandlerFunc) http.HandlerFunc {
	return http_middleware.Chain(f, http_middleware.Method("GET"), http_middleware.Logging())
}

func customMiddleware() http_middleware.Middleware {
	return http_middleware.CreateNewMiddleware(func(w http.ResponseWriter, r *http.Request) {
		defer func() { log.Println("Hi! from customMiddleware") }()
	})
}

func Hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Welcome to my website!")
}

func main() {
	port := flag.Int("p", 8002, "Port number")
	flag.Parse()

	http.HandleFunc("/", http_middleware.Chain(middlewares(Hello), customMiddleware()))

	publicDir := http.FileServer(http.Dir("public/"))
	http.HandleFunc("/public/", middlewares(http.StripPrefix("/public/", publicDir).ServeHTTP))
	http.HandleFunc("/favicon.ico", middlewares(publicDir.ServeHTTP))

	fmt.Println("Server listening on port:", *port)
	err := http.ListenAndServe(":"+fmt.Sprint(*port), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
