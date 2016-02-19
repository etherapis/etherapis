package contract

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"gopkg.in/inconshreveable/log15.v2"
)

// contractAddress is the static address on which the contract resides
var contractAddress = common.HexToAddress("0x406e9fff90231f97b2f0d2832001b49d57df4dd2")

// signFn is a signer function callback when the contract requires a method to
// sign the transaction before submission.
type signFn func(*types.Transaction) (*types.Transaction, error)

// Contract is the wrapper object that reflects the contract on the Ethereum
// network.
//
// Contract implements the proxy.Verifier and proxy.Charges interfaces.
type Contract struct {
	abi        abi.ABI
	blockchain *core.BlockChain

	filters *filters.FilterSystem
	mux     *event.TypeMux
	db      ethdb.Database

	subs      map[common.Hash]*Subscription
	channelMu sync.RWMutex

	// call key is a temporary key used to do calls
	callKey   *ecdsa.PrivateKey
	callMu    sync.Mutex
	callState func() *state.StateDB
}

// New initialises a new abi and returns the contract. It does not
// deploy the contract, hence the name.
func New(db ethdb.Database, mux *event.TypeMux, blockchain *core.BlockChain, callState func() *state.StateDB) (*Contract, error) {
	contract := Contract{
		blockchain: blockchain,
		subs:       make(map[common.Hash]*Subscription),
		filters:    filters.NewFilterSystem(mux),
		callState:  callState,
	}
	contract.callKey, _ = crypto.GenerateKey()

	var err error
	contract.abi, err = abi.JSON(strings.NewReader(jsonAbi))
	if err != nil {
		return nil, err
	}

	return &contract, nil
}

// Stop stops the contract from watching any events that might be present.
func (c *Contract) Stop() {
	c.filters.Stop()
}

// Exists returns whether there exists a channel between transactor and beneficiary.
func (c *Contract) Exists(from common.Address, serviceId *big.Int) (exists bool) {
	if err := c.Call(&exists, "isValidSubscription", c.SubscriptionId(from, serviceId)); err != nil {
		log15.Warn("exists", "error", err)
	}

	return exists
}

// Validate validates the ECDSA (curve=secp256k1) signature with the given input
// where H=KECCAK(from, to, amount) and the validation must satisfy:
// channel_owner == ECRECOVER(H, S) where S is the given signature signed by
// the sender.
func (c *Contract) ValidateSig(from common.Address, serviceId *big.Int, nonce uint64, amount *big.Int, sig []byte) (validSig bool) {
	if len(sig) != 65 {
		// invalid signature
		return false
	}

	subscriptionId := c.SubscriptionId(from, serviceId)
	signature := bytesToSignature(sig)

	if err := c.Call(&validSig, "verifySignature", subscriptionId, nonce, amount, signature.v, signature.r, signature.s); err != nil {
		log15.Warn("verifySignature", "error", err)
	}

	return validSig
}

// Service represents a service created by the contract.
//
// Service is used to unmarshal contract return values in.
type Service struct {
	Name             string         `json:"name"`
	Endpoint         string         `json:"endpoint"`
	Price            *big.Int       `json:"price"`
	CancellationTime *big.Int       `json:"cancellationTime"`
	Owner            common.Address `json:"owner"`
	Enabled          bool           `json:"enabled"`
}

// Verify verifies a payment with the given input parameters. It validates the
// the signature using off-chain verification methods using the EVM and the
// managed contract.
func (c *Contract) Verify(from common.Address, serviceId *big.Int, nonce uint64, amount *big.Int, sig []byte) (bool, bool) {
	if len(sig) != 65 {
		// invalid signature
		return false, false
	}

	subscriptionId := c.SubscriptionId(from, serviceId)
	signature := bytesToSignature(sig)

	var validPayment bool
	if err := c.Call(&validPayment, "verifyPayment", subscriptionId, nonce, amount, signature.v, signature.r, signature.s); err != nil {
		log15.Warn("verifyPayment", "error", err)
	}

	var subscription Subscription
	if err := c.Call(&subscription, "getSubscription", subscriptionId); err != nil {
		log15.Warn("getService", "error", err)
	}

	return validPayment, subscription.Value.Cmp(amount) >= 0
}

// Price returns the price associated with the given serviceId.
func (c *Contract) Price(from common.Address, serviceId *big.Int) *big.Int {
	var service Service
	if err := c.Call(&service, "getService", serviceId); err != nil {
		log15.Warn("getService", "error", err)
	}
	return service.Price
}

// Nonce returns the nonce associated to a subscription.
func (c *Contract) Nonce(from common.Address, serviceId *big.Int) (nonce *big.Int) {
	if err := c.Call(&nonce, "getSubscriptionNonce", c.SubscriptionId(from, serviceId)); err != nil {
		log15.Warn("getSubscriptionNonce", "error", err)
	}
	return nonce
}

