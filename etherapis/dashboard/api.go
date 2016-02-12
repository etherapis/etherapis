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

// newAPIServeMux creates an etherapis API endpoint to serve RESTful requests,
// and returns the HTTP route multipelxer to embed in the main handler.
func newAPIServeMux(base string, ethereum *geth.API) *http.ServeMux {
	// Create an API to expose various internals
	handler := &api{
		ethereum: ethereum,
	}
	// Register all the API handler endpoints
	router := http.NewServeMux()

	router.HandleFunc(base+"accounts", handler.Accounts)
	router.HandleFunc(base+"ethereum/peers", handler.PeerInfos)
	router.HandleFunc(base+"ethereum/syncing", handler.SyncStatus)
	router.HandleFunc(base+"ethereum/head", handler.HeadBlock)

	return router
}

// Accounts retrieves the accounts currently owned by the node.
func (a *api) Accounts(w http.ResponseWriter, r *http.Request) {
	reply, err := a.ethereum.Request("eth_accounts", nil)
	if err != nil {
		log15.Error("Failed to retrieve owned accounts", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// PeerInfos retrieves the currently connected peers and returns them in their
// raw Ethereum API reply form.
func (a *api) PeerInfos(w http.ResponseWriter, r *http.Request) {
	reply, err := a.ethereum.Request("admin_peers", nil)
	if err != nil {
		log15.Error("Failed to retrieve connected peers", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// SyncStatus retrieves the current sync status and returns it in its raw
// Ethereum API reply form.
func (a *api) SyncStatus(w http.ResponseWriter, r *http.Request) {
	reply, err := a.ethereum.Request("eth_syncing", nil)
	if err != nil {
		log15.Error("Failed to retrieve sync status", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}

// HeadBlock retrieves the current head block and returns it in its raw
// Ethereum API reply form.
func (a *api) HeadBlock(w http.ResponseWriter, r *http.Request) {
	reply, err := a.ethereum.Request("eth_getBlockByNumber", []interface{}{"latest", false})
	if err != nil {
		log15.Error("Failed to retrieve head block", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}
