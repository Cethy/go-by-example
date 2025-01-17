package main

import (
	"flag"
	"fmt"
	"http-proxy-cache/handler"
	httpmiddleware "http-server-middleware/http-middleware"
	"log"
	"net/http"
)

func listenHttp(urlMode bool, port int) {
	httpServer := http.NewServeMux()
	httpServer.HandleFunc("/", httpmiddleware.Chain(handler.GetProxyCacheHandler(urlMode), httpmiddleware.LoggingPre("[HTTP ]"), httpmiddleware.PostLogging()))

	fmt.Println("http listening on port:", port)

	err := http.ListenAndServe(":"+fmt.Sprint(port), httpServer)
	if err != nil {
		log.Fatal("[http ]", err)
		return
	}
}

func listenHttps(urlMode bool, port int, certFile string, keyFile string) {
	httpsServer := http.NewServeMux()
	httpsServer.HandleFunc("/", httpmiddleware.Chain(handler.GetProxyCacheHandler(urlMode), httpmiddleware.LoggingPre("[HTTPS]"), httpmiddleware.PostLogging()))

	fmt.Println("https listening on port:", port)

	err := http.ListenAndServeTLS(":"+fmt.Sprint(port), certFile, keyFile, httpsServer)
	if err != nil {
		log.Fatal("[https]", err)
		return
	}
}

func main() {
	certFile := flag.String("cert", "", "Certificate file path")
	keyFile := flag.String("key", "", "Private key file path")

	port := flag.Int("p", 8005, "Port number")
	portSSL := flag.Int("ps", 8006, "Port number")
	urlMode := flag.Bool("url", false, "tells the executable how the proxy should read the target url ; false for default proxy mode, true for url mode : 'http://localhost:8003/?url=https://example.com'")
	flag.Parse()

	if *certFile != "" && *keyFile != "" {
		go listenHttps(*urlMode, *portSSL, *certFile, *keyFile)
	}
	listenHttp(*urlMode, *port)
}
