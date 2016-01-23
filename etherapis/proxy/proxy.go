// Package proxy implements the HTTP payment proxy between a locally exposed
// endpoint and the public internet.
package proxy

import (
	"fmt"
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
	fmt.Fprintln(w, "Authorization approved, forwarding request...")
}
