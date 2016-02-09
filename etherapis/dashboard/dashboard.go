// Package dashboard provides the administrative and exploratory web interface.
package dashboard

import (
	"net/http"

	"github.com/etherapis/etherapis/etherapis/geth"
	"github.com/gorilla/mux"
	"gopkg.in/inconshreveable/log15.v2"
)

// New creates an HTTP route multiplexer injected with all the various components
// required to run the dashboard: static assets and API endpoints.
func New(ethereum *geth.API) *mux.Router {
	// Create an API to expose various internals
	api := &api{
		ethereum: ethereum,
	}
	// Register all the route handlers
	router := mux.NewRouter()
	router.HandleFunc("/api/ethereum/peers", api.PeersHandler)

	return router
}

// api is a wrapper around the etherapis components, exposing various parts of
// each submodule.
type api struct {
	ethereum *geth.API
}

// PeersHandler retrieves the currently connected peers and returns them.
func (a *api) PeersHandler(w http.ResponseWriter, r *http.Request) {
	reply, err := a.ethereum.Request("admin_peers", nil)
	if err != nil {
		log15.Error("Failed to retrieve connected peers", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}
