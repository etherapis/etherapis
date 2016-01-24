package main

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
	"github.com/gophergala2016/etherapis/etherapis/geth"
	"github.com/gophergala2016/etherapis/etherapis/proxy"
)

var (
	datadirFlag  = flag.String("datadir", "", "Path where to put the client data (\"\" = $HOME/.etherapis)")
	loglevelFlag = flag.Int("loglevel", 3, "Log level to use for displaying system events")
	syncFlag     = flag.Duration("sync", 5*time.Minute, "Oldest allowed sync state before resync")
	proxyFlag    = flag.String("proxy", "", "Payment proxy configs ext-port:int-port:type (e.g. 80:8080:call,81:8081:data)")
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
		return
	}
	// Assemble and start the Ethereum client
	log15.Info("Booting Ethereum client...")
	client, err := geth.New(datadir, geth.TestNet)
	if err != nil {
		log15.Crit("Failed to create Ethereum client", "error", err)
		return
	}
	if err := client.Start(); err != nil {
		log15.Crit("Failed to start Ethereum client", "error", err)
		return
	}
	api, err := client.Attach()
	if err != nil {
		log15.Crit("Failed to attach to node", "error", err)
		return
	}
	// Wait for network connectivity and monitor synchronization
	log15.Info("Searching for network peers...")
	server := client.Stack().Server()
	for len(server.Peers()) == 0 {
		time.Sleep(100 * time.Millisecond)
	}
	go monitorSync(api)

	// Make sure we're at least semi recent on the chain before continuing
	waitSync(*syncFlag, api)

	// If we're running a proxy, start processing external requests
	if *proxyFlag != "" {
		// Create the payment vault to hold the various authorizations
		vault := proxy.NewVault()

		for i, config := range strings.Split(*proxyFlag, ",") {
			// Split the proxy configuration
			parts := strings.Split(config, ":")
			if len(parts) != 3 {
				log15.Crit("Invalid proxy config", "config", config)
				return
			}
			extPort, err := strconv.Atoi(parts[0])
			if err != nil || extPort < 0 || extPort > 65535 {
				log15.Crit("Invalid external port number", "port", parts[0])
				return
			}
			intPort, err := strconv.Atoi(parts[1])
			if err != nil || intPort < 0 || intPort > 65535 {
				log15.Crit("Invalid internal port number", "port", parts[1])
				return
			}
			var kind proxy.ProxyType
			switch strings.ToLower(parts[2]) {
			case "call":
				kind = proxy.CallProxy
			case "data":
				kind = proxy.DataProxy
			default:
				log15.Crit("Unsupported proxy type", "type", parts[2], "allowed", []string{"call", "data"})
				return
			}
			// Create and start the new proxy
			gateway := proxy.New(i, extPort, intPort, kind, new(testVerifier), vault)
			go func() {
				if err := gateway.Start(); err != nil {
					log15.Crit("Failed to start proxy", "error", err)
					os.Exit(-1)
				}
			}()
		}
		// Wait indefinitely, for now at least
		for {
			time.Sleep(time.Second)
		}
	}
	// Clean up for now
	log15.Info("Terminating Ethereum client...")
	if err := client.Stop(); err != nil {
		log15.Crit("Failed to terminate Ethereum client", "error", err)
		return
	}
}

type testVerifier struct{}

func (v *testVerifier) Exists(from, to common.Address) bool { return from != to }
func (v *testVerifier) Verify(from, to common.Address, amount uint64, signature []byte) (bool, bool) {
	if len(signature) == 0 {
		return false, false
	}
	if amount > 10000000 {
		return true, false
	}
	return true, true
}
