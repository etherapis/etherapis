// Package proxy implements the HTTP payment proxy between a locally exposed
// endpoint and the public internet.
package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"sync/atomic"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"
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
	extPort int       // External port number to accept requests on
	intPort int       // Internal port to forward requests to
	kind    ProxyType // Proxy payment authorization type

	verifier Verifier // Payment verifier that looks into the Ethereum state to validate a transaction
	vault    *Vault   // Payment vault storing previously accepted authorizations

	logger log15.Logger // ID-embedded contextual logger
	autoid uint32       // Auto ID to assign to the next request (log tracking)
}

// New creates a new payment proxy between an internal and external world.
func New(id int, extPort, intPort int, kind ProxyType, verifier Verifier, vault *Vault) *Proxy {
	return &Proxy{
		extPort:  extPort,
		intPort:  intPort,
		kind:     kind,
		verifier: verifier,
		vault:    vault,
		logger:   log15.New("proxy-id", id),
	}
}

// Start boots up the proxy, opening up the HTTP listeners towards the internally
// available service.
//
// Note, the method will block forever.
func (p *Proxy) Start() error {
	p.logger.Info("Starting up proxy", "external-port", p.extPort, "internal-port", p.intPort, "type", p.kind)
	return http.ListenAndServe(fmt.Sprintf(":%d", p.extPort), p)
}

// ServeHTTP implements http.Handler, extracting and validating payment headers
// contained within the HTTP request. If payment information is accepted, the
// request is passed on to the internal service for execution. Otherwise the proxy
// short circuits the request and sends back an appropriate error.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqLogger := p.logger.New("request-id", atomic.AddUint32(&p.autoid, 1))
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Short circuit CORS pre-flight requests
	if r.Method == "OPTIONS" {
		reqLogger.Debug("Allowing CORS pre-flight request")
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		return
	}
	// Allow head requests through for data APIs to query the content size
	if r.Method == "HEAD" {
		reqLogger.Debug("Allowing data HEAD request")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")

		res, err := p.service(r)
		if err != nil {
			reqLogger.Error("Failed to process API request", "error", err)
			http.Error(w, "Failed to execute request", http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		p.forward(w, res)
		return
	}
	// Retrieve the authorization header from the original request
	authHeader := r.Header.Get(AuthorizationHeader)
	reqLogger.Debug("Received an API request", "authorization", authHeader)

	// Ensure that all payment information are present
	if authHeader == "" {
		p.fail(w, &verification{Error: "Missing HTTP header: " + AuthorizationHeader})
		return
	}
	auth := new(authorization)
	if err := json.Unmarshal([]byte(authHeader), auth); err != nil {
		p.fail(w, &verification{Error: "Invalid authorization header: " + err.Error()})
		return
	}
	// Process the request and payment based on the proxy type
	switch p.kind {
	case CallProxy:
		// Make sure the consumer authorized the payment for this call
		consumer, provider := common.HexToAddress(auth.Consumer), common.HexToAddress(auth.Provider)

		if !p.verifier.Exists(consumer, provider) {
			p.fail(w, &verification{Unknown: true, Error: "Non existent API subscription"})
			return
		}
		valid, funded := p.verifier.Verify(consumer, provider, new(big.Int).SetUint64(auth.Amount), common.Hex2Bytes(auth.Signature))
		if !valid {
			p.fail(w, &verification{Error: "Invalid authorization signature"})
			return
		}
		if !funded {
			p.fail(w, &verification{Error: "Not enough funds available"})
			return
		}
		if prev := p.vault.Fetch(provider, consumer); prev != nil && prev.Amount+1 > auth.Amount {
			p.fail(w, &verification{
				Error:      "Not enough funds authorized",
				Authorized: prev.Amount,
				Proof:      prev.Signature,
				Need:       prev.Amount + 1,
			})
			return
		}
		p.vault.Store(auth)

		// Execute the API internally and proxy the response
		reqLogger.Debug("Payment accepted for API invocation")
		res, err := p.service(r)
		if err != nil {
			reqLogger.Error("Failed to process API request", "error", err)
			http.Error(w, "Failed to execute request", http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		p.forward(w, res)

	case DataProxy:
		// Since we're paying by the data, retrieve the amount first
		res, err := p.service(r)
		if err != nil {
			reqLogger.Error("Failed to process API request", "error", err)
			http.Error(w, "Failed to execute request", http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		// Make sure the user authorized payment for all the requested data
		data := res.ContentLength
		if data > 0 /* TODO */ {
			reqLogger.Debug("Payment accepted for API stream", "data", data)
			p.forward(w, res)
		}
	}
}

// fail sends out a authorization failure event to the client containing all the
// necessary infos in both the header and the body to allow various parsings.
func (p *Proxy) fail(w http.ResponseWriter, result *verification) {
	failure := result.Marshal()

	w.Header().Set(VerificationHeader, failure)
	http.Error(w, failure, http.StatusBadRequest)
}

// service executes the API request in the internal API, and returns the reply,
// which will either be forwarded as is, or charged per data rate.
func (p *Proxy) service(r *http.Request) (*http.Response, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", p.intPort, r.URL.String()), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header = r.Header
	return http.DefaultClient.Do(req)
}

// forward proxies an internal API response into the externl request's writer.
func (p *Proxy) forward(w http.ResponseWriter, res *http.Response) {
	for key, values := range res.Header {
		w.Header().Set(key, values[0])
	}
	io.Copy(w, res.Body)
}
