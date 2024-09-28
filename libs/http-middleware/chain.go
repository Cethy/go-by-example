package http_middleware

import "net/http"

func Handle(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}

func reverseHandle(f http.HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (int, error) {
		f(w, r)
		return 200, nil
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return Handle(f)
}

// ChainOG does NOT support post middlewares !
func ChainOG(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	return Chain(reverseHandle(f), middlewares...)
}
