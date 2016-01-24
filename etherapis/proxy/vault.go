// Contains the payment vault which accumulates the previously accepted payment
// to be able to charge them upon request.

package proxy

import (
	"sync"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"
)

// Vault is a payment aggregator collecting the individual accepted payments made
// by various clients. It can at any point return the most recent payment + proof,
// and can also charge the payments via the Ethereum network.
type Vault struct {
	vaults map[common.Address]*accountVault // Individual vaults per provider account
	lock   sync.RWMutex                     // Lock protecting the vaults
}

// NewVault create and returns an empty vault ready for accepting payments.
func NewVault() *Vault {
	return &Vault{
		vaults: make(map[common.Address]*accountVault),
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

// accountVault is the payment aggregator for a single owned account.
type accountVault struct {
	auths map[common.Address]*authorization // Individual authorizations per client
	lock  sync.RWMutex                      // Lock protecting the authorizations

}

// Store inserts a new authorized payment into the account vault for later
// redemption.
func (v *accountVault) Store(auth *authorization) {
	v.lock.Lock()
	defer v.lock.Unlock()

	consumer := common.HexToAddress(auth.Consumer)
	v.auths[consumer] = auth
}

// Fetch retrieves the last known authorization made by a particular consumer to
// this particular account vault.
func (v *accountVault) Fetch(consumer common.Address) *authorization {
	v.lock.RLock()
	defer v.lock.RUnlock()

	return v.auths[consumer]
}
