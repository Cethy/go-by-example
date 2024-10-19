---
Order: 9
Title: Static Website Generator
Summary: Opinionated Website Generator (powering this website)
ImgSrc: https://images.unsplash.com/photo-1611647832580-377268dba7cb?ixid=M3w2NjYzMTJ8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MjkyNzkxOTh8&ixlib=rb-4.0.3
---

# Static Website Generator

## Instructions

Make a program capable of generating this very website.

It must be capable of:

- copying dependency files (images, css files, js files, ...)
- assemble html files from fragments of repeated components like header or footer.
- generate article html files based on the projects README files built at this time
- provides some QoL scripts for developers, such as
  - generator
  - serve
  - watch
- the main component (generator) should be easily extendable

## Key Features

- generator extensibility
- fragment transformer
- article transformer
- article validator (ensure mandatory metadata are present and add a random unsplash cover image if none is present)

## Usage

```shell
# go get github.com/fsnotify/fsnotify
# go get github.com/yuin/goldmark-meta

# build
go run static-website-generator/generator.go
# serve
go run static-website-generator/serve.go
```

## TODO
- [ ] use cobra
