// Contains the Go implementation of the Ethereum RPC specs.

package geth

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// API is a Go wrapper around the JSON RPC API exposed by a Geth node.
//
// In theory implements:
//  - https://github.com/ethereum/wiki/wiki/JSON-RPC
//  - https://github.com/ethereum/go-ethereum/wiki/Go-ethereum-management-API's
//
// In practice it implements the bare essentials needed for this hackathon :D
type API struct {
	client rpc.Client // RPC client with a live connection to an Ethereum node
	autoid uint32     // ID number to use for the next API request
	lock   sync.Mutex // Singleton access until we get to request multiplexing
}

// request is a JSON RPC request package assembled internally from the client
// method calls.
type request struct {
	JsonRpc string        `json:"jsonrpc"` // Version of the JSON RPC protocol, always set to 2.0
	Id      int           `json:"id"`      // Auto incrementing ID number for this request
	Method  string        `json:"method"`  // Remote procedure name to invoke on the server
	Params  []interface{} `json:"params"`  // List of parameters to pass through (keep types simple)
}

// response is a JSON RPC response package sent back from the API server.
type response struct {
	JsonRpc string          `json:"jsonrpc"` // Version of the JSON RPC protocol, always set to 2.0
	Id      int             `json:"id"`      // Auto incrementing ID number for this request
	Error   json.RawMessage `json:"error"`   // Any error returned by the remote side
	Result  json.RawMessage `json:"result"`  // Whatever the remote side sends us in reply
}

// request forwards an API request to the RPC server, and parses the response.
//
// This is currently painfully non-concurrent, but it will have to do until we
// find the time for niceties like this :P
func (api *API) Request(method string, params []interface{}) (json.RawMessage, error) {
	api.lock.Lock()
	defer api.lock.Unlock()

	// Ugly hack to serialize an empty list properly
	if params == nil {
		params = []interface{}{}
	}
	// Assemble the request object
	req := &request{
		JsonRpc: "2.0",
		Id:      int(atomic.AddUint32(&api.autoid, 1)),
		Method:  method,
		Params:  params,
	}
	if err := api.client.Send(req); err != nil {
		return nil, err
	}
	res := new(response)
	if err := api.client.Recv(res); err != nil {
		return nil, err
	}
	if len(res.Error) > 0 {
		return nil, fmt.Errorf("remote error: %s", string(res.Error))
	}
	return res.Result, nil
}

// BlockNumber retrieves the current head number of the blockchain.
func (api *API) BlockNumber() (uint64, error) {
	res, err := api.Request("eth_blockNumber", nil)
	if err != nil {
		return 0, err
	}
	var hex string
	if err := json.Unmarshal(res, &hex); err != nil {
		return 0, err
	}
	return new(big.Int).SetBytes(common.FromHex(hex)).Uint64(), nil
}

// GetBlockTime retrieves the block of the given number from the canonical chain.
func (api *API) GetBlockTime(num uint64) (time.Time, error) {
	res, err := api.Request("eth_getBlockByNumber", []interface{}{rpc.NewHexNumber(num), true})
	if err != nil {
		return time.Time{}, err
	}
	result := make(map[string]json.RawMessage)
	if err := json.Unmarshal(res, &result); err != nil {
		return time.Time{}, err
	}
	var hex string
	if err := json.Unmarshal(result["timestamp"], &hex); err != nil {
		return time.Time{}, err
	}
	return time.Unix(new(big.Int).SetBytes(common.FromHex(hex)).Int64(), 0), nil
}

// SyncStatus is the current state the network sync is in.
type SyncStatus struct {
	StartingBlock uint64
	CurrentBlock  uint64
	HighestBlock  uint64
}

// Syncing returns the current sync status of the node, or nil if the node is not
// currently synchronizing with the network.
func (api *API) Syncing() (*SyncStatus, error) {
	// Execute the request and check if syncing is not running
	res, err := api.Request("eth_syncing", nil)
	if err != nil {
		return nil, err
	}
	var running bool
	if err := json.Unmarshal(res, &running); err == nil {
		return nil, nil
	}
	// Sync is running, extract the current status
	result := make(map[string]string)
	if err := json.Unmarshal(res, &result); err != nil {
		return nil, err
	}
	return &SyncStatus{
		StartingBlock: new(big.Int).SetBytes(common.FromHex(result["startingBlock"])).Uint64(),
		CurrentBlock:  new(big.Int).SetBytes(common.FromHex(result["currentBlock"])).Uint64(),
		HighestBlock:  new(big.Int).SetBytes(common.FromHex(result["highestBlock"])).Uint64(),
	}, nil
}

// Accounts retrieves the currently available Ethereum accounts.
func (api *API) Accounts() ([]common.Address, error) {
	return nil, nil
}
