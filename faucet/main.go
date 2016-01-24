// Package faucet implements an Ethereum test network account faucet.
package faucet

import (
	"fmt"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
)

func init() {
	http.HandleFunc("/faucet/web", webFaucet)
	http.HandleFunc("/faucet/cli", cliFaucet)
	http.HandleFunc("/faucet/web/fund", webFaucetFund)
	http.HandleFunc("/faucet/cli/fund", cliFaucetFund)
}

// Account is an pre-funded Ethereum demo account for the EtherAPIs demos.
type Account struct {
	Key  string    // Private key file in either the .js form or the Geth encoded key
	Used time.Time // Last time assigned (reuse accounts to prevent someone exhausting them)
}

// webAccountKey returns the key used for all web account entries.
func webAccountKey(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "Account", "web", 0, nil)
}

// cliAccountKey returns the key used for all cli account entries.
func cliAccountKey(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "Account", "cli", 0, nil)
}

// webFaucet retrieves a JavaScript Ethereum account for browser web use and
// sends in to the user.
func webFaucet(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Fetch the least recently used account from the database
	query := datastore.NewQuery("Account").Ancestor(webAccountKey(ctx)).Order("Used").Limit(1)

	accounts := make([]Account, 0, 1)
	keys, err := query.GetAll(ctx, &accounts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Update the chosen account's usage time
	accounts[0].Used = time.Now()
	if _, err := datastore.Put(ctx, keys[0], &accounts[0]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the account to the user
	fmt.Fprint(w, accounts[0].Key)
}

// cliFaucet retrieves a Geth Ethereum account for console use and sends in to
// the user.
func cliFaucet(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Fetch the least recently used account from the database
	query := datastore.NewQuery("Account").Ancestor(cliAccountKey(ctx)).Order("Used").Limit(1)

	accounts := make([]Account, 0, 1)
	keys, err := query.GetAll(ctx, &accounts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Update the chosen account's usage time
	accounts[0].Used = time.Now()
	if _, err := datastore.Put(ctx, keys[0], &accounts[0]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the account to the user
	fmt.Fprint(w, accounts[0].Key)
}

// webFaucetFund uploads a new website demo account into the faucet server,
// serving that too from this point onward.
func webFaucetFund(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	acc := Account{
		Key:  r.FormValue("key"),
		Used: time.Now(),
	}
	key := datastore.NewIncompleteKey(ctx, "Account", webAccountKey(ctx))
	if _, err := datastore.Put(ctx, key, &acc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "funded")
}

// cliFaucetFund uploads a new console demo account into the faucet server,
// serving that too from this point onward.
func cliFaucetFund(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	acc := Account{
		Key:  r.FormValue("key"),
		Used: time.Now(),
	}
	key := datastore.NewIncompleteKey(ctx, "Account", cliAccountKey(ctx))
	if _, err := datastore.Put(ctx, key, &acc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "funded")
}
