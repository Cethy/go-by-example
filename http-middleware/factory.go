package http_middleware

import "net/http"

func CreateNewMiddleware(middleware func(w http.ResponseWriter, r *http.Request)) Middleware {
	// Create a new Middleware
	return func(next http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc which is called by the server eventually
		return func(w http.ResponseWriter, r *http.Request) {
			// ... do http-middleware things
			middleware(w, r)
			// Call the next http-middleware/handler in chain
			next(w, r)
		}
	}
}
