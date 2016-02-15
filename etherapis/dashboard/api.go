// Contains the etherapis RESTful HTTP API endpoint.

package dashboard

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/etherapis/etherapis/etherapis/channels"
	"github.com/etherapis/etherapis/etherapis/geth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/gorilla/mux"
	"gopkg.in/inconshreveable/log15.v2"
)

// api is a wrapper around the etherapis components, exposing various parts of
// each submodule.
type api struct {
	ethereum *eth.Ethereum
	gethAPI  *geth.API
	contract *channels.Subscriptions
}

// newAPIServeMux creates an etherapis API endpoint to serve RESTful requests,
// and returns the HTTP route multipelxer to embed in the main handler.
func newAPIServeMux(base string, contract *channels.Subscriptions, ethereum *eth.Ethereum, gethAPI *geth.API) *mux.Router {
	// Create an API to expose various internals
	handler := &api{
		ethereum: ethereum,
		gethAPI:  gethAPI,
		contract: contract,
	}
	// Register all the API handler endpoints
	router := mux.NewRouter()

	router.HandleFunc(base+"accounts", handler.Accounts)
	router.HandleFunc(base+"ethereum/peers", handler.PeerInfos)
	router.HandleFunc(base+"ethereum/syncing", handler.SyncStatus)
	router.HandleFunc(base+"ethereum/head", handler.HeadBlock)
	router.HandleFunc(base+"services/{addresses}", handler.Services)
	router.HandleFunc(base+"subscriptions/{address}", handler.Subscriptions)

	return router
}

// Services retrieves the given address' services.
func (a *api) Services(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addresses, exist := vars["addresses"]
	if !exist {
		log15.Error("Failed to retrieve services", "error", "no address specified")
		http.Error(w, "no address specified", http.StatusInternalServerError)
		return
	}

	var services []channels.Service
	// addresses is a comma separated list of addresseses
	for _, addr := range strings.Split(addresses, ",") {
		srvs, err := a.contract.Services(common.HexToAddress(addr))
		if err != nil {
			log15.Error("Failed to retrieve services", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		services = append(services, srvs...)
	}

	out, err := json.Marshal(services)
	if err != nil {
		log15.Error("Failed to retrieve services", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(out)
}

// Subscriptions retrieves the given address' subscriptions.
func (a *api) Subscriptions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr, exist := vars["address"]
	if !exist {
		log15.Error("Failed to retrieve subscriptions", "error", "no address specified")
		http.Error(w, "no address specified", http.StatusInternalServerError)
		return
	}

	services, err := a.contract.Subscriptions(common.HexToAddress(addr))
	if err != nil {
		log15.Error("Failed to retrieve subscriptions", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(services)
	if err != nil {
		log15.Error("Failed to marshal subscriptions", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(out)
}

// Accounts retrieves the Ethereum accounts currently configured to be used by the
// payment proxies and/or the marketplace and subscriptions.
func (a *api) Accounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := a.ethereum.AccountManager().Accounts()
	if err != nil {
		log15.Error("Failed to retrieve accounts", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out, err := json.Marshal(accounts)
	if err != nil {
		log15.Error("Failed to marshal account list", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(out)
}

// PeerInfos retrieves the currently connected peers and returns them in their
// raw Ethereum API reply form.
func (a *api) PeerInfos(w http.ResponseWriter, r *http.Request) {
	reply, err := a.gethAPI.Request("admin_peers", nil)
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
	reply, err := a.gethAPI.Request("eth_syncing", nil)
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
	reply, err := a.gethAPI.Request("eth_getBlockByNumber", []interface{}{"latest", false})
	if err != nil {
		log15.Error("Failed to retrieve head block", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(reply)
}
