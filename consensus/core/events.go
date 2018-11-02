package core

import (
	"github.com/Lywel/go-ibft/consensus"
)

// RequestEvent is emitted for a proposal to be handled
type RequestEvent struct {
	Proposal consensus.Proposal
}

// BacklogEvent is an internal event used to store an event for latter processing
type BacklogEvent struct {
	Message *message
}

// StateEvent is emmitted when a peer joins the network
type StateEvent struct {
	valSet consensus.ValidatorSet
	state  State
}
