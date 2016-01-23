// Package geth implements the wrapper around to go-ethereum client
package geth

import (
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/accounts"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/core"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/core/state"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/crypto"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/eth"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/node"
)

// New creates a Ethereum client, pre-configured to one of the supported networks.
func New(datadir string, net EthereumNetwork) (*node.Node, error) {
	// Tag the data dir with the network name
	switch net {
	case MainNet:
		datadir = filepath.Join(datadir, "mainnet")
	case TestNet:
		datadir = filepath.Join(datadir, "testnet")
	default:
		return nil, fmt.Errorf("unsupported network: %v", net)
	}
	// Select the bootstrap nodes based on the network
	bootnodes := utils.FrontierBootNodes
	if net == TestNet {
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
		NetworkId:      int(net),
		AccountManager: accounts.NewManager(keystore),

		// Blatantly initialize the gas oracle to the defaults from go-ethereum
		GpoMinGasPrice:          new(big.Int).Mul(big.NewInt(50), common.Shannon),
		GpoMaxGasPrice:          new(big.Int).Mul(big.NewInt(500), common.Shannon),
		GpoFullBlockRatio:       80,
		GpobaseStepDown:         10,
		GpobaseStepUp:           100,
		GpobaseCorrectionFactor: 110,
	}
	// Override any default configs in the test net
	if net == TestNet {
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
	return stack, nil
}
