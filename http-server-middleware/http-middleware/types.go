package http_middleware

import "net/http"

type HandlerFunc func(w http.ResponseWriter, r *http.Request) (status int, err error)

type Middleware func(HandlerFunc) HandlerFunc
