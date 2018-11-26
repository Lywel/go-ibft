package backend

import (
	"crypto/ecdsa"

	"bitbucket.org/ventureslash/go-gossipnet"
	"bitbucket.org/ventureslash/go-ibft"
	"bitbucket.org/ventureslash/go-ibft/backend/events"
	"bitbucket.org/ventureslash/go-ibft/core"
	"bitbucket.org/ventureslash/go-ibft/crypto"
)

// Backend initializes the core, holds the keys and currenncy logic
type Backend struct {
	privateKey      *ecdsa.PrivateKey
	address         ibft.Address
	network         *gossipnet.Node
	core            ibft.Engine
	coreRunning     bool
	ibftEventsIn    chan core.Event
	ibftEventsOut   chan core.Event
	manager         events.Manager
	proposalManager ibft.ProposalManager
	stillConnecting int
}

// Config is the backend configuration struct
type Config struct {
	LocalAddr   string
	RemoteAddrs []string
}

// New returns a new Backend
func New(config *Config,
	privateKey *ecdsa.PrivateKey,
	proposalManager ibft.ProposalManager,
	eventProxy func(chan core.Event, chan core.Event) (chan core.Event, chan core.Event)) *Backend {

	network := gossipnet.New(config.LocalAddr, config.RemoteAddrs)
	in := make(chan core.Event, 256)
	out := make(chan core.Event, 256)
	pin, pout := eventProxy(in, out)

	backend := &Backend{
		privateKey:      privateKey,
		address:         crypto.PubkeyToAddress(privateKey.PublicKey),
		network:         network,
		ibftEventsIn:    in,
		ibftEventsOut:   out,
		manager:         events.New(network, pin, pout, len(config.RemoteAddrs)),
		proposalManager: proposalManager,
		stillConnecting: len(config.RemoteAddrs),
	}

	backend.core = core.New(backend, proposalManager)
	return backend
}

func (b *Backend) Network() map[ibft.Address]string {
	return b.core.NetworkMap()
}

// PrivateKey returns the private key
func (b *Backend) PrivateKey() *ecdsa.PrivateKey {
	return b.privateKey
}

// Start implements Engine.Start
func (b *Backend) Start() {
	if b.coreRunning {
		return
	}
	b.manager.Start(b.address)
	b.network.Start()
	b.core.Start()
	b.coreRunning = true
}

// Stop implements Engine.Stop
func (b *Backend) Stop() {
	if !b.coreRunning {
		return
	}
	b.network.Stop()
	b.core.Stop()
	b.coreRunning = false
}

// EventsInChan returns a channel receiving network events
func (b *Backend) EventsInChan() chan core.Event {
	return b.ibftEventsIn
}

// EventsOutChan returns a channel used to emit events to the network
func (b *Backend) EventsOutChan() chan core.Event {
	return b.ibftEventsOut
}
