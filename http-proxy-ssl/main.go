package main

import (
	"flag"
	"fmt"
	httpmiddleware "go-by-example/libs/http-middleware"
	httpproxy "go-by-example/libs/http-proxy"
	"log"
	"net/http"
)

func main() {
	certFile := flag.String("cert", "", "Certificate file path")
	keyFile := flag.String("key", "", "Private key file path")

	port := flag.Int("p", 8003, "Port number")
	portSSL := flag.Int("ps", 8004, "Port number")
	urlMode := flag.Bool("url", false, "tells the executable how the proxy should read the target url ; false for default proxy mode, true for url mode : 'http://localhost:8003/?url=https://example.com'")
	flag.Parse()

	if *certFile == "" || *keyFile == "" {
		log.Fatal("cert and key params are mandatory")
	}

	httpServer := http.NewServeMux()
	httpServer.HandleFunc("/", httpmiddleware.Chain(httpproxy.GetProxyHandler(*urlMode), httpmiddleware.LoggingPre("[HTTP ]"), httpmiddleware.PostLogging()))

	httpsServer := http.NewServeMux()
	httpsServer.HandleFunc("/", httpmiddleware.Chain(httpproxy.GetProxyHandler(*urlMode), httpmiddleware.LoggingPre("[HTTPS]"), httpmiddleware.PostLogging()))

	fmt.Println("http listening on port:", *port)
	fmt.Println("https listening on port:", *portSSL)

	go func() {
		err := http.ListenAndServe(":"+fmt.Sprint(*port), httpServer)
		if err != nil {
			log.Fatal("[http]", err)
			return
		}
	}()
	err := http.ListenAndServeTLS(":"+fmt.Sprint(*portSSL), *certFile, *keyFile, httpsServer)
	if err != nil {
		log.Fatal("[https]", err)
		return
	}
}
