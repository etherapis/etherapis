package channels

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
var contractAddress = common.HexToAddress("0x1e0757a6d4a211514028f98d6deb25e10a41247f")

// signFn is a signer function callback when the contract requires a method to
// sign the transaction before submission.
type signFn func(*types.Transaction) (*types.Transaction, error)

// Subscriptions is the channels contract reflecting that on the ethereum network. The
// channels contract handles all validation and verifications of payments and
// allows you to redeem cheques.
//
// Subscriptions implements the proxy.Verifier and proxy.Charges interfaces.
type Subscriptions struct {
	abi        abi.ABI
	blockchain *core.BlockChain

	filters *filters.FilterSystem
	mux     *event.TypeMux
	db      ethdb.Database

	channels  map[common.Hash]*Subscription
	channelMu sync.RWMutex

	// call key is a temporary key used to do calls
	callKey *ecdsa.PrivateKey

	callState func() *state.StateDB
}

// Fetch initialises a new abi and returns the contract. It does not
// deploy the contract, hence the name.
func Fetch(db ethdb.Database, mux *event.TypeMux, blockchain *core.BlockChain, callState func() *state.StateDB) (*Subscriptions, error) {
	contract := Subscriptions{
		blockchain: blockchain,
		channels:   make(map[common.Hash]*Subscription),
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

func (c *Subscriptions) Stop() {
	c.filters.Stop()
}

// Exists returns whether there exists a channel between transactor and beneficiary.
func (c *Subscriptions) Exists(from common.Address, serviceId *big.Int) (exists bool) {
	if err := c.Call(&exists, "isValidSubscription", c.SubscriptionId(from, serviceId)); err != nil {
		log15.Warn("exists", "error", err)
	}

	return exists
}

// Validate validates the ECDSA (curve=secp256k1) signature with the given input
// where H=KECCAK(from, to, amount) and the validation must satisfy:
// channel_owner == ECRECOVER(H, S) where S is the given signature signed by
// the sender.
func (c *Subscriptions) ValidateSig(from common.Address, serviceId *big.Int, nonce uint64, amount *big.Int, sig []byte) (validSig bool) {
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

type Service struct {
	Name, Endpoint          string
	Price, CancellationTime *big.Int
}

type Sub struct {
	From         common.Address
	ServiceId    *big.Int
	Nonce, Value *big.Int
	Cancelled    bool
	ClosedAt     *big.Int
}

func (c *Subscriptions) Verify(from common.Address, serviceId *big.Int, nonce uint64, amount *big.Int, sig []byte) (bool, bool) {
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

	var subscription Sub
	if err := c.Call(&subscription, "getSubscription", subscriptionId); err != nil {
		log15.Warn("getService", "error", err)
	}

	return validPayment, subscription.Value.Cmp(amount) >= 0
}

func (c *Subscriptions) Price(from common.Address, serviceId *big.Int) *big.Int {
	var service Service
	if err := c.Call(&service, "getService", serviceId); err != nil {
		log15.Warn("getService", "error", err)
	}
	return service.Price
}

func (c *Subscriptions) Nonce(from common.Address, serviceId *big.Int) (nonce *big.Int) {
	if err := c.Call(&nonce, "getSubscriptionNonce", c.SubscriptionId(from, serviceId)); err != nil {
		log15.Warn("getSubscriptionNoce", "error", err)
	}
	return nonce
}

// Claim redeems a given signature using the canonical channel. It creates an
// Ethereum transaction and submits it to the Ethereum network.
//
// Chaim returns the unsigned transaction and an error if it failed.
func (c *Subscriptions) Claim(signer common.Address, from common.Address, serviceId *big.Int, nonce uint64, amount *big.Int, sig []byte) (*types.Transaction, error) {
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

// helper forwarder
func (c *Subscriptions) Call(r interface{}, methodName string, v ...interface{}) error {
	return c.abi.Call(c.exec, r, methodName, v...)
}

// SubscriptionId returns the canonical channel name for transactor and beneficiary
func (c *Subscriptions) SubscriptionId(from common.Address, serviceId *big.Int) common.Hash {
	var id []byte
	if err := c.Call(&id, "makeSubscriptionId", from, serviceId); err != nil {
		log15.Warn("makeSubscriptionId", "error", err)
	}
	return common.BytesToHash(id)
}

// exec is the executer function callback for the abi `Call` method.
func (c *Subscriptions) exec(input []byte) []byte {
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
func (c *Subscriptions) Subscribe(key *ecdsa.PrivateKey, serviceId *big.Int, amount, price *big.Int, cb func(*Subscription)) (*types.Transaction, error) {
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

		channel, exist := c.channels[subscriptionId]
		if !exist {
			channel = NewSubscription(c, subscriptionId, from, serviceId, nonce)
			c.channels[subscriptionId] = channel
		}
		cb(channel)
	}

	c.filters.Add(filter)

	return transaction, nil
}

type Subscription struct {
	Id        common.Hash
	key       *ecdsa.PrivateKey
	from      common.Address
	serviceId *big.Int
	nonce     *big.Int

	channels *Subscriptions
}

// NewSubscription returns a new payment channel.
func NewSubscription(c *Subscriptions, id common.Hash, from common.Address, serviceId *big.Int, nonce *big.Int) *Subscription {
	return &Subscription{
		Id:        id,
		from:      from,
		serviceId: serviceId,
		channels:  c,
	}
}

type Cheque struct {
	Sig           []byte
	From          common.Address
	ServiceId     *big.Int
	Nonce, Amount *big.Int
}

// SignPayment returns a signed transaction on the current payment channel.
func (c *Subscription) SignPayment(amount *big.Int) (Cheque, error) {
	sig, err := crypto.Sign(sha3(c.Id[:], c.from[:], c.serviceId.Bytes(), c.nonce.Bytes(), amount.Bytes()), c.key)
	if err != nil {
		return Cheque{}, err
	}
	return Cheque{Sig: sig, From: c.from, ServiceId: c.serviceId, Nonce: c.nonce, Amount: amount}, nil
}

const jsonAbi = `[{"constant":false,"inputs":[{"name":"serviceId","type":"uint256"}],"name":"subscribe","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscription","outputs":[{"name":"from","type":"address"},{"name":"serviceId","type":"uint256"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"cancelled","type":"bool"},{"name":"closedAt","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"name","type":"string"},{"name":"endpoint","type":"string"},{"name":"price","type":"uint256"},{"name":"cancellationTime","type":"uint256"}],"name":"addService","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"serviceLength","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"claim","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"verifyPayment","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionClosedAt","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionNonce","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"from","type":"address"},{"name":"serviceId","type":"uint256"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"}],"name":"getHash","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionOwner","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"reclaim","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionValue","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"from","type":"address"},{"name":"serviceId","type":"uint256"}],"name":"makeSubscriptionId","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"deposit","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"cancel","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionServiceId","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"isValidSubscription","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"serviceId","type":"uint256"}],"name":"getService","outputs":[{"name":"name","type":"string"},{"name":"endpoint","type":"string"},{"name":"price","type":"uint256"},{"name":"cancellationTime","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"verifySignature","outputs":[{"name":"","type":"bool"}],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"name","type":"string"},{"indexed":true,"name":"owner","type":"address"},{"indexed":false,"name":"serviceId","type":"uint256"}],"name":"NewService","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"serviceId","type":"uint256"},{"indexed":false,"name":"subscriptionId","type":"bytes32"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"NewSubscription","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"subscriptionId","type":"bytes32"}],"name":"Deposit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionId","type":"bytes32"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"Redeem","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionId","type":"bytes32"},{"indexed":false,"name":"closedAt","type":"uint256"}],"name":"Cancel","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionId","type":"bytes32"}],"name":"Reclaim","type":"event"}]`
