// etherapis implements the Ether APIs marketplace gateway.
package etherapis

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/etherapis/etherapis/etherapis/contract"
	"github.com/etherapis/etherapis/etherapis/geth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"gopkg.in/inconshreveable/log15.v2"
)

// EtherAPIs is the main logic behind the entire system.
type EtherAPIs struct {
	client   *geth.Geth         // Embedded Ethereum client
	ethereum *eth.Ethereum      // Actual Ethereum protocol within the client
	rpcapi   *geth.API          // In-process RPC interface to the embedded client
	password string             // Master password to use to handle local accounts
	contract *contract.Contract // Ethereum contract handling consensus stuff
}

// New creates a new Ether APIs instance, connects with it to the Ethereum network
// via an embedded Geth instance and attaches an RPC in-process endpoint to it.
func New(datadir string, network geth.EthereumNetwork, address common.Address) (*EtherAPIs, error) {
	// Create a Geth instance and boot it up
	client, err := geth.New(datadir, network)
	if err != nil {
		return nil, err
	}
	if err := client.Start(); err != nil {
		return nil, err
	}
	// Retrieve the underlying ethereum service and attach global RPC interface
	var ethereum *eth.Ethereum
	if err := client.Stack().Service(&ethereum); err != nil {
		return nil, err
	}
	api, err := client.Attach()
	if err != nil {
		return nil, err
	}
	// Assemble an interface around the consensus contract
	contract, err := contract.New(ethereum.ChainDb(), ethereum.EventMux(), ethereum.BlockChain(), ethereum.Miner().PendingState)
	if err != nil {
		return nil, err
	}
	// Assemble and return the Ether APIs instance
	return &EtherAPIs{
		client:   client,
		ethereum: ethereum,
		rpcapi:   api,
		contract: contract,
	}, nil
}

// Close terminates the EtherAPIs instance along with all held resources.
func (eapis *EtherAPIs) Close() error {
	return eapis.client.Stop()
}

// Unlock iterates over all the known accounts and tries to unlock them using the
// provided master password.
func (eapis *EtherAPIs) Unlock(password string) error {
	// Retrieve the list of known accounts
	manager := eapis.ethereum.AccountManager()

	accounts, err := manager.Accounts()
	if err != nil {
		return err
	}
	// Unlock each of them using the master password
	for _, account := range accounts {
		address := fmt.Sprintf("0x%x", account.Address)

		log15.Debug("Unlocking account...", "account", address)
		if err := manager.Unlock(account.Address, password); err != nil {
			return fmt.Errorf("password rejected for %s", address)
		}
	}
	// All accounts unlocked successfully, accept the master password
	eapis.password = password

	return nil
}

// CreateAccount generates a new random account and returns it.
func (eapis *EtherAPIs) CreateAccount() (common.Address, error) {
	account, err := eapis.ethereum.AccountManager().NewAccount(eapis.password)
	if err != nil {
		return common.Address{}, err
	}
	if err := eapis.ethereum.AccountManager().Unlock(account.Address, eapis.password); err != nil {
		panic(fmt.Sprintf("Newly created account failed to unlock: %v", err))
	}
	return account.Address, err
}

// ImportAccount inserts an encrypted external account into the local keystore
// by first decrypting it, and then inserting using the local master password.
func (eapis *EtherAPIs) ImportAccount(keyjson []byte, password string) (common.Address, error) {
	key, err := crypto.DecryptKey(keyjson, password)
	if err != nil {
		return common.Address{}, err
	}
	if eapis.ethereum.AccountManager().HasAccount(key.Address) {
		return common.Address{}, errors.New("account already exists")
	}
	if err := eapis.client.Keystore().StoreKey(key, eapis.password); err != nil {
		return common.Address{}, err
	}
	if err := eapis.ethereum.AccountManager().Unlock(key.Address, eapis.password); err != nil {
		panic(fmt.Sprintf("Newly imported account failed to unlock: %v", err))
	}
	return key.Address, nil
}

// ExportAccount retrieves an account from the key store and exports it using
// a different password.
func (eapis *EtherAPIs) ExportAccount(account common.Address, password string) ([]byte, error) {
	key, err := eapis.client.Keystore().GetKey(account, eapis.password)
	if err != nil {
		return nil, err
	}
	return crypto.EncryptKey(key, password, crypto.StandardScryptN, crypto.StandardScryptP)
}

// DeleteAccount irreversibly deletes an account from the key store.
func (eapis *EtherAPIs) DeleteAccount(account common.Address) error {
	return eapis.ethereum.AccountManager().DeleteAccount(account, eapis.password)
}

// Account represents an ethereum account.
type Account struct {
	Nonce   uint64 `json:"nonce"`
	Balance string `json:"balance"`
}

// GetAccount returns the data associated with the account.
func (eapis *EtherAPIs) GetAccount(account common.Address) Account {
	state := eapis.ethereum.Miner().PendingState()
	return Account{Nonce: state.GetNonce(account), Balance: state.GetBalance(account).String()}
}

// Accounts retrieves the list of accounts known to etherapis.
func (eapis *EtherAPIs) Accounts() ([]common.Address, error) {
	accounts, err := eapis.ethereum.AccountManager().Accounts()
	if err != nil {
		return nil, err
	}
	addresses := make([]common.Address, len(accounts))
	for i, account := range accounts {
		addresses[i] = account.Address
	}
	return addresses, nil
}

// Contract retrieves the Ether APIs Ethereum contract to access the consensus data.
func (eapis *EtherAPIs) Contract() *contract.Contract {
	return eapis.contract
}

// CallRPC is a temporary helper method to pass an RPC call to the underlying
// go-ethereum server. It returns the exact raw response, no parsing done.
func (eapis *EtherAPIs) CallRPC(method string, params []interface{}) (json.RawMessage, error) {
	return eapis.rpcapi.Request(method, params)
}
