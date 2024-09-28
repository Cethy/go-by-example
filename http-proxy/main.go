package main

import (
	"flag"
	"fmt"
	http_middleware "go-by-example/libs/http-middleware"
	http_proxy "go-by-example/libs/http-proxy"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("p", 8003, "Port number")
	urlMode := flag.Bool("url", false, "tells the executable how the proxy should read the target url ; false for default proxy mode, true for url mode : 'http://localhost:8003/?url=https://example.com'")
	flag.Parse()

	http.HandleFunc("/", http_middleware.Chain(http_proxy.GetProxyHandler(*urlMode), http_middleware.Logging()))

	fmt.Println("Server listening on port:", *port)
	err := http.ListenAndServe(":"+fmt.Sprint(*port), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
