package crypto

import (
	"crypto/ecdsa"
	"errors"
	eth "github.com/ethereum/go-ethereum/crypto"
	"reflect"
)

var (
	errInvalidSignature = errors.New("invalid signature")
)

func Sign(data []byte, privkey *ecdsa.PrivateKey) ([]byte, error) {
	hashData := eth.Keccak256([]byte(data))
	return eth.Sign(hashData, privkey)
}

func CheckSignature(data, sig []byte, pubkey *ecdsa.PublicKey) error {
	hashData := eth.Keccak256([]byte(data))
	signer, err := eth.SigToPub(hashData, sig)
	if err != nil {
		return err
	}

	// Compare derived addresses
	if !reflect.DeepEqual(signer, pubkey) {
		return errInvalidSignature
	}
	return nil
}

const (
	AddressLength = 20
)

type Address [AddressLength]byte

func PubkeyToAddress(p ecdsa.PublicKey) Address {
	ethAddress := eth.PubkeyToAddress(p)
	var a Address
	copy(a[AddressLength-len(ethAddress):], ethAddress[:])
	return a
}
