package main

import (
	"flag"
	"http2curl/impl"
)

func main() {
	var webPattern bool
	flag.BoolVar(&webPattern, "w", false, "use the web pattern. open 127.0.0.1:4877 in your browser")
	flag.Parse()
	if webPattern {
		impl.Web()
	} else {
		impl.Cli()
	}
}
