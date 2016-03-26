// Package etherapis implements the Ether APIs marketplace gateway.
package etherapis

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/etherapis/etherapis/etherapis/contract"
	"github.com/etherapis/etherapis/etherapis/geth"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/inconshreveable/log15.v2"
)

// EtherAPIs is the main logic behind the entire system.
type EtherAPIs struct {
	client   *geth.Geth     // Embedded Ethereum client
	ethereum *eth.Ethereum  // Actual Ethereum protocol within the client
	eventmux *event.TypeMux // Event multiplexer to announce various happenings
	rpcapi   *geth.API      // In-process RPC interface to the embedded client
	password string         // Master password to use to handle local accounts

	contract *contract.EtherAPIs // Ethereum contract handling the matketplace consensus
	signer   bind.SignerFn       // Signer to authorize transactions on the contract

	txlock sync.Mutex // Serializes transaction creation to avoid nonce collisions
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
	rpcClient, err := client.Stack().Attach()
	if err != nil {
		return nil, err
	}
	contract, err := contract.NewEtherAPIs(address, backends.NewRPCBackend(rpcClient))
	if err != nil {
		return nil, err
	}
	// Assemble and return the Ether APIs instance
	return &EtherAPIs{
		client:   client,
		ethereum: ethereum,
		eventmux: client.Stack().EventMux(),
		rpcapi:   api,
		contract: contract,
		signer: func(from common.Address, tx *types.Transaction) (*types.Transaction, error) {
			signature, err := ethereum.AccountManager().Sign(accounts.Account{Address: from}, tx.SigHash().Bytes())
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signature)
		},
	}, nil
}

// Deploy deploys a fresh instance of the EtherAPIs contract, returning the
// transaction seeded into the network.
func (eapis *EtherAPIs) Deploy(from common.Address) (common.Address, *types.Transaction, error) {
	rpc, err := eapis.client.Stack().Attach()
	if err != nil {
		return common.Address{}, nil, err
	}
	auth := &bind.TransactOpts{
		From:   from,
		Signer: eapis.signer,
	}
	target, tx, etherapis, err := contract.DeployEtherAPIs(auth, backends.NewRPCBackend(rpc))
	if err == nil {
		eapis.contract = etherapis
	}
	return target, tx, err
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
	go eapis.eventmux.Post(NewAccountEvent{account.Address})

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
	go eapis.eventmux.Post(NewAccountEvent{key.Address})

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
	if err := eapis.ethereum.AccountManager().DeleteAccount(account, eapis.password); err != nil {
		return err
	}
	go eapis.eventmux.Post(DroppedAccountEvent{account})
	return nil
}

// RetrieveAccount returns the data associated with the account.
func (eapis *EtherAPIs) RetrieveAccount(account common.Address) Account {
	state, _ := eapis.ethereum.BlockChain().State()
	pendBlock, pendState := eapis.ethereum.Miner().Pending()

	txs := []*Transaction{}
	for _, tx := range pendBlock.Transactions() {
		from, _ := tx.From()
		if from == account || (tx.To() != nil && *tx.To() == account) {
			var to common.Address
			if tx.To() != nil {
				to = *tx.To()
			}
			txs = append(txs, &Transaction{
				Hash:   tx.Hash(),
				From:   from,
				To:     to,
				Amount: tx.Value(),
				Fees:   new(big.Int).Mul(tx.Gas(), tx.GasPrice()),
			})
		}
	}
	return Account{
		Nonce:        pendState.GetNonce(account),
		Balance:      state.GetBalance(account),
		Change:       new(big.Int).Sub(pendState.GetBalance(account), state.GetBalance(account)),
		Transactions: txs,
	}
}

// ListAccounts retrieves the list of accounts known to etherapis.
func (eapis *EtherAPIs) ListAccounts() ([]common.Address, error) {
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

// Transfer initiates a value transfer from an origin account to a destination
// account.
func (eapis *EtherAPIs) Transfer(from, to common.Address, amount *big.Int) (common.Hash, error) {
	// Make sure we actually own the origin account and have a valid destination
	accman := eapis.ethereum.AccountManager()
	if !accman.HasAccount(from) {
		return common.Hash{}, fmt.Errorf("unknown account: 0x%x", from.Bytes())
	}
	if to == (common.Address{}) {
		return common.Hash{}, fmt.Errorf("missing destination account")
	}
	// Serialize transaction creations to avoid nonce clashes
	eapis.txlock.Lock()
	defer eapis.txlock.Unlock()

	// Assemble and create the new transaction
	var (
		txpool   = eapis.ethereum.TxPool()
		nonce    = txpool.State().GetNonce(from)
		gasLimit = params.TxGas
		gasPrice = eapis.ethereum.GpoMinGasPrice
	)
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)

	// Sign the transaction and inject into the local pool for propagation
	signature, err := accman.Sign(accounts.Account{Address: from}, tx.SigHash().Bytes())
	if err != nil {
		return common.Hash{}, err
	}
	signed, err := tx.WithSignature(signature)
	if err != nil {
		return common.Hash{}, err
	}
	txpool.SetLocal(signed)
	if err := txpool.Add(signed); err != nil {
		return common.Hash{}, err
	}
	return signed.Hash(), nil
}

