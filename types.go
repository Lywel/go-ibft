package ibft

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
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

// Address of client
type Address [AddressLength]byte

// PubkeyToAddress return the address corresponding to ecdsa public key
func PubkeyToAddress(p ecdsa.PublicKey) Address {
	ethAddress := crypto.PubkeyToAddress(p)
	var a Address
	copy(a[AddressLength-len(ethAddress):], ethAddress[:])
	return a
}

// GetBytes returns the address as bytes
func (a Address) GetBytes() [AddressLength]byte {
	return [AddressLength]byte(a)
}

// FromBytes poppulates the address from bytes in argument
func (a *Address) FromBytes(data []byte) {
	copy(a[AddressLength-len(data):], data)
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


type PreprepareRaw struct {
	View *View
	Proposal []interface{} `rlp:"tail"`
}


// Preprepare include the proposal and the current view
type Preprepare struct {
	View     *View
	Proposal Proposal
}

func (b *Preprepare) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{b.View, b.Proposal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *PreprepareRaw) DecodeRLP(s *rlp.Stream) error {

	var preprepareRaw *PreprepareRaw
	if err := s.Decode(&preprepareRaw); err != nil {
		return err
	}
	// TODO call decodeProposal
	/*var preprepare struct {
		View     *View
		Proposal *types.Block
	}

	if err := s.Decode(&preprepare); err != nil {
		return err
	}
	b.View, b.Proposal = preprepare.View, preprepare.Proposal*/

	return nil
}

// Subject include the proposal digest and the current view
type Subject struct {
	View   *View
	Digest []byte
}
