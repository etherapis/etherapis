package main

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/core/state"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/crypto"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/ethdb"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/event"
)

func TestContract(t *testing.T) {
	var mux event.TypeMux

	db, _ := ethdb.NewMemDatabase()
	statedb, _ := state.New(common.Hash{}, db)
	stateFn := func() *state.StateDB {
		return statedb
	}

	key1, _ := crypto.GenerateKey()
	key2, _ := crypto.GenerateKey()
	to := crypto.PubkeyToAddress(key2.PublicKey)

	contract, err := GetContract(db, &mux, stateFn)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("methods:", contract.abi.Methods)
	fmt.Println("events:", contract.abi.Events)
	fmt.Println("\n\n")

	tx, err := contract.NewChannel(key1, to, new(big.Int).Mul(big.NewInt(10), common.Ether), func(c *Channel) {
		fmt.Println("new  channel created", c)
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tx)
}
