// Contains the etherapis RESTful HTTP API endpoint.

package dashboard

import (
	"net/http"

	"github.com/etherapis/etherapis/etherapis/geth"
	"gopkg.in/inconshreveable/log15.v2"
)

// api is a wrapper around the etherapis components, exposing various parts of
// each submodule.
type api struct {
	ethereum *geth.API
}

// newApi creates an etherapis API endpoint to serve RESTful requests.
func newApi(ethereum *geth.API) *api {
	return &api{
		ethereum: ethereum,
	}
}

// PeersHandler retrieves the currently connected peers and returns them in their
// raw Ethereum API reply form.
func (a *api) PeersHandler(w http.ResponseWriter, r *http.Request) {
	reply, err := a.ethereum.Request("admin_peers", nil)
	if err != nil {
		log15.Error("Failed to retrieve connected peers", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// SyncingHandler retrieves the current sync status and returns it in its raw
// Ethereum API reply form.
func (a *api) SyncingHandler(w http.ResponseWriter, r *http.Request) {
	reply, err := a.ethereum.Request("eth_syncing", nil)
	if err != nil {
		log15.Error("Failed to retrieve sync status", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// HeadHandler retrieves the current head block and returns it in its raw
// Ethereum API reply form.
func (a *api) HeadHandler(w http.ResponseWriter, r *http.Request) {
	reply, err := a.ethereum.Request("eth_getBlockByNumber", []interface{}{"latest", false})
	if err != nil {
		log15.Error("Failed to retrieve head block", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}
