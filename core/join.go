package core

import (
	"bitbucket.org/ventureslash/go-ibft"
)

func (c *core) handleJoin(src ibft.Address) {
	if c.isProposer() {
		c.logger.Log("handle join from", src)
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
	// TODO: fix
	if dest == c.address {
		c.logger.Log("received state")
		c.valSet = valset
		c.valSet.AddValidator(c.address)
		c.current = newRoundState(view, nil, valset, nil)
		c.logger.Log("view", view)
		// TODO: start consensus
	}
}
