package crypto

import (
	"crypto/ecdsa"

	eth "github.com/ethereum/go-ethereum/crypto"
	"github.com/Lywel/ibft-go/consensus"
)

func Sign(data []byte, privkey *ecdsa.PrivateKey) ([]byte, error) {
	hashData := eth.Keccak256([]byte(data))
	return eth.Sign(hashData, privkey)
}

func PubkeyToAddress(p ecdsa.PublicKey) consensus.Address {
	ethAddress := eth.PubkeyToAddress(p)
	var a Address
	copy(a[AddressLength-len(ethAddress):], ethAddress[:])
	return a
}
