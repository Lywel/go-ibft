package core

import (
	"bitbucket.org/ventureslash/go-ibft"
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
			err := c.checkRequest(r)
			if err == errFutureMessage {
				c.storeRequest(r)
			} else if err != nil {
				c.logger.Log("handle event request", "err", err)
			}
		case BacklogEvent:
			_, src := c.valSet.GetByAddress(ev.Message.Address)
			c.handleCheckedMsg(ev.Message, src)
		case JoinEvent:
			if ev.Address != c.address {
				c.logger.Log("New peer:", ev.Address)
				c.networkMap[ev.Address] = ev.NetworkAddr
				c.handleJoin(ev.Address)
			}

		case StateEvent:
			c.logger.Log("received stateEvent", "view", ev.View)
			c.handleStateEvent(ev.ValSet, ev.View, ev.Dest)
		case AddValidatorEvent:
			if ev.Address != c.address {
				res := c.valSet.AddValidator(ev.Address)
				if res {
					c.logger.Log("adding validator", ev.Address)
					c.setValidatorTimeout(ev.Address)
				}
			}
		}

	}
	c.logger.Log("End of handle events")
}

// decodes message, checks and calls handleCheckedMsg
func (c *core) handleMsg(payload []byte) error {
	msg := new(message)
	if err := msg.FromPayload(payload, c.ValidateFn); err != nil {
		c.logger.Log("failed to decode message from payload", "err", err)
		return err
	}
	_, src := c.valSet.GetByAddress(msg.Address)
	if src == nil {
		c.logger.Log("invalid address in message", "msg", msg)
		return errUnauthorized
	}
	return c.handleCheckedMsg(msg, src)
}

// handles the message, and stores it in backlog if needed
func (c *core) handleCheckedMsg(msg *message, src *ibft.Validator) error {

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
		c.logger.Log("round change event received, ignore")
		// TODO call handleroundchange
		return nil
	default:
		return errInvalidMessage
	}

}
