package http_proxy

import (
	"errors"
	httpmiddleware "go-by-example/libs/http-middleware"
	"io"
	"net/http"
	"net/url"
)

func GetProxyHandler(urlMode bool) httpmiddleware.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (status int, err error) {
		destURL, destErr := url.Parse(r.URL.String())
		if urlMode {
			destURL, destErr = url.Parse(r.URL.Query().Get("url"))
			if destURL.String() == "" {
				http.Error(w, "Parameter missing", http.StatusBadRequest)
				return http.StatusBadRequest, errors.New("parameter missing")
			}
		}
		if destErr != nil {
			http.Error(w, "Parameter malformed", http.StatusBadRequest)
			return http.StatusBadRequest, destErr
		}

		res, getErr := http.Get(destURL.String())
		if getErr != nil {
			http.Error(w, "Wrong parameter format or bad reply from target destination", http.StatusBadRequest)
			return http.StatusBadRequest, getErr
		}

		_, err = io.Copy(w, res.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return http.StatusBadRequest, err
		}
		return http.StatusOK, nil
	}
}
