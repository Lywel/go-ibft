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

// GetSignatureAddress returns an address from a signature
func GetSignatureAddress(data []byte, sig []byte) (ibft.Address, error) {
	// 1. Keccak data
	hashData := crypto.Keccak256([]byte(data))
	// 2. Recover public key
	pubkey, err := crypto.SigToPub(hashData, sig)
	if err != nil {
		return ibft.Address{}, err
	}
	return PubkeyToAddress(*pubkey), nil
}
