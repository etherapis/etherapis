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
	ID           *big.Int       `json:"id"`
	Name         string         `json:"name"`
	Owner        common.Address `json:"owner"`
	Endpoint     string         `json:"endpoint"`
	Price        *big.Int       `json:"price"`
	Cancellation *big.Int       `json:"cancellation"`
	Enabled      bool           `json:"enabled"`
	Deleted      bool           `json:"deleted"`
}
