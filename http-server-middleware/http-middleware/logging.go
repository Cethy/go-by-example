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

// LoggingPre with prefix
func LoggingPre(prefix string) Middleware {
	return CreateNewMiddleware(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() { log.Println(prefix, time.Since(start), "|", r.Method, r.URL) }()
	})
}

// PostLogging post handler
func PostLogging() Middleware {
	return CreateNewPostMiddleware(func(w http.ResponseWriter, r *http.Request, status int, err error) {
		if err != nil {
			log.Println(status, err)
			return
		}
		log.Println(status, "OK")
	})
}
