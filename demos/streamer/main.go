// This demo app is a simple video streamer that serves up content on an HTTP
// socket.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	rootFlag = flag.String("root", "", "Document root from where to stream files")
	portFlag = flag.Int("port", 8081, "Port number on which to listen for HTTP requests")
)

func main() {
	flag.Parse()

	if *rootFlag == "" {
		log.Fatalf("Please specify a path to serve files from (--root)")
	}
	http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), http.FileServer(http.Dir(*rootFlag)))
}
