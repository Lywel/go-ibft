package core

import (
	"github.com/Lywel/go-ibft/consensus"
)

func (c *core) sendPreprepare(request *consensus.Request) {
	if request.Proposal.Number().Cmp(c.current.sequence) == 0 && c.isProposer() {
		c.logger.Log("Send preprepare")
		preprepare, err := Encode(&consensus.Preprepare{
			View:     c.currentView(),
			Proposal: request.Proposal,
		})
		if err != nil {
			c.logger.Log("failed to encode")
			return
		}
		c.broadcast(&message{
			Type: typePreprepare,
			Msg:  preprepare,
		})
	}
}

func (c *core) handlePreprepare(msg *message, src *consensus.Validator) error {
	c.logger.Log("Handle preprepare from", src.String())
	var preprepare *consensus.Preprepare

	// Decode msg.Msg and fill preprepare
	err := msg.Decode(&preprepare)
	if err != nil {
		return errFailedDecodePreprepare
	}
	if err := c.checkMessage(typePreprepare, preprepare.View); err != nil {
		c.logger.Log("check message failed")
		return err
	}

	if !c.valSet.IsProposer(src.Address()) {
		return errNotFromProposer
	}
	if err := c.verify(preprepare.Proposal); err != nil {
		c.logger.Log("failed to verify proposal")
		// TODO
		// if it's a future block, we will handle it again after the duration
		// else sendNextRoundChange

		return err
	}
	if c.state == StateAcceptRequest {
		c.setState(StatePreprepared)
		c.current.Preprepare = preprepare
		c.sendPrepare()
	}
	return nil
}
