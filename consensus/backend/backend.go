package backend

import (
	"crypto/ecdsa"
	"github.com/Lywel/go-gossipnet"
	"github.com/Lywel/go-ibft/consensus"
	"github.com/Lywel/go-ibft/consensus/backend/crypto"
)

// Backend is
type Backend struct {
	privateKey *ecdsa.PrivateKey
	Address    consensus.Address
	Network    *gossipnet.Node
}

// Config is the backend configuration struct
type Config struct {
	LocalAddr   string
	RemoteAddrs []string
}

// New returns a new Backend
func New(config *Config, privateKey *ecdsa.PrivateKey) Backend {
	network := gossipnet.New(config.LocalAddr, config.RemoteAddrs)
	return Backend{
		privateKey: privateKey,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		Network:    network,
	}
}

// Sign implements Backend.Sign
func (b *Backend) Sign(data []byte) ([]byte, error) {
	return crypto.Sign(data, b.privateKey)
}
