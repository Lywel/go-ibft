package core

import (
	"math/big"
	"sync"

	"bitbucket.org/ventureslash/go-ibft"
)

// map messages to validator address
type messageSet struct {
	view       *ibft.View
	valSet     *ibft.ValidatorSet
	messagesMu *sync.Mutex
	messages   map[ibft.Address]*message
}

// newMessageSet construct a new message set to accumulate messages for given
// sequence/view number.
func newMessageSet(valSet *ibft.ValidatorSet) *messageSet {
	msgSet := &messageSet{
		view: &ibft.View{
			Round:    new(big.Int),
			Sequence: new(big.Int),
		},
		messagesMu: new(sync.Mutex),
		messages:   make(map[ibft.Address]*message),
		valSet:     valSet,
	}
	return msgSet
}

// View returns the current view
func (ms *messageSet) View() *ibft.View {
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
func (ms *messageSet) Get(addr ibft.Address) *message {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	return ms.messages[addr]
}

func (ms *messageSet) Values() (result []*message) {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()

	for _, v := range ms.messages {
		result = append(result, v)
	}

	return result
}

func (ms *messageSet) verifyAddr(msg *message) bool {
	if _, v := ms.valSet.GetByAddress(msg.Address); v != nil {
		return true
	}
	return false
}
