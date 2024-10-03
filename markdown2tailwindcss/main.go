package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	goldmarkTailwindcss "go-by-example/libs/goldmark-tailwindcss"
	"log"
	"net/http"
	"strings"
)

var mdConverter = goldmark.New(
	goldmark.WithExtensions(
		meta.Meta,
	),
	goldmark.WithRenderer(
		renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(html.NewRenderer(), 1000))),
	),
	goldmark.WithParser(
		parser.NewParser(parser.WithBlockParsers(parser.DefaultBlockParsers()...),
			parser.WithInlineParsers(parser.DefaultInlineParsers()...),
			parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
			parser.WithASTTransformers(
				util.Prioritized(goldmarkTailwindcss.NewTransformer(), 100),
			),
		),
	),
)

var mdContent = []byte(`# h1
## h2
### h3
#### h4
##### h5
###### h6

para **bold** *italic* [link](http://example.com) 
graph

- list item 1
- list item 2
    - list item 2.1
    - list item 2.2
- list item 3`)

var indexFile = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>markdown2tailwindcss</title>
    <link rel="icon" type="image/x-icon" href="/favicon.ico">
    <script src="/public/tailwindcss-3.4.5.min.js"></script>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap" rel="stylesheet">
    <style>
        :root {
            font-family: 'Inter', sans-serif;
            color: white;
        }
    </style>
</head>

<body class="bg-gradient-to-b from-slate-800 to-slate-950 antialiased px-4 w-full max-w-[48rem] mx-auto">
<section class="w-full min-h-[calc(100vh)]">
    <article class="flex flex-col gap-4 py-10">
        {content}
    </article>
</section>
</body>
</html>`

func main() {
	port := flag.Int("p", 8008, "Port number")
	flag.Parse()

	var buf bytes.Buffer
	if err := mdConverter.Convert(mdContent, &buf); err != nil {
		panic(err)
	}

	response := strings.ReplaceAll(indexFile, "{content}", string(buf.Bytes()))

	//http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)

		fmt.Fprintf(w, response)
	})

	publicDir := http.FileServer(http.Dir("public/"))
	http.Handle("/public/", http.StripPrefix("/public/", publicDir))
	http.Handle("/favicon.ico", publicDir)

	fmt.Println("Server listening on port:", *port)
	err := http.ListenAndServe(":"+fmt.Sprint(*port), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
