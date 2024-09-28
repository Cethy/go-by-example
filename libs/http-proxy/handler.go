package http_proxy

import (
	"errors"
	httpmiddleware "go-by-example/libs/http-middleware"
	"io"
	"net/http"
	"net/url"
)

func getUrl(r *http.Request, urlMode bool) (string, error) {
	destURL, err := url.Parse(r.URL.String())
	if urlMode {
		destURL, err = url.Parse(r.URL.Query().Get("url"))
	}
	if err != nil {
		//http.Error(w, "Parameter malformed", http.StatusBadRequest)
		return "", err
	}
	if destURL.String() == "" {
		//http.Error(w, "Parameter missing", http.StatusBadRequest)
		return "", errors.New("parameter missing")
	}

	return destURL.String(), nil
}

func GetProxyHandler(urlMode bool) httpmiddleware.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (status int, err error) {
		destUrl, err := getUrl(r, urlMode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return http.StatusBadRequest, err
		}

		res, getErr := http.Get(destUrl)
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