// Claim redeems a given signature using the canonical channel. It creates an
// Ethereum transaction and submits it to the Ethereum network.
//
// Chaim returns the unsigned transaction and an error if it failed.
func (c *Contract) Claim(signer common.Address, from common.Address, serviceId *big.Int, nonce uint64, amount *big.Int, sig []byte) (*types.Transaction, error) {
	if len(sig) != 65 {
		return nil, fmt.Errorf("Invalid signature. Signature requires to be 65 bytes")
	}

	subscriptionId := c.SubscriptionId(from, serviceId)
	signature := bytesToSignature(sig)

	txData, err := c.abi.Pack("claim", subscriptionId, nonce, amount, signature.v, signature.r, signature.s)
	if err != nil {
		return nil, err
	}

	statedb, _ := c.blockchain.State()
	gasPrice := big.NewInt(50000000000)
	gasLimit := big.NewInt(250000)
	tx := types.NewTransaction(statedb.GetNonce(signer), contractAddress, new(big.Int), gasLimit, gasPrice, txData)
	return tx, nil
}

// Call calls the contract method with v as input and sets the output to r.
//
// Call returns an error if the ABI failed parsing the input or if no method is
// present.
func (c *Contract) Call(r interface{}, method string, v ...interface{}) error {
	return c.abi.Call(c.exec, r, method, v...)
}

// SubscriptionId returns the canonical channel name for transactor and beneficiary
func (c *Contract) SubscriptionId(from common.Address, serviceId *big.Int) common.Hash {
	var id []byte
	if err := c.Call(&id, "makeSubscriptionId", from, serviceId); err != nil {
		log15.Warn("makeSubscriptionId", "error", err)
	}
	return common.BytesToHash(id)
}

// exec is the executer function callback for the abi `Call` method.
func (c *Contract) exec(input []byte) []byte {
	c.callMu.Lock()
	defer c.callMu.Unlock()

	ret, err := runtime.Call(contractAddress, input, &runtime.Config{
		GetHashFn: core.GetHashFn(c.blockchain.CurrentBlock().ParentHash(), c.blockchain),
		State:     c.callState(),
	})
	if err != nil {
		log15.Warn("execution failed", "error", err)
		return nil
	}

	return ret
}

// Start Go API. Not important for this version
func (c *Contract) Subscribe(key *ecdsa.PrivateKey, serviceId *big.Int, amount, price *big.Int, cb func(*Subscription)) (*types.Transaction, error) {
	from := crypto.PubkeyToAddress(key.PublicKey)

	data, err := c.abi.Pack("subscribe", serviceId)
	if err != nil {
		return nil, err
	}

	statedb, err := c.blockchain.State()
	if err != nil {
		return nil, err
	}

	transaction, err := types.NewTransaction(statedb.GetNonce(from), contractAddress, amount, big.NewInt(600000), big.NewInt(50000000000), data).SignECDSA(key)
	if err != nil {
		return nil, err
	}

	evId := c.abi.Events["NewSubscription"].Id()
	filter := filters.New(c.db)
	filter.SetAddresses([]common.Address{contractAddress})
	filter.SetTopics([][]common.Hash{ // TODO refactor, helper
		[]common.Hash{evId},
		[]common.Hash{from.Hash()},
		[]common.Hash{common.BigToHash(serviceId)},
	})
	filter.SetBeginBlock(0)
	filter.SetEndBlock(-1)
	filter.LogCallback = func(log *vm.Log, removed bool) {
		// TODO: do to and from validation here
		/*
			from := log.Topics[1]
			to := log.Topics[2]
		*/
		subscriptionId := common.BytesToHash(log.Data[0:31])
		nonce := common.BytesToBig(log.Data[31:])

		c.channelMu.Lock()
		defer c.channelMu.Unlock()

		channel, exist := c.subs[subscriptionId]
		if !exist {
			channel = NewSubscription(c, subscriptionId, from, serviceId, nonce)
			c.subs[subscriptionId] = channel
		}
		cb(channel)
	}

	c.filters.Add(filter, filters.PendingLogFilter)

	return transaction, nil
}

// Services returns the services associated with the given account in addr.
func (s *Contract) Services(addr common.Address) ([]Service, error) {
	var len *big.Int
	if err := s.Call(&len, "userServicesLength", addr); err != nil {
		return nil, err
	}

	services := make([]Service, len.Uint64())
	for i := 0; i < int(len.Uint64()); i++ {
		var idx *big.Int
		if err := s.Call(&idx, "userServices", addr, i); err != nil {
			return nil, err
		}

		if err := s.Call(&services[i], "getService", idx); err != nil {
			return nil, err
		}
	}

	return services, nil
}

