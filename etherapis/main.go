package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/accounts"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/core"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/core/types"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/eth"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
	"github.com/gophergala2016/etherapis/etherapis/channels"
	"github.com/gophergala2016/etherapis/etherapis/geth"
	"github.com/gophergala2016/etherapis/etherapis/proxy"
)

var (
	// General flags
	datadirFlag  = flag.String("datadir", "", "Path where to put the client data (\"\" = $HOME/.etherapis)")
	loglevelFlag = flag.Int("loglevel", 3, "Log level to use for displaying system events")
	syncFlag     = flag.Duration("sync", 5*time.Minute, "Oldest allowed sync state before resync")

	// Management commands
	importFlag   = flag.String("import", "", "Path to the demo account to import")
	accountsFlag = flag.Bool("accounts", false, "Lists the available accounts for micro payments")
	queryFlag    = flag.String("query", "", "Address of the provider to check for a subscription")

	subAccFlag  = flag.Int("subacc", 0, "Own account index with which to subscribe (list with --accounts)")
	subToFlag   = flag.String("subto", "", "Address of the provider to subscribe to")
	subFundFlag = flag.Float64("subfund", 1, "Initial ether value to fund the subscription with")

	// Proxy flags
	proxyFlag  = flag.String("proxy", "", "Payment proxy configs ext-port:int-port:type (e.g. 80:8080:call,81:8081:data)")
	chargeFlag = flag.Duration("charge", time.Minute, "Auto charge interval to collect pending fees")

	// Testing and admin flags
	testFlag    = flag.Bool("test", false, "Runs using the default test vectors for signing and verifying signatures")
	accGenFlag  = flag.Int("gen", 0, "Generates a batch of (empty) demo accounts and dumps their keys")
	accLiveFlag = flag.Bool("live", false, "Requests live account generation (funded and uloaded)")
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
	contract, err := channels.Fetch(eth.ChainDb(), eth.EventMux(), eth.BlockChain())
	if err != nil {
		log15.Crit("Failed to get contract", "error", err)
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

	case *accGenFlag > 0:
		// We're generating and dumping demo accounts
		var nonce uint64
		var bank accounts.Account

		if *accLiveFlag {
			// If we want to fund generated accounts, make sure we can
			accounts, err := eth.AccountManager().Accounts()
			if err != nil || len(accounts) == 0 {
				log15.Crit("Failed to retrieve funding account", "accounts", len(accounts), "error", err)
				return
			}
			bank = accounts[0]
			if err := eth.AccountManager().Unlock(bank.Address, "gophergala"); err != nil {
				log15.Crit("Failed to unlock funding account", "account", fmt.Sprintf("0x%x", bank.Address), "error", err)
				return
			}
			state, _ := eth.BlockChain().State()
			nonce = state.GetNonce(bank.Address)

			log15.Info("Funding demo accounts with", "bank", fmt.Sprintf("0x%x", bank.Address), "nonce", nonce)
		}
		// Start generating the actual accounts
		log15.Info("Generating demo accounts", "count", *accGenFlag)
		for i := 0; i < *accGenFlag; i++ {
			// Generate a new account
			account, err := eth.AccountManager().NewAccount("pass")
			if err != nil {
				log15.Crit("Failed to generate new account", "error", err)
				return
			}
			// Export it's private key
			keyPath := fmt.Sprintf("0x%x.key", account.Address)
			if err := eth.AccountManager().Export(keyPath, account.Address, "pass"); err != nil {
				log15.Crit("Failed to export account", "account", fmt.Sprintf("0x%x", account.Address), "error", err)
				return
			}
			// Clean up so it doesn't clutter out accounts
			if err := eth.AccountManager().DeleteAccount(account.Address, "pass"); err != nil {
				log15.Crit("Failed to delete account", "account", fmt.Sprintf("0x%x", account.Address), "error", err)
				return
			}
			// If we're just testing, stop here
			if !*accLiveFlag {
				log15.Info("Account generated and exported", "path", keyPath)
				continue
			}
			// Oh boy, live accounts, send some ether to it and upload to the faucet
			allowance := new(big.Int).Mul(big.NewInt(10), common.Ether)
			price := new(big.Int).Mul(big.NewInt(50), common.Shannon)

			tx := types.NewTransaction(nonce, account.Address, allowance, big.NewInt(21000), price, nil)
			sig, err := eth.AccountManager().Sign(bank, tx.SigHash().Bytes())
			if err != nil {
				log15.Crit("Failed to sign funding transaction", "error", err)
				return
			}
			stx, err := tx.WithSignature(sig)
			if err != nil {
				log15.Crit("Failed to assemble funding transaction", "error", err)
				return
			}
			if err := eth.TxPool().Add(stx); err != nil {
				log15.Crit("Failed to execute transfer", "error", err)
				return
			}
			nonce++
			log15.Info("Account successfully funded", "account", fmt.Sprintf("0x%x", account.Address))

			// Upload the account to the faucet server
			key, err := ioutil.ReadFile(keyPath)
			if err != nil {
				log15.Crit("Failed to load private key", "error", err)
				return
			}
			res, err := http.Get("https://etherapis.appspot.com/faucet/fund?key=" + string(key))
			if err != nil {
				log15.Crit("Failed to upload private key to faucet", "error", err)
				return
			}
			res.Body.Close()

			log15.Info("Account uploaded to faucet", "account", fmt.Sprintf("0x%x", account.Address))
			os.Remove(keyPath)
		}
		// Just wait a bit to ensure transactions get propagated into the network
		log15.Info("Sleeping to ensure transaction propagation")
		time.Sleep(10 * time.Second)
		return

	case len(*queryFlag) > 0:
		// Check whether any of our accounts are subscribed to this particular service
		accounts, err := eth.AccountManager().Accounts()
		if err != nil || len(accounts) == 0 {
			log15.Crit("Failed to retrieve account", "accounts", len(accounts), "error", err)
			return
		}
		provider := common.HexToAddress(*queryFlag)
		log15.Info("Checking subscription status", "service", fmt.Sprintf("0x%x", provider))
		for i, account := range accounts {
			// Check if a subscription exists
			if !contract.Exists(account.Address, provider) {
				fmt.Printf("Account #%d: [0x%x]: not subscribed.\n", i, account.Address)
				continue
			}
			// Retrieve the current balance on the subscription
			ethers := contract.Call("getChannelValue", contract.ChannelId(account.Address, provider)).(*big.Int)
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
		provider := common.HexToAddress(*subToFlag)
		if contract.Exists(account.Address, provider) {
			log15.Error("Account already subscribed", "index", *subAccFlag, "account", fmt.Sprintf("0x%x", account.Address), "service", fmt.Sprintf("0x%x", provider))
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

		log15.Info("Subscribing to new payment channel", "account", fmt.Sprintf("0x%x", account.Address), "service", fmt.Sprintf("0x%x", provider), "ethers", *subFundFlag)
		pend := make(chan *channels.Channel)
		tx, err := contract.NewChannel(key.PrivateKey, provider, amount, common.Shannon, func(sub *channels.Channel) { pend <- sub })
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

		to := accounts[1].Address
		log15.Info("making channel name...", "from", from.Hex(), "to", to.Hex(), "ID", contract.ChannelId(from, to).Hex())
		log15.Info("checking existence...", "exists", contract.Exists(from, to))

		amount := big.NewInt(10)
		hash := contract.Call("getHash", from, to, 0, amount).([]byte)
		log15.Info("signing data", "to", to.Hex(), "amount", amount, "hash", common.ToHex(hash))

		sig, err := eth.AccountManager().Sign(accounts[0], hash)
		if err != nil {
			log15.Crit("signing vailed", "err", err)
			return
		}
		log15.Info("verifying signature", "sig", common.ToHex(sig))

		if contract.ValidateSig(from, to, 0, amount, sig) {
			log15.Info("signature was valid and was verified by the EVM")
		} else {
			log15.Crit("signature was invalid")
		}

		log15.Info("verifying payment", "sig", common.ToHex(sig))
		if contract.Validate(from, to, 0, amount, sig) {
			log15.Info("payment was valid and was verified by the EVM")
		} else {
			log15.Crit("payment was invalid")
		}

		log15.Info("verifying invalid payment", "nonce", 1)
		if contract.Validate(from, to, 1, amount, sig) {
			log15.Crit("payment was valid")
		} else {
			log15.Info("payment was invalid")
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

func (c *testCharger) Charge(from, to common.Address, amount *big.Int, signature []byte) (common.Hash, error) {
	return common.HexToHash("0x7426125287fe1dfa9acc6d79008f0dc9a7e0c292b3387040e37c2a71518d711a"), nil
}

type Charger struct {
	txPool         *core.TxPool
	channels       *channels.Channels
	accountManager *accounts.Manager
	signer         accounts.Account
}

func NewCharger(signer accounts.Account, txPool *core.TxPool, channels *channels.Channels, am *accounts.Manager) *Charger {
	return &Charger{txPool: txPool, channels: channels, accountManager: am, signer: signer}
}

func (c *Charger) Charge(from, to common.Address, amount *big.Int, signature []byte) (common.Hash, error) {
	tx, err := c.channels.Claim(c.signer.Address, from, to, amount, signature)
	if err != nil {
		return common.Hash{}, err
	}

	sig, err := c.accountManager.Sign(c.signer, tx.Hash().Bytes())
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
