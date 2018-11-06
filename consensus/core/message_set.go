package core

import (
	"math/big"
	"sync"

	"bitbucket.org/ventureslash/go-ibft/consensus"
)

// map messages to validator address
type messageSet struct {
	view       *consensus.View
	valSet     *consensus.ValidatorSet
	messagesMu *sync.Mutex
	messages   map[consensus.Address]*message
}

// newMessageSet construct a new message set to accumulate messages for given
// sequence/view number.
func newMessageSet(valSet *consensus.ValidatorSet) *messageSet {
	msgSet := &messageSet{
		view: &consensus.View{
			Round:    new(big.Int),
			Sequence: new(big.Int),
		},
		messagesMu: new(sync.Mutex),
		messages:   make(map[consensus.Address]*message),
		valSet:     valSet,
	}
	return msgSet
}

// View returns the current view
func (ms *messageSet) View() *consensus.View {
	return ms.view
}

// Add a message to the set if the messages is from one of the validator
func (ms *messageSet) Add(msg *message) error {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	if !ms.verifyAddr(msg) {
		return errNotAValidatorAddress
	}
	ms.messages[msg.Address] = msg
	return nil
}

// Size returns the size of the messages map
func (ms *messageSet) Size() int {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	return len(ms.messages)
}

// Get the message from an address
func (ms *messageSet) Get(addr consensus.Address) *message {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	return ms.messages[addr]
}

func (ms *messageSet) verifyAddr(msg *message) bool {
	if _, v := ms.valSet.GetByAddress(msg.Address); v != nil {
		return true
	}
	return false
}
