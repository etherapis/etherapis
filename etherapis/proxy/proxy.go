// Package proxy implements the HTTP payment proxy between a locally exposed
// endpoint and the public internet.
package proxy

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
)

// ProxyType is the various types of proxies that can be created.
type ProxyType int

const (
	CallProxy ProxyType = iota // Payment units are authorized and charged per API call
	DataProxy                  // Payment units are authorized and charged per data traffic
)

// Proxy is a payment gateway between arbitrarily many internal services and
// the outside world. Its role is to broker API requests between them, while
// at the same time enforcing payment authorizations.
type Proxy struct {
	extPort int          // External port number to accept requests on
	intPort int          // Internal port to forward requests to
	kind    ProxyType    // Proxy payment authorization type
	logger  log15.Logger // ID-embedded contextual logger
}

// New creates a new payment proxy between an internal and external world.
func New(id int, extPort, intPort int, kind ProxyType) *Proxy {
	return &Proxy{
		extPort: extPort,
		intPort: intPort,
		logger:  log15.New("proxy-id", id),
	}
}

// Start boots up the proxy, opening up the HTTP listeners towards the internally
// available service.
//
// Note, the method will block forever.
func (p *Proxy) Start() error {
	p.logger.Info("Starting up proxy", "external-port", p.extPort, "internal-port", p.intPort)
	return http.ListenAndServe(fmt.Sprintf(":%d", p.extPort), p)
}

// ServeHTTP implements http.Handler, extracting and validating payment headers
// contained within the HTTP request. If payment information is accepted, the
// request is passed on to the internal service for execution. Otherwise the proxy
// short circuits the request and sends back an appropriate error.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Retrieve all the headers from the original request
	headers := r.Header
	var (
		sub = headers.Get(SubscriptionHeader)
		sum = headers.Get(AuthorizationHeader)
		sig = headers.Get(SignatureHeader)
	)
	p.logger.Debug("Received an API request", "subscription", sub, "authorized", sum, "signature", sig)

	// Ensure that all payment information are present
	if sub == "" {
		http.Error(w, "Missing HTTP header: "+SubscriptionHeader, http.StatusBadRequest)
		return
	}
	if sum == "" {
		http.Error(w, "Missing HTTP header: "+AuthorizationHeader, http.StatusBadRequest)
		return
	}
	if sig == "" {
		http.Error(w, "Missing HTTP header: "+SignatureHeader, http.StatusBadRequest)
		return
	}
	// Read the entire request and relay it into the internal API
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read input request", http.StatusPartialContent)
		return
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", p.intPort, r.URL.String()), bytes.NewReader(body))
	if err != nil {
		p.logger.Error("Failed to assemble internal HTTP request", "error", err)
		http.Error(w, "Failed to execute request", http.StatusInternalServerError)
		return
	}
	req.Header = headers
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		p.logger.Error("Failed to process inernal HTTP request", "error", err)
		http.Error(w, "Failed to execute request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	for key, values := range res.Header {
		w.Header().Set(key, values[0])
	}
	io.Copy(w, res.Body)
}
