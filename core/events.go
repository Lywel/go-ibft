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

type EncodedRequestEvent struct {
	Proposal []byte
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
	Address     ibft.Address
	NetworkAddr string
}

type AddValidatorEvent struct {
	Address ibft.Address
}

type RemoveValidatorEvent struct {
	Address ibft.Address
}

type ValidatorSetEvent struct {
	ValSet *ibft.ValidatorSet
	Dest   ibft.Address
}

type CustomEvent struct {
	Type uint
	Msg  []byte
}

type TimeOutEvent struct{}
