package http_middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging logs all requests with its path and the time it took to process
func Logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do http-middleware things
			start := time.Now()
			defer func() { log.Println(time.Since(start), "|", r.Method, r.URL) }()

			// Call the next http-middleware/handler in chain
			f(w, r)
		}
	}
}
