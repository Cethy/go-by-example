module http-proxy-ssl

go 1.23.1

replace http-server-middleware => ../http-server-middleware

replace http-proxy => ../http-proxy

require (
	http-proxy v0.0.0-00010101000000-000000000000
	http-server-middleware v0.0.0-00010101000000-000000000000
)
