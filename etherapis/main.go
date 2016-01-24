package main

import (
	"flag"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/eth"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
	"github.com/gophergala2016/etherapis/etherapis/channels"
	"github.com/gophergala2016/etherapis/etherapis/geth"
	"github.com/gophergala2016/etherapis/etherapis/proxy"
)

var (
	datadirFlag  = flag.String("datadir", "", "Path where to put the client data (\"\" = $HOME/.etherapis)")
	loglevelFlag = flag.Int("loglevel", 3, "Log level to use for displaying system events")
	syncFlag     = flag.Duration("sync", 5*time.Minute, "Oldest allowed sync state before resync")
	proxyFlag    = flag.String("proxy", "", "Payment proxy configs ext-port:int-port:type (e.g. 80:8080:call,81:8081:data)")
	chargeFlag   = flag.Duration("charge", time.Minute, "Auto charge interval to collect pending fees")
	testFlag     = flag.Bool("test", false, "Runs using the default test vectors for signing and verifying signatures")
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

	var eth *eth.Ethereum
	err = client.Stack().Service(&eth)
	if err != nil {
		log15.Crit("Failed to fetch eth service", "error", err)
		return
	}
	channels, err := channels.Fetch(eth.ChainDb(), eth.EventMux(), eth.BlockChain())
	if err != nil {
		log15.Crit("Failed to get contract", "error", err)
	}

	if *testFlag {
		accounts, err := eth.AccountManager().Accounts()
		if err != nil {
			log15.Crit("Failed retrieving accounts", "err", err)
		}
		if len(accounts) < 2 {
			log15.Crit("Test vectors requires at least 2 accounts", "len", len(accounts))
			return
		}

		log15.Info("Attempting channel test vectors...")
		from := accounts[0].Address
		eth.AccountManager().Unlock(from, "")

		to := accounts[1].Address
		log15.Info("making channel name...", "from", from.Hex(), "to", to.Hex(), "ID", common.ToHex(channels.ChannelId(from, to)))
		log15.Info("checking existence...", "exists", channels.Exists(from, to))

		amount := big.NewInt(100)
		hash := channels.Call("getHash", from, to, 0, amount).([]byte)
		log15.Info("signing data", "to", to.Hex(), "amount", amount, "hash", common.ToHex(hash))

		sig, err := eth.AccountManager().Sign(accounts[0], hash)
		if err != nil {
			log15.Crit("signing vailed", "err", err)
			return
		}
		log15.Info("verifying signature", "sig", common.ToHex(sig))

		if channels.Validate(from, to, amount, sig) {
			log15.Info("signature was valid and was verified by the EVM")
		} else {
			log15.Crit("signature was invalid")
		}
	}

	// If we're running a proxy, start processing external requests
	if *proxyFlag != "" {
		// Create the payment vault to hold the various authorizations
		vault := proxy.NewVault(new(testCharger))
		vault.AutoCharge(*chargeFlag)

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

type testCharger struct{}

func (c *testCharger) Charge(from, to common.Address, amount uint64, signature []byte) (string, error) {
	return "0x7426125287fe1dfa9acc6d79008f0dc9a7e0c292b3387040e37c2a71518d711a", nil
}