// CreateService registers a new service into the API marketplace.
func (eapis *EtherAPIs) CreateService(owner common.Address, name, url string, model *big.Int, price *big.Int, cancel uint64) (*types.Transaction, error) {
	auth := &bind.TransactOpts{
		From:   owner,
		Signer: eapis.signer,
	}
	return eapis.contract.AddService(auth, name, url, model, price, big.NewInt(int64(cancel)))
}

// DeleteService sets a service to deleted, permanently removing it from all listings.
func (eapis *EtherAPIs) DeleteService(owner common.Address, id *big.Int) (*types.Transaction, error) {
	auth := &bind.TransactOpts{
		From:   owner,
		Signer: eapis.signer,
	}
	return eapis.contract.DeleteService(auth, id)
}

// LockService disables a service for future subscriptions.
func (eapis *EtherAPIs) LockService(owner common.Address, id *big.Int) (*types.Transaction, error) {
	auth := &bind.TransactOpts{
		From:   owner,
		Signer: eapis.signer,
	}
	return eapis.contract.DisableService(auth, id)
}

// UnlockService enabled a service for future subscriptions.
func (eapis *EtherAPIs) UnlockService(owner common.Address, id *big.Int) (*types.Transaction, error) {
	auth := &bind.TransactOpts{
		From:   owner,
		Signer: eapis.signer,
	}
	return eapis.contract.EnableService(auth, id)
}

// Services retrieves a map of locally owned services, grouped by owner account.
func (eapis *EtherAPIs) Services() (map[common.Address][]*Service, error) {
	// Fetch all the accounts owned by this node
	addresses, err := eapis.ListAccounts()
	if err != nil {
		return nil, err
	}
	// For each address, retrieves all the registered services
	services := make(map[common.Address][]*Service)
	for _, address := range addresses {
		// Retrieve the number of services belonging to a user (currently creating included)
		count, err := eapis.contract.UserServicesLength(&bind.CallOpts{Pending: true}, address)
		if err != nil {
			return nil, fmt.Errorf("userServicesLength: %v", err)
		}
		// Make sure an empty list is always reported
		services[address] = []*Service{}

		// Retrieve each of the services individually
		for i := int64(0); i < count.Int64(); i++ {
			// Retrieve the users Nth service, both pending and stable
			id, err := eapis.contract.UserServices(&bind.CallOpts{Pending: true}, address, big.NewInt(i))
			if err != nil {
				return nil, fmt.Errorf("userServices: %v", err)
			}
			pending, err := eapis.contract.GetService(&bind.CallOpts{Pending: true}, id)
			if err != nil {
				return nil, fmt.Errorf("getService: %v", err)
			}
			stable, err := eapis.contract.GetService(&bind.CallOpts{Pending: false}, id)
			creating := err != nil

			if !creating && stable.Deleted {
				continue // Ignore long deleted services
			}

			// Convert to out internal type and accumulate
			services[address] = append(services[address], &Service{
				ID:           id,
				Name:         pending.Name,
				Owner:        pending.Owner,
				Endpoint:     pending.Endpoint,
				Model:        pending.Model,
				Price:        pending.Price,
				Cancellation: pending.Cancellation,
				Enabled:      pending.Enabled,

				Creating: creating,
				Changing: pending.Enabled != stable.Enabled,
				Deleting: pending.Deleted,
			})
		}
	}
	return services, nil
}

// Marketplace retrieves all the available services from the marketplace. It only
// ever operates on the currently final contract state (i.e. no pending data).
func (eapis *EtherAPIs) Marketplace() ([]*Service, error) {
	// Retrieve the total number of services in the marketplace
	count, err := eapis.contract.ServicesLength(nil)
	if err != nil {
		return nil, fmt.Errorf("servicesLength: %v", err)
	}
	// Retrieve each of the services individually
	services := make([]*Service, 0, count.Int64())
	for i := int64(0); i < count.Int64(); i++ {
		id := big.NewInt(i)

		service, err := eapis.contract.GetService(nil, id)
		if err != nil {
			return nil, fmt.Errorf("getService: %v", err)
		}
		if !service.Enabled || service.Deleted {
			continue
		}
		services = append(services, &Service{
			ID:           id,
			Name:         service.Name,
			Owner:        service.Owner,
			Endpoint:     service.Endpoint,
			Model:        service.Model,
			Price:        service.Price,
			Cancellation: service.Cancellation,
		})
	}
	return services, nil
}

// Geth retrieves the Ethereum client through which to interact with the underlying
// peer-to-peer networking layer.
func (eapis *EtherAPIs) Geth() *geth.Geth {
	return eapis.client
}

// Ethereum retrieves the Ethereum protocol running within the connected client.
func (eapis *EtherAPIs) Ethereum() *eth.Ethereum {
	return eapis.ethereum
}

// Contract retrieves the Ether APIs Ethereum contract to access the consensus data.
func (eapis *EtherAPIs) Contract() *contract.EtherAPIs {
	return eapis.contract
}

// CallRPC is a temporary helper method to pass an RPC call to the underlying
// go-ethereum server. It returns the exact raw response, no parsing done.
func (eapis *EtherAPIs) CallRPC(method string, params []interface{}) (json.RawMessage, error) {
	return eapis.rpcapi.Request(method, params)
}
