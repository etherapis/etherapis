// This demo app is a relatively simple Geo IP lookup service, which actually
// blatantly piggie-backs freegeoip.net for the IP lookups.

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

var (
	portFlag = flag.Int("port", 8080, "Port number on which to listen for HTTP requests")
)

func init() {
	http.HandleFunc("/", lookup)
}

// lookup is an HTTP handler that parses an IP address from the URL parameters,
// queries freegeoip for the geographical data and returns it to the user.
func lookup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Query freegeoip for the address lookup
	address := r.URL.Query().Get("ip")

	res, err := http.Get("https://freegeoip.net/json/" + address)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to resolve IP address: %v", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Stream any output they send to our requester
	io.Copy(w, res.Body)
}

func main() {
	flag.Parse()
	http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), nil)
}
