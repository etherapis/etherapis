package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
	"github.com/gophergala2016/etherapis/etherapis/geth"
)

var (
	datadirFlag  = flag.String("datadir", "", "Path where to put the client data (\"\" = $HOME/.etherapis)")
	loglevelFlag = flag.Int("loglevel", 3, "Log level to use for displaying system events")
)

func main() {
	// Parse and handle the command line flags
	flag.Parse()

	log15.Root().SetHandler(log15.LvlFilterHandler(log15.Lvl(*loglevelFlag), log15.StderrHandler))

	datadir := *datadirFlag
	if datadir == "" {
		datadir = filepath.Join(os.Getenv("HOME"), ".etherapis")
	}
	if err := os.MkdirAll(datadir, 0700); err != nil {
		log15.Crit("Failed to create data directory: %v", err)
	}
	// Assemble and start the Ethereum client
	log15.Info("Booting Ethereum client...")
	client, err := geth.New(datadir, geth.TestNet)
	if err != nil {
		log15.Crit("Failed to create Ethereum client", "error", err)
	}
	if err := client.Start(); err != nil {
		log15.Crit("Failed to start Ethereum client", "error", err)
	}
	api, err := client.Attach()
	if err != nil {
		log15.Crit("Failed to attach to node", "error", err)
	}
	// Wait for network connectivity and monitor synchronization
	log15.Info("Searching for network peers...")
	server := client.Stack().Server()
	for len(server.Peers()) == 0 {
		time.Sleep(100 * time.Millisecond)
	}
	go monitorSync(api)

	// Make sure we're at least semi recent on the chain before continuing
	waitSync(5*time.Minute, api)
	log15.Info("Connected and synced with the network!")

	// Clean up for now
	log15.Info("Terminating Ethereum client...")
	if err := client.Stop(); err != nil {
		log15.Crit("Failed to terminate Ethereum client", "error", err)
	}
}
