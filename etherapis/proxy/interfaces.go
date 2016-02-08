// Contains the interfaces for payment authorization verification.

package proxy

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Verifier is an interface that accepts the details of a payment authorization
// and returns whether the sender is allowed to make the payment or not.
type Verifier interface {
	// Exists checks whether there's a live payment channel already set up between
	// the sender and recipient.
	Exists(from, to common.Address) bool

	// Verify checks whether the authorization is cryptographically valid, and also
	// whether there are enough funds in the payment channel to process this payment.
	Verify(from, to common.Address, nonce uint64, amount *big.Int, signature []byte) (bool, bool)

	// Price returns the price provided by the signature of (from || to).
	Price(from, to common.Address) *big.Int
	// Nonce returns the nonce provided by the signature of (from || to).
	Nonce(from, to common.Address) *big.Int
}

// Charger chaaaaarges! :D Fun's aside, this interfaces provides the capability
// to redeem an authorized payment by the underlying framework.
type Charger interface {
	// Charge calls down into the underlying Ethereum contract layer and executes
	// a payment charging transaction. It returns the hex encoded transaction id
	// to enable later verification.
	Charge(from, to common.Address, nonce uint64, amount *big.Int, signature []byte) (common.Hash, error)
}
