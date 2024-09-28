package http_proxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func failure(w http.ResponseWriter, error string, code int) {
	http.Error(w, error, code)
	log.Println(code, error)
}
func success(w http.ResponseWriter, src io.Reader) {
	_, err := io.Copy(w, src)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Println(200, "OK")
}

func GetProxyHandler(urlMode bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		destURL, destErr := url.Parse(r.URL.String())
		if urlMode {
			destURL, destErr = url.Parse(r.URL.Query().Get("url"))
			if destURL.String() == "" {
				failure(w, "Parameter missing", http.StatusBadRequest)
				return
			}
		}
		if destErr != nil {
			failure(w, "Parameter malformed", http.StatusBadRequest)
			return
		}

		res, getErr := http.Get(destURL.String())
		if getErr != nil {
			failure(w, "Wrong parameter format or bad reply from target destination", http.StatusBadRequest)
			return
		}

		success(w, res.Body)
	}
}
