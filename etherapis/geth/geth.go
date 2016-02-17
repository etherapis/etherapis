// Package geth implements the wrapper around to go-ethereum client
package geth

import (
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/node"
)

// Geth is a wrapper around the Ethereum Go client.
type Geth struct {
	stack    *node.Node      // Ethereum network node / protocol stack
	keystore crypto.KeyStore // Keystore to retrieve private keys from
}

// New creates a Ethereum client, pre-configured to one of the supported networks.
func New(datadir string, network EthereumNetwork) (*Geth, error) {
	// Tag the data dir with the network name
	switch network {
	case MainNet:
		datadir = filepath.Join(datadir, "mainnet")
	case TestNet:
		datadir = filepath.Join(datadir, "testnet")
	default:
		return nil, fmt.Errorf("unsupported network: %v", network)
	}
	// Select the bootstrap nodes based on the network
	bootnodes := utils.FrontierBootNodes
	if network == TestNet {
		bootnodes = utils.TestNetBootNodes
	}
	// Configure the node's service container
	stackConf := &node.Config{
		DataDir:        datadir,
		Name:           common.MakeName(NodeName, NodeVersion),
		BootstrapNodes: bootnodes,
		ListenAddr:     fmt.Sprintf(":%d", NodePort),
		MaxPeers:       NodeMaxPeers,
	}
	// Configure the bare-bone Ethereum service
	keystore := crypto.NewKeyStorePassphrase(filepath.Join(datadir, "keystore"), crypto.StandardScryptN, crypto.StandardScryptP)
	ethConf := &eth.Config{
		FastSync:       true,
		DatabaseCache:  64,
		NetworkId:      int(network),
		AccountManager: accounts.NewManager(keystore),

		// Blatantly initialize the gas oracle to the defaults from go-ethereum
		GpoMinGasPrice:          new(big.Int).Mul(big.NewInt(50), common.Shannon),
		GpoMaxGasPrice:          new(big.Int).Mul(big.NewInt(500), common.Shannon),
		GpoFullBlockRatio:       80,
		GpobaseStepDown:         10,
		GpobaseStepUp:           100,
		GpobaseCorrectionFactor: 110,
	}
	// Override any default configs in the test network
	if network == TestNet {
		ethConf.NetworkId = 2
		ethConf.Genesis = core.TestNetGenesisBlock()
		state.StartingNonce = 1048576 // (2**20)
	}
	// Assemble and return the protocol stack
	stack, err := node.New(stackConf)
	if err != nil {
		return nil, fmt.Errorf("protocol stack: %v", err)
	}
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) { return eth.New(ctx, ethConf) }); err != nil {
		return nil, fmt.Errorf("ethereum service: %v", err)
	}
	return &Geth{
		stack:    stack,
		keystore: keystore,
	}, nil
}

// Start boots up the Ethereum protocol, starts interacting with the P2P network
// and opens up the IPC based JSON RPC API for accessing the exposed APIs.
func (g *Geth) Start() error {
	return g.stack.Start()
}

// Stop closes down the Ethereum client, along with any other resources it might
// be keeping around.
func (g *Geth) Stop() error {
	return g.stack.Stop()
}

// Stack is a quick hack to expose the internal Ethereum node implementation. We
// should probably remove this after the API interface is implemented, but until
// then it makes things simpler.
func (g *Geth) Stack() *node.Node {
	return g.stack
}

// Keystore is a quick hack to expose the internal Ethereum keystore. We should
//  remove this after the API interface is implemented, but until then it makes
// things simpler.
func (g *Geth) Keystore() crypto.KeyStore {
	return g.keystore
}

// Attach connects to the running node's IPC exposed APIs, and returns a Go API
// interface.
func (g *Geth) Attach() (*API, error) {
	client, err := g.stack.Attach()
	if err != nil {
		return nil, err
	}
	return &API{client: client}, nil
}
