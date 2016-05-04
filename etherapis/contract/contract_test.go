package contract

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestDeploymentAndIntegration(t *testing.T) {
	// generate new key
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	// init new simulated backend
	sim := backends.NewSimulatedBackend(core.GenesisAccount{Address: auth.From, Balance: big.NewInt(10000000000)})

	// deploy the contract
	_, _, api, err := DeployEtherAPIs(auth, sim)

	// use a session based approach so that we do not need
	// to repass these settings all the time.
	session := &EtherAPIsSession{
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
	_, err = session.AddService("etherapis", "https://etherapis.io", big.NewInt(0), big.NewInt(10), big.NewInt(432000))
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
