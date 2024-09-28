package main

import (
	"fmt"
	"go-by-example/http-middleware"
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

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my website!")
}

func main() {
	http.HandleFunc("/", http_middleware.Chain(middlewares(Hello), customMiddleware()))

	publicDir := http.FileServer(http.Dir("public/"))
	http.HandleFunc("/public/", middlewares(http.StripPrefix("/public/", publicDir).ServeHTTP))
	http.HandleFunc("/favicon.ico", middlewares(publicDir.ServeHTTP))

	http.ListenAndServe(":"+fmt.Sprint(8002), nil)
}
