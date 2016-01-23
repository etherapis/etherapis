package main

import (
	"strings"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"
)

type Contract struct {
	abi abi.ABI
}

// contractAddress is the static address on which the contract resides
var contractAddress = common.HexToAddress("0xaa")

// GetContract initialises a new abi and returns the contract. It does not
// deploy the contract, hence the name.
func GetContract() *Contract {
	var contract Contract

	var err error
	contract.abi, err = abi.JSON(strings.NewReader(jsonAbi))
	if err != nil {
		panic(err)
	}

	return &contract
}

const jsonAbi = `[{"constant":false,"inputs":[],"name":"Channel","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"channel","type":"bytes32"}],"name":"isValidChannel","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"channel","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"claim","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"}],"name":"createChannel","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"channel","type":"bytes32"},{"name":"recipient","type":"address"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"}],"name":"getHash","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"channels","outputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"validUntil","type":"uint256"},{"name":"valid","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"channel","type":"bytes32"}],"name":"getChannelValidUntil","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"channel","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"verify","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"channel","type":"bytes32"}],"name":"reclaim","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"channel","type":"bytes32"}],"name":"getChannelOwner","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[{"name":"channel","type":"bytes32"}],"name":"deposit","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"channel","type":"bytes32"}],"name":"getChannelValue","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":false,"name":"channel","type":"bytes32"}],"name":"NewChannel","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"channel","type":"bytes32"}],"name":"Deposit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"who","type":"address"},{"indexed":true,"name":"channel","type":"bytes32"}],"name":"Claim","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"channel","type":"bytes32"}],"name":"Reclaim","type":"event"}]`
