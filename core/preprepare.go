package core

import (
	"bitbucket.org/ventureslash/go-ibft"
)

func (c *core) sendPreprepare(request *ibft.Request) {
	if request.Proposal.Number().Cmp(c.current.sequence) == 0 && c.isProposer() {
		c.logger.Info(c.address, ": Send preprepare")
		preprepare, err := Encode(&ibft.Preprepare{
			View:     c.currentView(),
			Proposal: request.Proposal,
		})
		if err != nil {
			c.logger.Warning(c.address, ": Failed to encode")
			return
		}

		c.broadcast(&message{
			Type: typePreprepare,
			Msg:  preprepare,
		})
	}
}

func (c *core) handlePreprepare(msg *message, src *ibft.Validator) error {
	c.logger.Info(c.address, ": Handle preprepare from ", src.String())

	var encodedPreprepare *ibft.EncodedPreprepare
	err := msg.Decode(&encodedPreprepare)
	if err != nil {
		c.logger.Warning("failed to decode msg ", "err ", err)
		return err
	}

	proposal, err := c.proposalManager.DecodeProposal(encodedPreprepare.Prop)
	if err != nil {
		c.logger.Warning("failed to decode proposal ", "err ", err)
		return err
	}

	var preprepare = &ibft.Preprepare{
		View:     encodedPreprepare.View,
		Proposal: proposal,
	}

	if err := c.checkMessage(typePreprepare, preprepare.View); err != nil {
		c.logger.Warning(c.address, ": Check message failed err ", err)
		return err
	}

	if !c.valSet.IsProposer(src.Address()) {
		return errNotFromProposer
	}
	if err := c.verify(preprepare.Proposal); err != nil {
		c.logger.Warning(c.address, ": Failed to verify proposal ", "err ", err)

		return err
	}
	if c.state == StateAcceptRequest {
		c.setState(StatePreprepared)
		c.current.Preprepare = preprepare
		c.sendPrepare()
	}
	return nil
}
