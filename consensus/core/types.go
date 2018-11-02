package core

import (
	"github.com/Lywel/go-ibft/consensus"
	"github.com/ethereum/go-ethereum/rlp"
)

// State is an enum representing the current state
type State uint

// State* are core states
const (
	StateAcceptRequest State = iota
	StatePreprepared
	StatePrepared
	StateCommitted
)

// Cmp compares s and y and returns:
//   -1 if s is the previous state of y
//    0 if s and y are the same state
//   +1 if s is the next state of y
func (s State) Cmp(y State) int {
	if uint64(s) < uint64(y) {
		return -1
	}
	if uint64(s) > uint64(y) {
		return 1
	}
	return 0
}

func (s State) String() string {
	if s == StateAcceptRequest {
		return "Accept request"
	} else if s == StatePreprepared {
		return "Preprepared"
	} else if s == StatePrepared {
		return "Prepared"
	} else if s == StateCommitted {
		return "Committed"
	} else {
		return "Unknown"
	}
}

const (
	typePreprepare = iota
	typePrepare
	typeCommit
	typeRoundChange
)

type message struct {
	Type      int
	Msg       []byte
	Address   consensus.Address
	Signature []byte
}

func (m *message) FromPayload(b []byte, validateFn func([]byte, []byte) (consensus.Address, error)) error {
	if err := rlp.DecodeBytes(b, &m); err != nil {
		return err
	}
	payloadNoSig, err := m.PayloadNoSig()
	if err != nil {
		return err
	}
	_, err = validateFn(payloadNoSig, m.Signature)
	return err
}

func (m *message) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

func (m *message) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&message{
		Type:      m.Type,
		Msg:       m.Msg,
		Address:   m.Address,
		Signature: []byte{},
	})
}

// Decode a message following the rlp standard
func (m *message) Decode(val interface{}) error {
	return rlp.DecodeBytes(m.Msg, val)
}

// Encode a message with the rlp standard
func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
