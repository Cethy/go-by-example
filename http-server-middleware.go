package main

import (
	"fmt"
	"go-by-example/http-middleware"
	"net/http"
)

func middlewares(f http.HandlerFunc) http.HandlerFunc {
	return http_middleware.Chain(f, http_middleware.Method("GET"), http_middleware.Logging())
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my website!")
}

func main() {
	http.HandleFunc("/", middlewares(Hello))

	publicDir := http.FileServer(http.Dir("public/"))
	http.HandleFunc("/public/", middlewares(http.StripPrefix("/public/", publicDir).ServeHTTP))
	http.HandleFunc("/favicon.ico", middlewares(publicDir.ServeHTTP))

	http.ListenAndServe(":"+fmt.Sprint(8002), nil)
}
