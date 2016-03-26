package etherapis

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Account represents an Ethereum account.
type Account struct {
	Nonce        uint64         `json:"nonce"`
	Balance      *big.Int       `json:"balance"`
	Change       *big.Int       `json:"change"`
	Transactions []*Transaction `json:"transactions"`
}

// Transaction represents an Ethereum transaction.
type Transaction struct {
	Hash   common.Hash    `json:"hash"`
	From   common.Address `json:"from"`
	To     common.Address `json:"to"`
	Amount *big.Int       `json:"amount"`
	Fees   *big.Int       `json:"fees"`
}

// Service represents an Ether APIs service created by the contract.
type Service struct {
	ID           *big.Int       `json:"id"`           // Globally unique identifier for the service
	Name         string         `json:"name"`         // Name assigned to the service by its provider
	Owner        common.Address `json:"owner"`        // Address of the service provider
	Endpoint     string         `json:"endpoint"`     // Endpoint and/or website to reach the service
	Model        *big.Int       `json:"model"`        // Payment model to charge based on (0 = call, 1 = data, 2 = time)
	Price        *big.Int       `json:"price"`        // Price per unit (defined by the payment model)
	Cancellation *big.Int       `json:"cancellation"` // Minimum time before unused funds are released
	Enabled      bool           `json:"enabled"`      // Whether the contract accepts subscriptions or not

	Creating bool `json:"creating"` // Whether the ervice registration is currently being executed
	Changing bool `json:"changing"` // Whether the contract enabled/disabled state is currently changing
	Deleting bool `json:"deleting"` // Whether the contract is currently being deleted
}
