package core

import (
	"bitbucket.org/ventureslash/go-ibft"
)

func (c *core) handleJoin(src ibft.Address) {
	if c.isProposer() {
		c.logger.Info(c.address, ": Handle join from ", src)
		c.eventsOut <- AddValidatorEvent{
			Address: src,
		}
		c.sendState(src)
	}
}

func (c *core) sendState(src ibft.Address) {
	c.eventsOut <- StateEvent{
		ValSet: c.valSet,
		View:   c.currentView(),
		Dest:   src,
	}
}

func (c *core) handleStateEvent(valset *ibft.ValidatorSet, view *ibft.View,
	dest ibft.Address) {
	// TODO: fix and add security
	if dest == c.address {
		c.logger.Info(c.address, ": Received state")
		c.valSet = valset
		c.valSet.AddValidator(c.address)
		c.initTimeouts()
		c.current = newRoundState(view, nil, c.valSet, nil)
		c.logger.Info(c.address, ": view ", view)
		// c.setState(StateAcceptRequest)
	}
}
