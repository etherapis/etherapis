package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/eth"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
	"github.com/gophergala2016/etherapis/etherapis/geth"
)

func main() {
	datadir, err := ioutil.TempDir("", "etherapis-")
	if err != nil {
		log15.Crit("Failed to create temporary datadir", "error", err)
	}
	defer os.RemoveAll(datadir)

	log15.Info("Booting Ethereum client...")
	client, err := geth.New(datadir, geth.TestNet)
	if err != nil {
		log15.Crit("Failed to create Ethereum client", "error", err)
	}
	if err := client.Start(); err != nil {
		log15.Crit("Failed to start Ethereum client", "error", err)
	}

	log15.Info("Searching for network peers...")
	ethereum := new(eth.Ethereum)
	if err := client.Stack().Service(&ethereum); err != nil {
		log15.Crit("Failed to retrieve Ethereum service", "error", err)
	}
	server := client.Stack().Server()
	for len(server.Peers()) == 0 {
		time.Sleep(time.Second)
	}
	log15.Info("Connected to the network, syncing...")
	for {
		head := ethereum.BlockChain().CurrentFastBlock()
		log15.Info("Synchronizing network...", "peers", len(server.Peers()), "block", head.NumberU64(), "hash", fmt.Sprintf("%x", head.Hash().Bytes()[:4]))
		time.Sleep(time.Second)
	}
	log15.Info("Terminating Ethereum client...")
	if err := client.Stop(); err != nil {
		log15.Crit("Failed to terminate Ethereum client", "error", err)
	}
}
