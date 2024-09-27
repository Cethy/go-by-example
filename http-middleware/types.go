package http_middleware

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

/*
func createNewMiddleware() Middleware {

	// Create a new Middleware
	middleware := func(next http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc which is called by the server eventually
		handler := func(w http.ResponseWriter, r *http.Request) {

			// ... do http-middleware things

			// Call the next http-middleware/handler in chain
			next(w, r)
		}

		// Return newly created handler
		return handler
	}

	// Return newly created http-middleware
	return middleware
}
*/
