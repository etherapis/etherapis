// Contains the payment vault which accumulates the previously accepted payment
// to be able to charge them upon request.

package proxy

import (
	"math/big"
	"sync"
	"time"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
)

// Vault is a payment aggregator collecting the individual accepted payments made
// by various clients. It can at any point return the most recent payment + proof,
// and can also charge the payments via the Ethereum network.
type Vault struct {
	vaults  map[common.Address]*accountVault // Individual vaults per provider account
	charger Charger                          // Paymen charger to redeem pending payments
	lock    sync.RWMutex                     // Lock protecting the vaults
}

// NewVault create and returns an empty vault ready for accepting payments.
func NewVault(charger Charger) *Vault {
	return &Vault{
		vaults:  make(map[common.Address]*accountVault),
		charger: charger,
	}
}

// Store inserts a new authorized payment into the vault for later redemption.
func (v *Vault) Store(auth *authorization) {
	v.lock.Lock()
	defer v.lock.Unlock()

	provider := common.HexToAddress(auth.Provider)
	if _, ok := v.vaults[provider]; !ok {
		v.vaults[provider] = &accountVault{
			auths: make(map[common.Address]*authorization),
			pends: make(map[common.Address]*authorization),
		}
	}
	v.vaults[provider].Store(auth)
}

// Fetch retrieves the last known authorization made by a particular consumer to
// a particular provider.
func (v *Vault) Fetch(provider, consumer common.Address) *authorization {
	v.lock.RLock()
	defer v.lock.RUnlock()

	if account, ok := v.vaults[provider]; ok {
		return account.Fetch(consumer)
	}
	return nil
}

// Charge will redeem all pending payments from the vault.
func (v *Vault) Charge() {
	v.lock.RLock()
	defer v.lock.RUnlock()

	for _, account := range v.vaults {
		account.Charge(v.charger)
	}
}

// AutoCharge starts a periodical automatic charging to claim collected funds.
func (v *Vault) AutoCharge(interval time.Duration) {
	log15.Info("Automatic payment charging enabled", "interval", interval)
	go func() {
		for {
			time.Sleep(interval)
			log15.Debug("Running automatic payment charging...")
			v.Charge()
			log15.Debug("Rescheduling automatic charges...", "at", time.Now().Add(interval))
		}
	}()
}

// accountVault is the payment aggregator for a single owned account.
type accountVault struct {
	auths map[common.Address]*authorization // Individual authorizations per client
	pends map[common.Address]*authorization // Authorizations that have not yet been charged

	lock sync.RWMutex // Lock protecting the authorizations
}

// Store inserts a new authorized payment into the account vault for later
// redemption.
func (v *accountVault) Store(auth *authorization) {
	v.lock.Lock()
	defer v.lock.Unlock()

	consumer := common.HexToAddress(auth.Consumer)
	v.auths[consumer] = auth
	v.pends[consumer] = auth
}

// Fetch retrieves the last known authorization made by a particular consumer to
// this particular account vault.
func (v *accountVault) Fetch(consumer common.Address) *authorization {
	v.lock.RLock()
	defer v.lock.RUnlock()

	return v.auths[consumer]
}

// Charge will redeem all pending payments made to this provider since startup.
func (v *accountVault) Charge(charger Charger) {
	v.lock.RLock()
	for _, auth := range v.pends {
		tx, err := charger.Charge(common.HexToAddress(auth.Consumer), common.HexToAddress(auth.Provider), auth.Nonce, new(big.Int).SetUint64(auth.Amount), common.FromHex(auth.Signature))
		if err != nil {
			log15.Error("Failed to charge payment", "authorization", auth, "error", err)
		} else {
			log15.Info("Payment charged", "tx", "http://testnet.etherscan.io/tx/"+tx.Hex())
		}
	}
	v.lock.RUnlock()

	v.lock.Lock()
	for consumer, auth := range v.auths {
		if v.pends[consumer] == auth {
			delete(v.pends, consumer)
		}
	}
	v.lock.Unlock()
}
