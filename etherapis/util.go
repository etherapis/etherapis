package main

import (
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"
	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/crypto"
)

// sha3 returns the canonical sha3 of the 32byte (padded) input
func sha3(in ...[]byte) []byte {
	out := make([]byte, len(in)*32)
	for i, input := range in {
		copy(out[i*32:i*32+32], common.LeftPadBytes(input, 32))
	}
	return crypto.Sha3(out)
}
