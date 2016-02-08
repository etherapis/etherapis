// This demo app is a simple video streamer that serves up content over an HTTP
// socket. An initial set of data assets is downloaded from external sources to
// use for the demo on startup. This can be done without serving too via --init,
// which can be useful for bundling into self contained docker containers.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Asset URLs to load into the streamer data directory.
	assets = []string{
		"https://github.com/etherapis/webmcoder/releases/download/v0.1/elephants-dream.webm",
	}
)

var (
	initFlag = flag.Bool("init", false, "Initialize the streamer with any external assets, don't run")
	portFlag = flag.Int("port", 8081, "Port number on which to listen for HTTP data requests")
	rootFlag = flag.String("root", "", "Document root from where to stream files")
)

// initialize iterates over the external assets and downloads them into the server
// data directory if they do not exist yet.
func initialize(docroot string, assets []string) error {
	// Iterate over and download all the missing assets
	for _, asset := range assets {
		// If the file already exists, skip
		parts := strings.Split(asset, "/")
		name := parts[len(parts)-1]
		path := filepath.Join(docroot, name)

		if _, err := os.Stat(path); err == nil {
			continue
		}
		// Otherwise download the remote asset
		log.Printf("Downloading missing asset %s...", name)
		output, err := os.Create(path)
		if err != nil {
			return err
		}
		defer output.Close()

		response, err := http.Get(asset)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		n, err := io.Copy(output, response.Body)
		if err != nil {
			return err
		}
		log.Printf("Downloaded %s: %d bytes.", name, n)
	}
	return nil
}

func main() {
	// Process and validate any command line flags
	flag.Parse()

	if *rootFlag == "" {
		log.Fatalf("Please specify a path to serve files from (--root)")
	}
	if err := os.MkdirAll(*rootFlag, 0640); err != nil {
		log.Fatalf("Failed to create specified document root: %v", err)
	}
	// Download the remote assets into the document root
	if err := initialize(*rootFlag, assets); err != nil {
		log.Fatalf("Failed to download remote assets: %v", err)
	}
	if *initFlag {
		// Only init requested, no actual runtime, stop
		return
	}
	// Start data streaming service
	log.Printf("Starting file streaming on port %d...\n", *portFlag)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), http.FileServer(http.Dir(*rootFlag))); err != nil {
		log.Fatalf("Failed to start file streamer: %v", err)
	}
}
