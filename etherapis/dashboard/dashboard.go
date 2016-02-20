// Package dashboard provides the administrative and exploratory web interface.
package dashboard

//go:generate go-bindata -pkg dashboard -o assets.go --prefix assets assets/...

import (
	"net/http"
	"strings"

	"github.com/etherapis/etherapis/etherapis"
)

// New creates an HTTP route multiplexer injected with all the various components
// required to run the dashboard: static assets and API endpoints.
func New(eapis *etherapis.EtherAPIs, assetsPath string) *http.ServeMux {
	router := http.NewServeMux()

	// Register the static asset handler
	if assetsPath != "" {
		router.Handle("/", http.FileServer(http.Dir(assetsPath)))
	} else {
		router.HandleFunc("/", handleAsset)
	}
	// Register the various API handlers
	router.Handle("/api/v0/", newAPIServeMux("/api/v0/", eapis))
	router.Handle("/api/v1/", newStateServer("/api/v1/", eapis))

	return router
}

// handleAsset returns static assets from the data built into the binary itself.
func handleAsset(w http.ResponseWriter, r *http.Request) {
	// Extract the file to retrieve from the URL
	path := r.URL.Path[1:]
	if path == "" {
		path = "index.html"
	}
	// Retrieve the asset and return it, or error out
	if data, err := Asset(path); err == nil {
		// Certain file types cause issues in the browser if not tagged correctly
		switch {
		case strings.HasSuffix(path, ".css"):
			w.Header().Set("Content-Type", "text/css")
		}
		// Write the data itself
		w.Write(data)
		return
	}
	http.NotFound(w, r)
}
