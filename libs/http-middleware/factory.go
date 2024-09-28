package http_middleware

import "net/http"

func CreateNewMiddleware(middleware func(w http.ResponseWriter, r *http.Request)) Middleware {
	// Create a new Middleware
	return func(next HandlerFunc) HandlerFunc {
		// Define the http.HandlerFunc which is called by the server eventually
		return func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			// ... do middleware things
			middleware(w, r)
			// Call the next http-middleware/handler in chain
			return next(w, r)
		}
	}
}

func CreateNewPostMiddleware(middleware func(w http.ResponseWriter, r *http.Request, status int, err error)) Middleware {
	// Create a new Middleware
	return func(next HandlerFunc) HandlerFunc {
		// Define the http.HandlerFunc which is called by the server eventually
		return func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			// Call the next http-middleware/handler in chain
			status, err = next(w, r)

			// ... do middleware things
			middleware(w, r, status, err)

			return status, err
		}
	}
}
