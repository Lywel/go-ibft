package ibft

import (
	"crypto/ecdsa"
	"encoding/hex"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 20
	// ValidatorTimeout is ??
	ValidatorTimeout = 120 * time.Second
	// RequestTimeout is ??
	RequestTimeout = 40 * time.Second
)

const (
	TypeJoinEvent = iota
	TypeValidatorSetEvent
	TypeRemoveValidatorEvent
	TypeCustomEvents
)

var (
	Big0 = big.NewInt(0)
	Big1 = big.NewInt(1)
)

// Core can be started and stoped
type Core interface {
	Start(valSet *ValidatorSet, view *View)
	Stop()
	NetworkMap() map[Address]string
}

// Hash is the common hash
type Hash [HashLength]byte

// Bytes returns a hash as bytes
func (h Hash) Bytes() []byte {
	return h[:]
}

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// RlpHash return the keccak256 hash of the rlp encoding of an interface
func RlpHash(x interface{}) Hash {
	var h common.Hash
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return BytesToHash(h.Bytes())
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
	Hash() Hash
	String() string
	EncodeRLP(w io.Writer) error
	DecodeRLP(s *rlp.Stream) error
	ExportAsRLPEncodedProposal() ([]byte, error)
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

// EncodeRLP implements RLPEncoder
func (b *Preprepare) EncodeRLP(w io.Writer) error {
	encodedProposal, err := b.Proposal.ExportAsRLPEncodedProposal()
	if err != nil {
		return err
	}
	return rlp.Encode(w, []interface{}{b.View, encodedProposal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *EncodedPreprepare) DecodeRLP(s *rlp.Stream) error {
	var encodedPreprepare struct {
		View *View
		Prop []byte
	}
	if err := s.Decode(&encodedPreprepare); err != nil {
		return err
	}
	var ep *EncodedProposal
	rlp.DecodeBytes(encodedPreprepare.Prop, &ep)
	b.View, b.Prop = encodedPreprepare.View, ep
	return nil
}

// Subject include the proposal digest and the current view
type Subject struct {
	View   *View
	Digest Hash
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (b *Subject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{b.View, b.Digest})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *Subject) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		View   *View
		Digest Hash
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}
	b.View, b.Digest = subject.View, subject.Digest
	return nil
}

// ProposalManager is able to decode verify and commit a Proposal
type ProposalManager interface {
	DecodeProposal(prop *EncodedProposal) (Proposal, error)
	Verify(proposal Proposal) error
	Commit(proposal Proposal) error
}
