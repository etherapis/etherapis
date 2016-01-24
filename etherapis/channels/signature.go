package channels

import "github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/github.com/ethereum/go-ethereum/common"

type signature struct {
	v    uint64
	r, s common.Hash
}

func bytesToSignature(sig []byte) signature {
	var signature signature
	signature.v = uint64(sig[64] + 27)
	signature.r = common.BytesToHash(sig[:32])
	signature.s = common.BytesToHash(sig[32:64])
	return signature
}
