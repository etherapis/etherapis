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
	http.HandleFunc("/faucet", faucet)
	http.HandleFunc("/faucet/fund", faucetFund)
}

// Account is an pre-funded Ethereum demo account for the EtherAPIs demos.
type Account struct {
	Key  string    // Private key file in either the .js form or the Geth encoded key
	Used time.Time // Last time assigned (reuse accounts to prevent someone exhausting them)
}

// accountKey returns the key used for all account entries.
func accountKey(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "Account", "testnet", 0, nil)
}

// faucet retrieves an Ethereum testnet private key pre-loaded with some initial
// funds.
func faucet(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Fetch the least recently used account from the database
	query := datastore.NewQuery("Account").Ancestor(accountKey(ctx)).Order("Used").Limit(1)

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

// faucetFund uploads a new demo account into the faucet server, serving that
// too from this point onward.
func faucetFund(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	acc := Account{
		Key:  r.FormValue("key"),
		Used: time.Now(),
	}
	key := datastore.NewIncompleteKey(ctx, "Account", accountKey(ctx))
	if _, err := datastore.Put(ctx, key, &acc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "funded")
}
