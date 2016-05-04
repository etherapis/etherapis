package main

import (
	"errors"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/etherapis/etherapis/etherapis"
	"github.com/etherapis/etherapis/etherapis/contract"
	"github.com/etherapis/etherapis/etherapis/dashboard"
	"github.com/etherapis/etherapis/etherapis/geth"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/logger/glog"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	// General flags
	datadirFlag         = flag.String("datadir", "", "Path where to put the client data (\"\" = $HOME/.etherapis)")
	loglevelFlag        = flag.Int("loglevel", 3, "Log level to use for displaying system events")
	loglevelGethFlag    = flag.Int("loglevelgeth", 0, "Log level to use for displaying go-ethereum system events")
	dashboardFlag       = flag.Int("dashboard", 8080, "Port number on which to run the dashboard (0 = disabled)")
	dashboardAssetsFlag = flag.String("dashboard-assets", "", "Path to the dashboard static assets to use (empty = built in assets)")
	passwordFlag        = flag.String("password", "", "Master password to use for account management")

	// Proxy flags
	proxyFlag  = flag.String("proxy", "", "Payment proxy configs ext-port:int-port:type (e.g. 80:8080:call,81:8081:data)")
	chargeFlag = flag.Duration("charge", time.Minute, "Auto charge interval to collect pending fees")

	// Testing and admin flags
	signFlag   = flag.String("sign", "", "Signs the given json input (e.g. {'provider':'0x', 'amount':0, 'nonce': 0}")
	deployFlag = flag.Bool("deploy", false, "Deploys a new version of the EtherAPIs contract")
)

func main() {
	// Parse and handle the command line flags
	flag.Parse()

	log15.Root().SetHandler(log15.LvlFilterHandler(log15.Lvl(*loglevelFlag), log15.StderrHandler))
	if *loglevelGethFlag > 0 {
		glog.SetV(*loglevelGethFlag)
		glog.SetToStderr(true)
	}
	datadir := *datadirFlag
	if datadir == "" {
		datadir = filepath.Join(os.Getenv("HOME"), ".etherapis")
	}
	if err := os.MkdirAll(datadir, 0700); err != nil {
		log15.Crit("Failed to create data directory: %v", err)
		return
	}
	// Start the Ether APIs client and unlock all used accounts
	log15.Info("Joining Ethereum network...")
	client, err := etherapis.New(datadir, geth.TestNet, common.HexToAddress("0x1b523270eac78b07cb6170c11abff9a1df39ca20"))
	if err != nil {
		log15.Crit("Failed to create Ether APIs client", "error", err)
		return
	}
	log15.Info("Unlocking Ether APIs accounts...")
	if err := client.Unlock(*passwordFlag); err != nil {
		log15.Crit("Failed to unlock accounts", "error", err)
		return
	}
	// Deploy a new contract if it was requested
	if *deployFlag {
		accounts, _ := client.ListAccounts()
		if len(accounts) == 0 {
			log15.Crit("Cannot deploy new contract without a valid account")
			return
		}
		log15.Warn("Deploying new EtherAPIs contract...", "owner", accounts[0].Hex())
		address, tx, err := client.Deploy(accounts[0])
		if err != nil {
			log15.Crit("Failed to deploy new contract", "error", err)
			return
		}
		log15.Warn("New contract deployed", "address", address.Hex(), "transaction", tx.Hash().Hex())
	}
	// Create the etherapis dashboard and run it
	if *dashboardFlag != 0 {
		log15.Info("Starting the EtherAPIs dashboard...", "url", fmt.Sprintf("http://localhost:%d", *dashboardFlag))
		go func() {
			http.Handle("/", dashboard.New(client, *dashboardAssetsFlag))
			if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", *dashboardFlag), nil); err != nil {
				log15.Crit("Failed to start dashboard", "error", err)
				os.Exit(-1)
			}
		}()
	}
	// Some leftovers from the gopher gala
	/*if *signFlag != "" {
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
	}*/
	// If the dashboard was opened, wait indefinitely
	if *dashboardFlag > 0 {
		for {
			time.Sleep(5 * time.Second)
		}
	}
	// Clean up for now
	log15.Info("Disconnecting from Ethereum network...")
	if err := client.Close(); err != nil {
		log15.Crit("Failed to terminate Ether APIs client", "error", err)
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
	contract       *contract.EtherAPIs
	accountManager *accounts.Manager
	signer         accounts.Account
}

func NewCharger(signer accounts.Account, txPool *core.TxPool, contract *contract.EtherAPIs, am *accounts.Manager) *Charger {
	return &Charger{txPool: txPool, contract: contract, accountManager: am, signer: signer}
}

func (c *Charger) Charge(from common.Address, serviceId *big.Int, nonce uint64, amount *big.Int, signature []byte) (common.Hash, error) {
	/*	tx, err := c.contract.Claim(c.signer.Address, from, serviceId, nonce, amount, signature)
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

		return signedTx.Hash(), nil*/
	return common.Hash{}, errors.New("not implemented")
}
