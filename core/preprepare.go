package core

import (
	"bitbucket.org/ventureslash/go-ibft"
)

func (c *core) sendPreprepare(request *ibft.Request) {
	if request.Proposal.Number().Cmp(c.current.sequence) == 0 && c.isProposer() {
		c.logger.Log("Send preprepare")
		preprepare, err := Encode(&ibft.Preprepare{
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

func (c *core) handlePreprepare(msg *message, src *ibft.Validator) error {
	c.logger.Log("Handle preprepare from", src.String())

	var preprepareRaw *PreprepareRaw

	msg.Decode(&preprepareRaw)

	proposal, err := c.backend.DecodeProposal(preprepareRaw.Proposal[0])

	if err != nil {
		return errFailedDecodePreprepare
	}

	var preprepare = &preprepare{
		View: preprepareRaw.View,
		Proposal: proposal
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
