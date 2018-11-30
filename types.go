package ibft

import (
	"crypto/ecdsa"
	"encoding/hex"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	AddressLength    = 20
	ValidatorTimeout = 20 * time.Second
	RequestTimeout   = 20 * time.Second
)

var (
	Big0 = big.NewInt(0)
	Big1 = big.NewInt(1)
)

// Engine can be started and stoped
type Core interface {
	Start()
	Stop()
	NetworkMap() map[Address]string
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
	return hex.EncodeToString(bytes[:6])
}

// Proposal interface to be used during the consensus
type Proposal interface {
	Number() *big.Int
	Hash() []byte
	String() string
	EncodeRLP(w io.Writer) error
	DecodeRLP(s *rlp.Stream) error
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

type EncodedProposal struct {
	Type uint
	Prop []byte
}

type EncodedPreprepare struct {
	View *View
	Prop *EncodedProposal
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
func (b *EncodedPreprepare) DecodeRLP(s *rlp.Stream) error {
	var encodedPreprepare struct {
		View *View
		Prop *EncodedProposal
	}
	if err := s.Decode(&encodedPreprepare); err != nil {
		return err
	}
	b.View, b.Prop = encodedPreprepare.View, encodedPreprepare.Prop
	return nil
}

// Subject include the proposal digest and the current view
type Subject struct {
	View   *View
	Digest []byte
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (b *Subject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{b.View, b.Digest})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *Subject) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		View   *View
		Digest []byte
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}
	b.View, b.Digest = subject.View, subject.Digest
	return nil
}

type ProposalManager interface {
	DecodeProposal(prop *EncodedProposal) (Proposal, error)
	Verify(proposal Proposal) error
	Commit(proposal Proposal) error
}
