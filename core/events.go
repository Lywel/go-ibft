package core

import (
	"bitbucket.org/ventureslash/go-ibft"
)

// RequestEvent is emitted for a proposal to be handled
type RequestEvent struct {
	Proposal ibft.Proposal
}

// BacklogEvent is an internal event used to store an event for latter processing
type BacklogEvent struct {
	Message *message
}

// StateEvent is emmitted when a peer joins the network
type StateEvent struct {
	valSet ibft.ValidatorSet
	state  State
}