// AllServices returns the services currently in existence.
func (s *Contract) AllServices() ([]Service, error) {
	var len *big.Int
	if err := s.Call(&len, "serviceLength"); err != nil {
		return nil, err
	}

	services := make([]Service, len.Uint64())
	for i := 0; i < int(len.Uint64()); i++ {
		if err := s.Call(&services[i], "getService", i); err != nil {
			return nil, err
		}
	}

	return services, nil
}

// Subscriptions returns all subscriptions currently in existence.
func (s *Contract) Subscriptions(addr common.Address) ([]Subscription, error) {
	var len *big.Int
	if err := s.Call(&len, "userSubscriptionsLength", addr); err != nil {
		return nil, err
	}

	subscriptions := make([]Subscription, len.Uint64())
	for i := 0; i < int(len.Uint64()); i++ {
		var hash common.Hash
		if err := s.Call(&hash, "userSubscriptions", addr, i); err != nil {
			return nil, err
		}

		if err := s.Call(&subscriptions[i], "getSubscription", hash); err != nil {
			return nil, err
		}
		subscriptions[i].Id = hash
	}

	return subscriptions, nil
}

// Subscription represents a user subscription on a service.
//
// Subscription is used to unmarshal contract return values in.
type Subscription struct {
	key *ecdsa.PrivateKey

	Id        common.Hash    `json:"id"`
	From      common.Address `json:"from"`
	ServiceId *big.Int       `json:"serviceId"`
	Nonce     *big.Int       `json:"nonce"`
	Value     *big.Int       `json:"value"`
	Cancelled bool           `json:"cancelled"`
	ClosedAt  *big.Int       `json:"closedAt"`

	contract *Contract
}

// NewSubscription returns a new payment channel.
func NewSubscription(contract *Contract, id common.Hash, from common.Address, serviceId *big.Int, nonce *big.Int) *Subscription {
	return &Subscription{
		Id:        id,
		From:      from,
		ServiceId: serviceId,
		contract:  contract,
	}
}

// SignPayment returns a signed transaction on the current payment channel.
func (c *Subscription) SignPayment(amount *big.Int) (Cheque, error) {
	sig, err := crypto.Sign(sha3(c.Id[:], c.From[:], c.ServiceId.Bytes(), c.Nonce.Bytes(), amount.Bytes()), c.key)
	if err != nil {
		return Cheque{}, err
	}
	return Cheque{Sig: sig, From: c.From, ServiceId: c.ServiceId, Nonce: c.Nonce, Amount: amount}, nil
}

// Checque is a signed payment and can be verified off-chain.
type Cheque struct {
	Sig           []byte
	From          common.Address
	ServiceId     *big.Int
	Nonce, Amount *big.Int
}

// contract ABI
const jsonAbi = `[{"constant":false,"inputs":[{"name":"serviceId","type":"uint256"}],"name":"subscribe","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"addr","type":"address"}],"name":"userServicesLength","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscription","outputs":[{"name":"from","type":"address"},{"name":"serviceId","type":"uint256"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"cancelled","type":"bool"},{"name":"closedAt","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"name","type":"string"},{"name":"endpoint","type":"string"},{"name":"price","type":"uint256"},{"name":"cancellationTime","type":"uint256"}],"name":"addService","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"serviceLength","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"claim","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"uint256"}],"name":"userServices","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"verifyPayment","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionClosedAt","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionNonce","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"from","type":"address"},{"name":"serviceId","type":"uint256"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"}],"name":"getHash","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionOwner","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"reclaim","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionValue","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"from","type":"address"},{"name":"serviceId","type":"uint256"}],"name":"makeSubscriptionId","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"deposit","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"cancel","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"uint256"}],"name":"userSubscriptions","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionServiceId","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"addr","type":"address"}],"name":"userSubscriptionsLength","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"isValidSubscription","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"serviceId","type":"uint256"}],"name":"getService","outputs":[{"name":"name","type":"string"},{"name":"owner","type":"address"},{"name":"endpoint","type":"string"},{"name":"price","type":"uint256"},{"name":"cancellationTime","type":"uint256"},{"name":"enabled","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"verifySignature","outputs":[{"name":"","type":"bool"}],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"name","type":"string"},{"indexed":true,"name":"owner","type":"address"},{"indexed":false,"name":"serviceId","type":"uint256"}],"name":"NewService","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"serviceId","type":"uint256"},{"indexed":false,"name":"subscriptionId","type":"bytes32"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"NewSubscription","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"subscriptionId","type":"bytes32"}],"name":"Deposit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionId","type":"bytes32"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"Redeem","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionId","type":"bytes32"},{"indexed":false,"name":"closedAt","type":"uint256"}],"name":"Cancel","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionId","type":"bytes32"}],"name":"Reclaim","type":"event"}]`
