package core

import (
	"log"

	"bitbucket.org/ventureslash/go-ibft"
	"github.com/ethereum/go-ethereum/rlp"
)

func (c *core) handleEvents() {
	c.wg.Add(1)
	defer c.wg.Done()
	for event := range c.eventsIn {
		switch ev := event.(type) {
		case MessageEvent:
			c.handleMsg(ev.Payload)
		case RequestEvent:
			r := &ibft.Request{
				Proposal: ev.Proposal,
			}
			err := c.handleRequest(r)
			if err == errFutureMessage {
				c.storeRequest(r)
			} else if err != nil {
				c.logger.Warning(c.address, ": Handle event request ", "err ", err)
			}
		case EncodedRequestEvent:
			var encodedProposal *ibft.EncodedProposal
			err := rlp.DecodeBytes(ev.Proposal, &encodedProposal)
			if err != nil {
				log.Print(err)
				continue
			}
			decodedProposal, err := c.proposalManager.DecodeProposal(encodedProposal)
			if err != nil {
				log.Print(err)
				continue
			}
			c.eventsIn <- RequestEvent{Proposal: decodedProposal}
		case BacklogEvent:
			_, src := c.valSet.GetByAddress(ev.Message.Address)
			c.handleCheckedMsg(ev.Message, src)
		case JoinEvent:
			if ev.Address != c.address {
				c.logger.Info(c.address, ": New peer: ", ev.Address)
				c.networkMap[ev.Address] = ev.NetworkAddr
				c.handleJoin(ev.Address)
			}

		case StateEvent:
			c.logger.Info(c.address, ": Received stateEvent ", "view ", ev.View)
			c.handleStateEvent(ev.ValSet, ev.View, ev.Dest)
		case AddValidatorEvent:
			if ev.Address != c.address {
				res := c.valSet.AddValidator(ev.Address)
				if res {
					c.logger.Info(c.address, ": Adding validator ", ev.Address)
					c.setValidatorTimeout(ev.Address)
				}
			}
		}

	}
	c.logger.Info(c.address, ": End of handle events")
}

// decodes message, checks and calls handleCheckedMsg
func (c *core) handleMsg(payload []byte) error {
	msg := new(message)
	if err := msg.FromPayload(payload, c.ValidateFn); err != nil {
		c.logger.Warning(c.address, ": Failed to decode message from payload ", "err ", err)
		return err
	}
	_, src := c.valSet.GetByAddress(msg.Address)
	if src == nil {
		c.logger.Warning(c.address, ": invalid address in message ", "msg ", msg)
		return errUnauthorized
	}
	return c.handleCheckedMsg(msg, src)
}

// handles the message, and stores it in backlog if needed
func (c *core) handleCheckedMsg(msg *message, src *ibft.Validator) error {
	c.setValidatorTimeout(msg.Address)

	// add message in backlog if it is a future message
	testBacklog := func(err error) error {
		if err == errFutureMessage {
			c.storeBacklog(msg, src)
		}
		return err
	}

	switch msg.Type {
	case typePreprepare:
		return testBacklog(c.handlePreprepare(msg, src))
	case typePrepare:
		return testBacklog(c.handlePrepare(msg, src))
	case typeCommit:
		return testBacklog(c.handleCommit(msg, src))
	case typeRoundChange:
		return testBacklog(c.handleRoundChange(msg, src))
	default:
		return errInvalidMessage
	}

}

func (c *core) handleTimeoutMsg() {
	if !c.waitingForRoundChange {
		maxRound := c.roundChangeSet.MaxRound(c.valSet.F() + 1)
		if maxRound != nil && maxRound.Cmp(c.current.round) > 0 {
			c.sendRoundChange(maxRound)
			return
		}
	}
	c.sendNextRoundChange()

}
