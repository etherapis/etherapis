// Contains the etherapis RESTful HTTP API endpoint.

package dashboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/etherapis/etherapis/etherapis"
	"github.com/etherapis/etherapis/etherapis/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"gopkg.in/inconshreveable/log15.v2"
)

// api is a wrapper around the etherapis components, exposing various parts of
// each submodule.
type api struct {
	eapis *etherapis.EtherAPIs
}

// newAPIServeMux creates an etherapis API endpoint to serve RESTful requests,
// and returns the HTTP route multipelxer to embed in the main handler.
func newAPIServeMux(base string, eapis *etherapis.EtherAPIs) *mux.Router {
	// Create an API to expose various internals
	handler := &api{
		eapis: eapis,
	}
	// Register all the API handler endpoints
	router := mux.NewRouter()

	router.HandleFunc(base+"accounts", handler.Accounts)
	router.HandleFunc(base+"accounts/{address:0x[0-9a-f]{40}}", handler.Account)
	router.HandleFunc(base+"accounts/{address:0x[0-9a-f]{40}}/{password:.+}", handler.AccountExport)
	router.HandleFunc(base+"services/{addresses}", handler.Services)
	router.HandleFunc(base+"services", handler.Services)
	router.HandleFunc(base+"subscriptions/{address}", handler.Subscriptions)

	return router
}

// Services returns the services for a given address or all services if
// no list of address is given.
func (a *api) Services(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		vars = mux.Vars(r)
	)
	services := make([]contract.Service, 0) // Initialize to serialize into [], not null

	// if there's an address present on the URL return the services
	// owned by this account.
	if addresses, exist := vars["addresses"]; exist {
		// addresses is a comma separated list of addresseses
		for _, addr := range strings.Split(addresses, ",") {
			srvs, err := a.eapis.Contract().Services(common.HexToAddress(addr))
			if err != nil {
				log15.Error("Failed to retrieve services", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			services = append(services, srvs...)
		}
	} else {
		// Get all services
		services, err = a.eapis.Contract().AllServices()
		if err != nil {
			log15.Error("Failed to retrieve services", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

	services, err := a.eapis.Contract().Subscriptions(common.HexToAddress(addr))
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
// payment proxies and/or the marketplace and subscriptions. In case of a HTTP POST,
// a new account is imported using the uploaded key file and access password.
func (a *api) Accounts(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST" && r.FormValue("action") == "create":
		// Create a brand new random account
		address, err := a.eapis.CreateAccount()
		if err != nil {
			log15.Error("Failed to generate account", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("0x%x", address)))

	case r.Method == "POST" && r.FormValue("action") == "import":
		// Import a new account specified by a form
		password := r.FormValue("password")
		account, _, err := r.FormFile("account")
		if err != nil {
			log15.Error("Failed to retrieve account", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		key, err := ioutil.ReadAll(account)
		if err != nil {
			log15.Error("Failed to read account", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		address, err := a.eapis.ImportAccount(key, password)
		if err != nil {
			log15.Error("Failed to import account", "error", err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.Write([]byte(fmt.Sprintf("0x%x", address)))

	default:
		http.Error(w, "Unsupported method: "+r.Method, http.StatusMethodNotAllowed)
	}
}

// Account is the RESTfull endpoint for account management.
func (a *api) Account(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Only reply for deletion requests and get request
	switch r.Method {
	case "DELETE":
		// Delete the account and return an error if something goes wrong
		if err := a.eapis.DeleteAccount(common.HexToAddress(params["address"])); err != nil {
			log15.Error("Failed to delete account", "address", params["address"], "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		log15.Error("Invalid method on account endpoint", "method", r.Method)
		http.Error(w, "Unsupported method: "+r.Method, http.StatusMethodNotAllowed)
		return
	}
}

// AccountExport handles the request to export an account with a particular
// password set.
func (a *api) AccountExport(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Export the key into a json key file and return and error if something goes wrong
	key, err := a.eapis.ExportAccount(common.HexToAddress(params["address"]), params["password"])
	if err != nil {
		log15.Error("Failed to export account", "address", params["address"], "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Pretty print the key json since the exporter sucks :P
	pretty := new(bytes.Buffer)
	if err := json.Indent(pretty, key, "", "  "); err != nil {
		log15.Error("Failed to pretty print key", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set the correct header to ensure download (i.e. no display) and dump the contents
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "inline; filename=\""+params["address"]+".json\"")
	w.Write(pretty.Bytes())
}
