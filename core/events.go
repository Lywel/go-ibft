package core

import (
	"bitbucket.org/ventureslash/go-ibft"
)

// Event is an ibft event
type Event interface{}

// Size of events chanels
const (
	EventChannelBufferSize = 256
)

// RequestEvent  is emitted for a proposal to be handled
type RequestEvent struct {
	Proposal ibft.Proposal
}

// BacklogEvent  is an internal event used to store an event for latter processing
type BacklogEvent struct {
	Message *message
}

// MessageEvent is emmitted during the IBFT consensus algo
type MessageEvent struct {
	Payload []byte
}

// JoinEvent is emmitted when a peer joins the network
type JoinEvent struct {
	Address ibft.Address
}

// StateEvent  is emmitted when a peer joins the network
type StateEvent struct {
	valSet *ibft.ValidatorSet
	view   *ibft.View
}
