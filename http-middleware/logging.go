package http_middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging() Middleware {
	return CreateNewMiddleware(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() { log.Println(time.Since(start), "|", r.Method, r.URL) }()
	})
}
