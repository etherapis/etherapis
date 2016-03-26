package contract

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestDeploymentAndIntegration(t *testing.T) {
	// generate new key
	key := crypto.NewKey(rand.Reader)
	// init new simulated backend
	sim := backends.NewSimulatedBackend(core.GenesisAccount{key.Address, big.NewInt(10000000000)})

	// created authenticator
	auth := bind.NewKeyedTransactor(key)

	// deploy the contract
	_, _, api, err := DeployEtherApis(auth, sim)

	// use a session based approach so that we do not need
	// to repass these settings all the time.
	session := &EtherApisSession{
		Contract: api,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: big.NewInt(400000),
		},
	}

	// add a new service
	_, err = session.AddService("etherapis", "https://etherapis.io", big.NewInt(10), big.NewInt(432000))
	if err != nil {
		t.Error(err)
	}

	// retrieve the service
	service, err := session.GetService(big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	if service.Enabled {
		t.Error("expected service to be disabled")
	}

	// enable the service
	_, err = session.EnableService(big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	service, err = session.GetService(big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	if !service.Enabled {
		t.Error("expected service to be enabled")
	}

	// flag deletion
	_, err = session.DeleteService(big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}

	service, err = session.GetService(big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	if !service.Deleted {
		t.Error("expected service to be deleted")
	}
}
