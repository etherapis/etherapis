package main

import (
	"crypto/ecdsa"
	"math/big"
	"strings"
	"sync"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/core/state"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/core/types"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/core/vm"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/crypto"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/eth/filters"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/ethdb"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/event"
)

// makeChannelName returns the canonical channel name based on the from and to
// paramaters.
func makeChannelName(from, to common.Address) []byte {
	return sha3(from[:], to[:])
}

type Contract struct {
	abi     abi.ABI
	stateFn func() *state.StateDB
	filters *filters.FilterSystem
	mux     *event.TypeMux
	db      ethdb.Database

	channels  map[common.Hash]*Channel
	channelMu sync.RWMutex
}

// contractAddress is the static address on which the contract resides
var contractAddress = common.HexToAddress("0xaa")

// GetContract initialises a new abi and returns the contract. It does not
// deploy the contract, hence the name.
func GetContract(db ethdb.Database, mux *event.TypeMux, stateFn func() *state.StateDB) (*Contract, error) {
	contract := Contract{
		stateFn: stateFn,
		filters: filters.NewFilterSystem(mux),
	}

	var err error
	contract.abi, err = abi.JSON(strings.NewReader(jsonAbi))
	if err != nil {
		return nil, err
	}

	return &contract, nil
}

func (c *Contract) Stop() {
	c.filters.Stop()
}

func (c *Contract) NewChannel(key *ecdsa.PrivateKey, to common.Address, amount *big.Int, cb func(*Channel)) (*types.Transaction, error) {
	from := crypto.PubkeyToAddress(key.PublicKey)

	data, err := c.abi.Pack("createChannel", to)
	if err != nil {
		return nil, err
	}

	statedb := c.stateFn()
	transaction, err := types.NewTransaction(statedb.GetNonce(from), to, amount, big.NewInt(250000), big.NewInt(50000000000), data).SignECDSA(key)
	if err != nil {
		return nil, err
	}

	evId := c.abi.Events["NewChannel"].Id
	filter := filters.New(c.db)
	filter.SetAddresses([]common.Address{contractAddress})
	filter.SetTopics([][]common.Hash{ // TODO refactor, helper
		[]common.Hash{evId},
		[]common.Hash{from.Hash()},
		[]common.Hash{to.Hash()},
	})
	filter.SetBeginBlock(0)
	filter.SetEndBlock(-1)
	filter.LogsCallback = func(logs vm.Logs) {
		// tere should really be only one log. TODO this part
		log := logs[0]
		// [ event_id, from, to, [channel_id, nonce ] ]
		// TODO: do to and from validation here
		/*
			from := log.Topics[1]
			to := log.Topics[2]
		*/
		channelId := common.BytesToHash(log.Data[0:31])
		//nonce := log.Data[31:]

		c.channelMu.Lock()
		defer c.channelMu.Unlock()

		channel, exist := c.channels[channelId]
		if !exist {
			channel = NewChannel(c, channelId, from, to)
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

	chContract *Contract
}

func NewChannel(c *Contract, id common.Hash, from, to common.Address) *Channel {
	return &Channel{
		Id:         id,
		from:       from,
		to:         to,
		chContract: c,
	}
}

const jsonAbi = `[{"constant":false,"inputs":[],"name":"Channel","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"channel","type":"bytes32"}],"name":"isValidChannel","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"channel","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"claim","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"}],"name":"createChannel","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"channel","type":"bytes32"},{"name":"recipient","type":"address"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"}],"name":"getHash","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"channels","outputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"validUntil","type":"uint256"},{"name":"valid","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"channel","type":"bytes32"}],"name":"getChannelValidUntil","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"channel","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"verify","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"channel","type":"bytes32"}],"name":"reclaim","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"channel","type":"bytes32"}],"name":"getChannelOwner","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[{"name":"channel","type":"bytes32"}],"name":"deposit","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"channel","type":"bytes32"}],"name":"getChannelValue","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"channel","type":"bytes32"}],"name":"NewChannel","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"channel","type":"bytes32"}],"name":"Deposit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"who","type":"address"},{"indexed":true,"name":"channel","type":"bytes32"}],"name":"Claim","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"channel","type":"bytes32"}],"name":"Reclaim","type":"event"}]`
