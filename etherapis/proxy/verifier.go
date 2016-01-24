// Contains the interfaces for payment authorization verification.

package proxy

import "github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"

// Verifier is an interface that accepts the details of a payment authorization
// and returns whether the sender is allowed to make the payment or not.
type Verifier interface {
	// Exists checks whether there's a live payment channel already set up between
	// the sender and recipient.
	Exists(from, to common.Address) bool

	// Verify checks whether the authorization is cryptographically valid, and also
	// whether there are enough funds in the payment channel to process this payment.
	Verify(from, to common.Address, amount uint64, signature []byte) (bool, bool)
}
