module static-website-generator

go 1.23.1

replace markdown2tailwindcss => ../markdown2tailwindcss

require (
	github.com/fsnotify/fsnotify v1.7.0
	github.com/hbagdi/go-unsplash v0.0.0-20230414214043-474fc02c9119
	github.com/spf13/cobra v1.8.1
	github.com/yuin/goldmark v1.7.8
	github.com/yuin/goldmark-meta v1.1.0
	golang.org/x/oauth2 v0.23.0
	markdown2tailwindcss v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.4.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
