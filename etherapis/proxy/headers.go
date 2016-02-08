// Contains the custom HTTP headers defined by the payment proxy.

package proxy

import (
	"encoding/json"

	"gopkg.in/inconshreveable/log15.v2"
)

const (
	AuthorizationHeader = "etherapi-authorization" // Client side payment authorization header
	VerificationHeader  = "etherapi-verification"  // Server side payment verification header
)

// authorization is the data content of a client-to-server payment authorization
// header, based on which the server may check for fund availability
type authorization struct {
	Consumer  string `json:"consumer"`  // API consumer authorizing the payment
	Provider  string `json:"provider"`  // API provider to which to authorize the payment to
	Nonce     uint64 `json:"nonce"`     // The nonce used for channel subscription
	Amount    uint64 `json:"amount"`    // Amount of calls/data to authorize
	Signature string `json:"signature"` // Secp256k1 elliptic curve signature
}

// verification is the data content of the server-to-client payment verification
// header,
type verification struct {
	Unknown    bool   `json:"unknown,omitempty"`    // Flag set when there's no valid subscription
	Authorized uint64 `json:"authorized,omitempty"` // Last successfully authorized payment amount
	Proof      string `json:"proof,omitempty"`      // Proof of the last authorization for client verification
	Need       uint64 `json:"need,omitempty"`       // Amount needed for this call (need = authorized => all ok)
	Nonce      uint64 `json:"nonce"`                // Nonce required for this call
	Error      string `json:"error,omitempty"`      // Error message specifically deemed to developer consumption
}

// Marshal flattens a verification string into a JSON encoded blob.
func (v *verification) Marshal() string {
	blob, err := json.Marshal(v)
	if err != nil {
		log15.Crit("Failed to marshal verification header", "error", err)
	}
	return string(blob)
}
