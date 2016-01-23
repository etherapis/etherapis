package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/eth"
	"github.com/gophergala2016/etherapis/etherapis/geth"
)

func main() {
	datadir, err := ioutil.TempDir("", "etherapis-")
	if err != nil {
		log.Fatalf("Failed to create temporary datadir: %v", err)
	}
	defer os.RemoveAll(datadir)

	log.Printf("Booting Ethereum client...")
	client, err := geth.New(datadir, geth.TestNet)
	if err != nil {
		log.Fatalf("Failed to create Ethereum client: %v", err)
	}
	if err := client.Start(); err != nil {
		log.Fatalf("Failed to start Ethereum client: %v", err)
	}

	log.Printf("Searching for network peers...")
	ethereum := new(eth.Ethereum)
	if err := client.Service(&ethereum); err != nil {
		log.Fatalf("Failed to retrieve Ethereum service: %v", err)
	}
	for len(client.Server().Peers()) == 0 {
		time.Sleep(time.Second)
	}
	log.Printf("Connected to the network, syncing...")
	for {
		head := ethereum.BlockChain().CurrentFastBlock()
		log.Printf("At block #%d [%x]", head.NumberU64(), head.Hash().Bytes()[:4])
		time.Sleep(time.Second)
	}
	log.Printf("Terminating Ethereum client...")
	if err := client.Stop(); err != nil {
		log.Fatalf("Failed to terminate Ethereum client: %v", err)
	}
}
