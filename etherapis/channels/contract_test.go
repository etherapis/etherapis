package channels

import "testing"

func TestChannels(t *testing.T) {
	/*
		var mux event.TypeMux

		db, _ := ethdb.NewMemDatabase()
		statedb, _ := state.New(common.Hash{}, db)
		stateFn := func() *state.StateDB {
			return statedb
		}

		key1, _ := crypto.GenerateKey()
		key2, _ := crypto.GenerateKey()
		to := crypto.PubkeyToAddress(key2.PublicKey)

		contract, err := Fetch(db, &mux, stateFn)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("methods:", contract.abi.Methods)
		fmt.Println("events:", contract.abi.Events)
		fmt.Println("\n\n")

		tx, err := contract.NewChannel(key1, to, new(big.Int).Mul(big.NewInt(10), common.Ether), func(c *Channel) {
			fmt.Println("new  channel created", c)
		})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(tx)
	*/
}
