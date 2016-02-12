// Package dashboard provides the administrative and exploratory web interface.
package dashboard

//go:generate go-bindata -pkg dashboard -o assets.go --prefix assets assets/...

import (
	"net/http"

	"github.com/etherapis/etherapis/etherapis/geth"
)

// New creates an HTTP route multiplexer injected with all the various components
// required to run the dashboard: static assets and API endpoints.
func New(ethereum *geth.API, assetsPath string) *http.ServeMux {
	router := http.NewServeMux()

	// Register the static asset handler
	if assetsPath != "" {
		router.Handle("/", http.FileServer(http.Dir(assetsPath)))
	} else {
		router.HandleFunc("/", handleAsset)
	}
	// Register the various API handlers
	router.Handle("/api/v0/", newAPIServeMux("/api/v0/", ethereum))

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
		w.Write(data)
		return
	}
	http.NotFound(w, r)
}
