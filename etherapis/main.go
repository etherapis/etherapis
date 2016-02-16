package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/etherapis/etherapis/etherapis/contract"
	"github.com/etherapis/etherapis/etherapis/dashboard"
	"github.com/etherapis/etherapis/etherapis/geth"
	"github.com/etherapis/etherapis/etherapis/proxy"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	// General flags
	datadirFlag         = flag.String("datadir", "", "Path where to put the client data (\"\" = $HOME/.etherapis)")
	loglevelFlag        = flag.Int("loglevel", 3, "Log level to use for displaying system events")
	syncFlag            = flag.Duration("sync", 5*time.Minute, "Oldest allowed sync state before resync")
	dashboardFlag       = flag.Int("dashboard", 0, "Port number on which to run the dashboard (0 = disabled)")
	dashboardAssetsFlag = flag.String("dashboard-assets", "", "Path to the dashboard static assets to use (empty = built in assets)")

	// Management commands
	importFlag   = flag.String("import", "", "Path to the demo account to import")
	accountsFlag = flag.Bool("accounts", false, "Lists the available accounts for micro payments")
	serviceFlag  = flag.String("service", "", "Id of the service to check for a subscription")

	subAccFlag  = flag.Int("subacc", 0, "Own account index with which to subscribe (list with --accounts)")
	subToFlag   = flag.String("subto", "", "Id of the service to subscribe to")
	subFundFlag = flag.Float64("subfund", 1, "Initial ether value to fund the subscription with")

	// Proxy flags
	proxyFlag  = flag.String("proxy", "", "Payment proxy configs ext-port:int-port:type (e.g. 80:8080:call,81:8081:data)")
	chargeFlag = flag.Duration("charge", time.Minute, "Auto charge interval to collect pending fees")

	// Testing and admin flags
	testFlag = flag.Bool("test", false, "Runs using the default test vectors for signing and verifying signatures")
	signFlag = flag.String("sign", "", "Signs the given json input (e.g. {'provider':'0x', 'amount':0, 'nonce': 0}")
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

	// Retrieve the Ether APIs marketplace and payment contract
	var eth *eth.Ethereum
	err = client.Stack().Service(&eth)
	if err != nil {
		log15.Crit("Failed to fetch eth service", "error", err)
		return
	}
	ethContract, err := contract.New(eth.ChainDb(), eth.EventMux(), eth.BlockChain(), eth.Miner().PendingState)
	if err != nil {
		log15.Crit("Failed to get contract", "error", err)
		return
	}

	// Create the etherapis dashboard and run it
	if *dashboardFlag != 0 {
		log15.Info("Starting the EtherAPIs dashboard...", "url", fmt.Sprintf("http://localhost:%d", *dashboardFlag))
		go func() {
			http.Handle("/", dashboard.New(ethContract, eth, api, *dashboardAssetsFlag))
			if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", *dashboardFlag), nil); err != nil {
				log15.Crit("Failed to start dashboard", "error", err)
				os.Exit(-1)
			}
		}()
	}
	// Make sure we're at least semi recent on the chain before continuing
	waitSync(*syncFlag, api)

	if *signFlag != "" {
		var message struct {
			ServiceId int64
			Nonce     uint64
			Amount    uint64
		}
		if err := json.Unmarshal([]byte(*signFlag), &message); err != nil {
			log15.Crit("Failed to decode data", "error", err)
			return
		}
		fmt.Println(message)
		accounts, err := eth.AccountManager().Accounts()
		if err != nil {
			log15.Crit("Failed retrieving accounts", "err", err)
		}
		if len(accounts) == 0 {
			log15.Crit("Signing data requires at least one account", "len", len(accounts))
			return
		}
		account := accounts[0]

		from := account.Address
		serviceId := big.NewInt(message.ServiceId)

		var hash []byte
		ethContract.Call(&hash, "getHash", from, serviceId, message.Nonce, message.Amount)
		log15.Info("getting hash", "hash", common.ToHex(hash))

		eth.AccountManager().Unlock(from, "")
		sig, err := eth.AccountManager().Sign(account, hash)
		if err != nil {
			log15.Crit("signing vailed", "err", err)
			return
		}
		log15.Info("generated signature", "sig", common.ToHex(sig))

		return
	}

	// Depending on the flags, execute different things
	switch {
	case *importFlag != "":
		// Account import, parse the provided .json file and ensure it's proper
		manager := eth.AccountManager()
		account, err := manager.Import(*importFlag, "")
		if err != nil {
			log15.Crit("Failed to import specified account", "path", *importFlag, "error", err)
			return
		}
		state, _ := eth.BlockChain().State()
		log15.Info("Account successfully imported", "account", fmt.Sprintf("0x%x", account.Address), "balance", state.GetBalance(account.Address))
		return
	case *accountsFlag:
		// Account listing requested, print all accounts and balances
		accounts, err := eth.AccountManager().Accounts()
		if err != nil || len(accounts) == 0 {
			log15.Crit("Failed to retrieve account", "accounts", len(accounts), "error", err)
			return
		}
		state, _ := eth.BlockChain().State()
		for i, account := range accounts {
			balance := float64(new(big.Int).Div(state.GetBalance(account.Address), common.Finney).Int64()) / 1000
			fmt.Printf("Account #%d: %f ether (http://testnet.etherscan.io/address/0x%x)\n", i, balance, account.Address)
		}
		return

	case len(*serviceFlag) > 0:
		// Check whether any of our accounts are subscribed to this particular service
		accounts, err := eth.AccountManager().Accounts()
		if err != nil || len(accounts) == 0 {
			log15.Crit("Failed to retrieve account", "accounts", len(accounts), "error", err)
			return
		}
		serviceId := common.String2Big(*serviceFlag)
		log15.Info("Checking subscription status", "service", fmt.Sprintf("%s", serviceId))
		for i, account := range accounts {
			// Check if a subscription exists
			if !ethContract.Exists(account.Address, serviceId) {
				fmt.Printf("Account #%d: [0x%x]: not subscribed.\n", i, account.Address)
				continue
			}
			// Retrieve the current balance on the subscription
			var ethers *big.Int
			ethContract.Call(&ethers, "getSubscriptionValue", ethContract.SubscriptionId(account.Address, serviceId))
			funds := float64(new(big.Int).Div(ethers, common.Finney).Int64()) / 1000

			fmt.Printf("Account #%d: [0x%x]: subscribed, with %v ether(s) left.\n", i, account.Address, funds)
		}
		return

	case len(*subToFlag) > 0:
		// Subscription requested, make sure all the details are provided
		accounts, err := eth.AccountManager().Accounts()
		if err != nil || len(accounts) < *subAccFlag {
			log15.Crit("Failed to retrieve account", "accounts", len(accounts), "requested", *subAccFlag, "error", err)
			return
		}
		account := accounts[*subAccFlag]

		// Check if a subscription exists
		serviceId := common.String2Big(*subToFlag)
		if ethContract.Exists(account.Address, serviceId) {
			log15.Error("Account already subscribed", "index", *subAccFlag, "account", fmt.Sprintf("0x%x", account.Address), "service", serviceId)
			return
		}
		// Try to subscribe and wait until it completes
		keystore := client.Keystore()
		key, err := keystore.GetKey(account.Address, "")
		if err != nil {
			log15.Crit("Failed to unlock account", "account", fmt.Sprintf("0x%x", account.Address), "error", err)
			return
		}
		amount := new(big.Int).Mul(big.NewInt(int64(1000000000**subFundFlag)), common.Shannon)

		log15.Info("Subscribing to new payment channel", "account", fmt.Sprintf("0x%x", account.Address), "service", serviceId, "ethers", *subFundFlag)
		pend := make(chan *contract.Subscription)
		tx, err := ethContract.Subscribe(key.PrivateKey, serviceId, amount, big.NewInt(1), func(sub *contract.Subscription) { pend <- sub })
		if err != nil {
			log15.Crit("Failed to create subscription", "error", err)
			return
		}
		if err := eth.TxPool().Add(tx); err != nil {
			log15.Crit("Failed to execute subscription", "error", err)
			return
		}
		log15.Info("Waiting for subscription to be finalized...", "tx", tx)
		log15.Info("Successfully subscribed", "channel", fmt.Sprintf("%x", (<-pend).Id))
		return
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

		serviceId := big.NewInt(0)
		log15.Info("making channel name...", "from", from.Hex(), "service-id", serviceId, "id", ethContract.SubscriptionId(from, serviceId).Hex())

		if !ethContract.Exists(from, serviceId) {
			log15.Crit("No subscription found")
			return
		}
		log15.Info("checking existence...", "exists", "OK")

		amount := big.NewInt(1)

		var hash []byte
		ethContract.Call(&hash, "getHash", from, serviceId, 1, amount)
		log15.Info("signing data", "service-id", serviceId, "amount", amount, "hash", common.ToHex(hash))

		sig, err := eth.AccountManager().Sign(accounts[0], hash)
		if err != nil {
			log15.Crit("signing vailed", "err", err)
			return
		}
		log15.Info("verifying signature", "sig", common.ToHex(sig))

		if ethContract.ValidateSig(from, serviceId, 1, amount, sig) {
			log15.Info("signature was valid and was verified by the EVM")
		} else {
			log15.Crit("signature was invalid")
			return
		}

		log15.Info("verifying payment", "sig", common.ToHex(sig))
		if valid, _ := ethContract.Verify(from, serviceId, 1, amount, sig); valid {
			log15.Info("payment was valid and was verified by the EVM")
		} else {
			log15.Crit("payment was invalid")
			return
		}

		log15.Info("verifying invalid payment", "nonce", 2)
		if valid, _ := ethContract.Verify(from, serviceId, 2, amount, sig); valid {
			log15.Crit("payment was valid")
			return
		} else {
			log15.Info("payment was invalid")
		}
	}

	// If we're running a proxy, start processing external requests
	if *proxyFlag != "" {
		// Subscription requested, make sure all the details are provided
		accounts, err := eth.AccountManager().Accounts()
		if err != nil || len(accounts) < *subAccFlag {
			log15.Crit("Failed to retrieve account", "accounts", len(accounts), "requested", *subAccFlag, "error", err)
			return
		}
		account := accounts[*subAccFlag]
		if err := eth.AccountManager().Unlock(account.Address, ""); err != nil {
			log15.Crit("Failed to unlock provider account", "account", fmt.Sprintf("0x%x", account.Address), "error", err)
			return
		}
		log15.Info("Setuping vault...", "owner", account.Address.Hex())

		// Create the payment vault to hold the various authorizations
		vault := proxy.NewVault(NewCharger(account, eth.TxPool(), ethContract, eth.AccountManager()))
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
			gateway := proxy.New(i, extPort, intPort, kind, ethContract, vault)
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
	// If the dashboard was opened, wait indefinitely
	if *dashboardFlag > 0 {
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
func (v *testVerifier) Verify(from, to common.Address, amount *big.Int, signature []byte) (bool, bool) {
	if len(signature) == 0 {
		return false, false
	}
	if amount.Cmp(big.NewInt(10000000)) > 0 {
		return true, false
	}
	return true, true
}

type testCharger struct{}

func (c *testCharger) Charge(from common.Address, serviceId *big.Int, amount *big.Int, signature []byte) (common.Hash, error) {
	return common.HexToHash("0x7426125287fe1dfa9acc6d79008f0dc9a7e0c292b3387040e37c2a71518d711a"), nil
}

type Charger struct {
	txPool         *core.TxPool
	contract       *contract.Contract
	accountManager *accounts.Manager
	signer         accounts.Account
}

func NewCharger(signer accounts.Account, txPool *core.TxPool, contract *contract.Contract, am *accounts.Manager) *Charger {
	return &Charger{txPool: txPool, contract: contract, accountManager: am, signer: signer}
}

func (c *Charger) Charge(from common.Address, serviceId *big.Int, nonce uint64, amount *big.Int, signature []byte) (common.Hash, error) {
	tx, err := c.contract.Claim(c.signer.Address, from, serviceId, nonce, amount, signature)
	if err != nil {
		return common.Hash{}, err
	}

	sig, err := c.accountManager.Sign(c.signer, tx.SigHash().Bytes())
	if err != nil {
		return common.Hash{}, err
	}

	signedTx, err := tx.WithSignature(sig)
	if err != nil {
		return common.Hash{}, err
	}

	err = c.txPool.Add(signedTx)
	if err != nil {
		return common.Hash{}, err
	}

	return signedTx.Hash(), nil
}
