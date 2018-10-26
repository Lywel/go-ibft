package backend

import (
	"crypto/ecdsa"
)

// Backend is
type Backend struct {
	privateKey *ecdsa.PrivateKey
	address    Address
}
