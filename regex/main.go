package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {
	var message string
	path := flag.String("path", "", "pathname to a src file")
	flag.StringVar(&message, "m", "", "message to regex")
	flag.Parse()

	if *path != "" {
		content, err := os.ReadFile(*path)
		if err != nil {
			panic(err)
		}
		message = string(content)
	}

	pattern := "(\\{[A-z]+\\})"
	r, _ := regexp.Compile(pattern)
	all := r.FindAllString(message, -1)
	for _, item := range all {
		fmt.Println(item)
	}
}
