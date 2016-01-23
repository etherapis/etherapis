// Package geth implements the wrapper around to go-ethereum client
package geth

import (
	"fmt"
	"math/big"
	"net"
	"path/filepath"
	"runtime"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/accounts"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/core"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/core/state"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/crypto"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/eth"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/node"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/rpc"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
)

// Geth is a wrapper around the Ethereum Go client.
type Geth struct {
	stack *node.Node // Ethereum network node / protocol stack

	ipcEndpoint string       // File system (or named pipe) path for the
	ipcSocket   net.Listener // IPC listener socket waiting for inbound connections
	ipcServer   *rpc.Server  // RPC server handling the IPC API requests
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
		stack: stack,
	}, nil
}

// Start boots up the Ethereum protocol, starts interacting with the P2P network
// and opens up the IPC based JSOn RPC API for accessing the exposed APIs.
func (g *Geth) Start() error {
	// Assemble and start the protocol stack
	err := g.stack.Start()
	if err != nil {
		return err
	}
	// Attach an IPC endpoint to the stack (this should really have been done by the stack..)
	if runtime.GOOS == "windows" {
		g.ipcEndpoint = common.DefaultIpcPath()
	} else {
		g.ipcEndpoint = filepath.Join(g.stack.DataDir(), "geth.ipc")
	}
	if g.ipcSocket, err = rpc.CreateIPCListener(g.ipcEndpoint); err != nil {
		return err
	}
	// Iterate over all the APIs provided by the stack and expose them
	g.ipcServer = rpc.NewServer()
	for _, api := range g.stack.APIs() {
		log15.Debug("Exposing/extending API...", "namespace", api.Namespace, "runner", fmt.Sprintf("%T", api.Service))
		g.ipcServer.RegisterName(api.Namespace, api.Service)
	}
	log15.Info("Starting IPC listener...", "endpoint", g.ipcEndpoint)
	go g.handleIPC()

	return nil
}

// Stop closes down the Ethereum client, along with any other resouces it might
// be keeping around.
func (g *Geth) Stop() error {
	// Silently ignore errors for this hackathon :P
	g.ipcSocket.Close()
	g.ipcServer.Stop()

	return g.stack.Stop()
}

// Stack is a quick hack to expose the internal Ethereum node implementation. We
// should probably remove this after the API interface is implemented, but until
// then it makes things simpler.
func (g *Geth) Stack() *node.Node {
	return g.stack
}

// Attach connects to the running node's IPC exposed APIs, and returns an Go API
// interface.
func (g *Geth) Attach() (*API, error) {
	client, err := rpc.NewIPCClient(g.ipcEndpoint)
	if err != nil {
		return nil, err
	}
	return &API{client: client}, nil
}

// handleIPC is the RPC over IPC handler to accept inbound API requests and run
// them until termination is requested.
func (g *Geth) handleIPC() {
	for {
		conn, err := g.ipcSocket.Accept()
		if err != nil {
			log15.Debug("Unable to accept connection (shutting down?)", "error", err)
			return
		}
		go g.ipcServer.ServeCodec(rpc.NewJSONCodec(conn))
	}
}
