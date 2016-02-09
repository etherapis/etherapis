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
var contractAddress = common.HexToAddress("0x187be5cddd95d1e3c4ce06b659c7c4ce8c696574")

// signFn is a signer function callback when the contract requires a method to
// sign the transaction before submission.
type signFn func(*types.Transaction) (*types.Transaction, error)

// Channels is the channels contract reflecting that on the ethereum network. The
// channels contract handles all validation and verifications of payments and
// allows you to redeem cheques.
//
// Channels implements the proxy.Verifier and proxy.Charges interfaces.
type Channels struct {
	abi        abi.ABI
	blockchain *core.BlockChain

	filters *filters.FilterSystem
	mux     *event.TypeMux
	db      ethdb.Database

	channels  map[common.Hash]*Channel
	channelMu sync.RWMutex

	// call key is a temporary key used to do calls
	callKey *ecdsa.PrivateKey

	callState func() *state.StateDB
}

// Fetch initialises a new abi and returns the contract. It does not
// deploy the contract, hence the name.
func Fetch(db ethdb.Database, mux *event.TypeMux, blockchain *core.BlockChain, callState func() *state.StateDB) (*Channels, error) {
	contract := Channels{
		blockchain: blockchain,
		channels:   make(map[common.Hash]*Channel),
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

func (c *Channels) Stop() {
	c.filters.Stop()
}

// Exists returns whether there exists a channel between transactor and beneficiary.
func (c *Channels) Exists(from, to common.Address) bool {
	return c.Call("isValidChannel", c.ChannelId(from, to)).(bool)
}

// Validate validates the ECDSA (curve=secp256k1) signature with the given input
// where H=KECCAK(from, to, amount) and the validation must satisfy:
// channel_owner == ECRECOVER(H, S) where S is the given signature signed by
// the sender.
func (c *Channels) ValidateSig(from, to common.Address, nonce uint64, amount *big.Int, sig []byte) bool {
	if len(sig) != 65 {
		// invalid signature
		return false
	}

	channelId := c.ChannelId(from, to)
	signature := bytesToSignature(sig)
	return c.Call("verifySignature", channelId, nonce, amount, signature.v, signature.r, signature.s).(bool)
}

func (c *Channels) Verify(from, to common.Address, nonce uint64, amount *big.Int, sig []byte) (bool, bool) {
	if len(sig) != 65 {
		// invalid signature
		return false, false
	}

	channelId := c.ChannelId(from, to)
	signature := bytesToSignature(sig)
	validPayment := c.Call("verifyPayment", channelId, nonce, amount, signature.v, signature.r, signature.s).(bool)
	enoughFunds := c.Call("getChannelValue", c.ChannelId(from, to)).(*big.Int).Cmp(amount) >= 0
	return validPayment, enoughFunds
}

func (c *Channels) Price(from, to common.Address) *big.Int {
	return c.Call("getChannelPrice", c.ChannelId(from, to)).(*big.Int)
}

func (c *Channels) Nonce(from, to common.Address) *big.Int {
	return c.Call("getChannelNonce", c.ChannelId(from, to)).(*big.Int)
}

// Claim redeems a given signature using the canonical channel. It creates an
// Ethereum transaction and submits it to the Ethereum network.
//
// Chaim returns the unsigned transaction and an error if it failed.
func (c *Channels) Claim(signer common.Address, from, to common.Address, nonce uint64, amount *big.Int, sig []byte) (*types.Transaction, error) {
	if len(sig) != 65 {
		return nil, fmt.Errorf("Invalid signature. Signature requires to be 65 bytes")
	}

	channelId := c.ChannelId(from, to)
	signature := bytesToSignature(sig)

	txData, err := c.abi.Pack("claim", channelId, nonce, amount, signature.v, signature.r, signature.s)
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
func (c *Channels) Call(methodName string, v ...interface{}) interface{} {
	return c.abi.Call(c.exec, methodName, v...)
}

// ChannelId returns the canonical channel name for transactor and beneficiary
func (c *Channels) ChannelId(from, to common.Address) common.Hash {
	return common.BytesToHash(c.Call("makeChannelId", from, to).([]byte))
}

// exec is the executer function callback for the abi `Call` method.
func (c *Channels) exec(input []byte) []byte {
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
func (c *Channels) NewChannel(key *ecdsa.PrivateKey, to common.Address, amount, price *big.Int, cb func(*Channel)) (*types.Transaction, error) {
	from := crypto.PubkeyToAddress(key.PublicKey)

	data, err := c.abi.Pack("createChannel", to, price)
	if err != nil {
		return nil, err
	}

	statedb, err := c.blockchain.State()
	if err != nil {
		return nil, err
	}

	transaction, err := types.NewTransaction(statedb.GetNonce(from), contractAddress, amount, big.NewInt(250000), big.NewInt(50000000000), data).SignECDSA(key)
	if err != nil {
		return nil, err
	}

	evId := c.abi.Events["NewChannel"].Id()
	filter := filters.New(c.db)
	filter.SetAddresses([]common.Address{contractAddress})
	filter.SetTopics([][]common.Hash{ // TODO refactor, helper
		[]common.Hash{evId},
		[]common.Hash{from.Hash()},
		[]common.Hash{to.Hash()},
	})
	filter.SetBeginBlock(0)
	filter.SetEndBlock(-1)
	filter.LogCallback = func(log *vm.Log, removed bool) {
		// TODO: do to and from validation here
		/*
			from := log.Topics[1]
			to := log.Topics[2]
		*/
		channelId := common.BytesToHash(log.Data[0:31])
		nonce := common.BytesToBig(log.Data[31:])

		c.channelMu.Lock()
		defer c.channelMu.Unlock()

		channel, exist := c.channels[channelId]
		if !exist {
			channel = NewChannel(c, channelId, from, to, nonce)
			c.channels[channelId] = channel
		}
		cb(channel)
	}

	c.filters.Add(filter)

	return transaction, nil
}

type Channel struct {
	Id       common.Hash
	key      *ecdsa.PrivateKey
	from, to common.Address
	nonce    *big.Int

	channels *Channels
}

// NewChannel returns a new payment channel.
func NewChannel(c *Channels, id common.Hash, from, to common.Address, nonce *big.Int) *Channel {
	return &Channel{
		Id:       id,
		from:     from,
		to:       to,
		channels: c,
	}
}

type Cheque struct {
	Sig           []byte
	From, To      common.Address
	Nonce, Amount *big.Int
}

// SignPayment returns a signed transaction on the current payment channel.
func (c *Channel) SignPayment(amount *big.Int) (Cheque, error) {
	sig, err := crypto.Sign(sha3(c.Id[:], c.from[:], c.to[:], c.nonce.Bytes(), amount.Bytes()), c.key)
	if err != nil {
		return Cheque{}, err
	}
	return Cheque{Sig: sig, From: c.from, To: c.to, Nonce: c.nonce, Amount: amount}, nil
}

const jsonAbi = `[{"constant":false,"inputs":[{"name":"serviceId","type":"uint256"}],"name":"subscribe","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionIdNonce","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"name","type":"string"},{"name":"endpoint","type":"string"},{"name":"price","type":"uint256"},{"name":"cancellationTime","type":"uint256"}],"name":"addService","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"serviceLength","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"claim","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"verifyPayment","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionIdClosedAt","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"from","type":"address"},{"name":"serviceId","type":"uint256"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"}],"name":"getHash","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"reclaim","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"from","type":"address"},{"name":"serviceId","type":"uint256"}],"name":"makeSubscriptionId","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"deposit","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionIdOwner","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"cancel","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"isValidSubscription","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"serviceId","type":"uint256"}],"name":"getService","outputs":[{"name":"","type":"string"},{"name":"","type":"string"},{"name":"","type":"uint256"},{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionIdValue","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"verifySignature","outputs":[{"name":"","type":"bool"}],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"name","type":"string"},{"indexed":true,"name":"owner","type":"address"},{"indexed":false,"name":"serviceId","type":"uint256"}],"name":"NewService","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"serviceId","type":"uint256"},{"indexed":false,"name":"subscriptionId","type":"bytes32"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"NewSubscription","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"subscriptionId","type":"bytes32"}],"name":"Deposit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionId","type":"bytes32"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"Redeem","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionId","type":"bytes32"},{"indexed":false,"name":"closedAt","type":"uint256"}],"name":"Cancel","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionId","type":"bytes32"}],"name":"Reclaim","type":"event"}]`
