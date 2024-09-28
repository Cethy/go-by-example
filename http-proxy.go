package main

import (
	"flag"
	"fmt"
	http_middleware "go-by-example/http-middleware"
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
func getProxyHandler(urlMode bool) http.HandlerFunc {
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

func main() {
	port := flag.Int("p", 8003, "Port number")
	urlMode := flag.Bool("url", false, "tells the executable how the proxy should read the target url ; false for default proxy mode, true for url mode : 'http://localhost:8003/?url=https://example.com'")
	flag.Parse()

	http.HandleFunc("/", http_middleware.Chain(getProxyHandler(*urlMode), http_middleware.Method("GET"), http_middleware.Logging()))

	fmt.Println("Server listening on port:", *port)
	err := http.ListenAndServe(":"+fmt.Sprint(*port), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
