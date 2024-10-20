module http-proxy-cache

go 1.23.1

replace http-proxy => ../http-proxy

replace http-server-middleware => ../http-server-middleware

require (
	http-proxy v0.0.0-00010101000000-000000000000
	http-server-middleware v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
)
