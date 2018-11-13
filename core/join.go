package core

import (
	"bitbucket.org/ventureslash/go-ibft"
)

func (c *core) handleJoin(src ibft.Address) {
	if c.isProposer() {
		c.eventsOut <- AddValidatorEvent{
			Address: src,
		}
		c.sendState(src)
	}
}

func (c *core) sendState(src ibft.Address) {
	c.eventsOut <- StateEvent{
		valSet: c.valSet,
		view:   c.currentView(),
		dest:   src,
	}
}

func (c *core) handleStateEvent(valset *ibft.ValidatorSet, view *ibft.View,
	dest ibft.Address) {
	if dest == c.address {
		c.valSet = valset
		c.valSet.AddValidator(c.address)
		c.current = newRoundState(view, nil, valset, nil)
		// TODO: start consensus
	}
}
