package crypto

import (
	"crypto/ecdsa"

	"bitbucket.org/ventureslash/go-ibft"
	"github.com/ethereum/go-ethereum/crypto"
)

// Sign returns the signature of data from from privateKey
func Sign(data []byte, privkey *ecdsa.PrivateKey) ([]byte, error) {
	hashData := crypto.Keccak256([]byte(data))
	return crypto.Sign(hashData, privkey)
}

// PubkeyToAddress returns an Address from a ecdsa.PublicKey
func PubkeyToAddress(p ecdsa.PublicKey) ibft.Address {
	ethAddress := crypto.PubkeyToAddress(p)
	var a ibft.Address
	copy(a[ibft.AddressLength-len(ethAddress):], ethAddress[:])
	return a
}