package http_middleware

import "net/http"

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
func Method(m string) Middleware {
	return CreateNewMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != m {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	})
}
