package consensus

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	"github.com/Lywel/go-gossipnet"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	AddressLength = 20
)

var (
	Big0 = big.NewInt(0)
	Big1 = big.NewInt(1)
)

// Engine can be started and stoped
type Engine interface {
	Start()
	Stop()
}

type Backend interface {
	Start()
	Stop()
	Network() *gossipnet.Node
	Address() Address
	Sign(data []byte) ([]byte, error)
	AddValidator(addr Address) bool
}

// Address of client
type Address [AddressLength]byte

// PubkeyToAddress return the address corresponding to ecdsa public key
func PubkeyToAddress(p ecdsa.PublicKey) Address {
	ethAddress := crypto.PubkeyToAddress(p)
	var a Address
	copy(a[AddressLength-len(ethAddress):], ethAddress[:])
	return a
}

func (a Address) String() string {
	var bytes [AddressLength]byte
	copy(bytes[:], a[:])
	return hex.EncodeToString(bytes[:])
}

// Proposal interface to be used during the consensus
type Proposal interface {
	Number() *big.Int
	Hash() []byte
	String() string
}

// Request is the original request of the client
// It is sent to handleRequest
type Request struct {
	Proposal Proposal
}

// View include round number and sequence number
type View struct {
	Round    *big.Int
	Sequence *big.Int
}

// Cmp compare v and y
// Priority: Sequence > Round
func (v *View) Cmp(y *View) int {
	if v.Sequence.Cmp(y.Sequence) != 0 {
		return v.Sequence.Cmp(y.Sequence)
	}
	if v.Round.Cmp(y.Round) != 0 {
		return v.Round.Cmp(y.Round)
	}
	return 0
}

// Preprepare include the proposal and the current view
type Preprepare struct {
	View     *View
	Proposal Proposal
}

// Subject include the proposal digest and the current view
type Subject struct {
	View   *View
	Digest []byte
}
