package backend

import (
	"crypto/ecdsa"
	"github.com/Lywel/go-ibft/consensus/"
	"github.com/Lywel/go-ibft/consensus/backend/crypto"
)

// Backend is
type Backend struct {
	privateKey *ecdsa.PrivateKey
	address    consensus.Address
}

func (b *Backend) Sign(data []byte) ([]byte, error) {
	return crypto.Sign(data, b.privateKey)
}
