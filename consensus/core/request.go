package core

import (
	"github.com/Lywel/ibft-go/consensus"
)

func (c *core) handleRequest(request *consensus.Request) error {
	if err := c.checkRequest(request); err != nil {
		if err == errInvalidMessage {
			c.logger.Log("invalid message")
			return err
		}
		c.logger.Log("unexpected request", "err", err, "number",
			request.Proposal.Number())
		return err
	}
	c.logger.Log("HandleRequest number", request.Proposal.Number(), "hash",
		request.Proposal.Hash())
	c.current.pendingRequest = request
	if c.state == StateAcceptRequest {
		c.sendPreprepare(request)
	}
	return nil
}

func (c *core) checkRequest(request *consensus.Request) error {
	if request == nil || request.Proposal == nil {
		return errInvalidMessage
	}

	if c := c.current.sequence.Cmp(request.Proposal.Number()); c > 0 {
		return errOldMessage
	} else if c < 0 {
		return errFutureMessage
	}
	return nil
}

func (c *core) storeRequest(request *consensus.Request) {
	c.logger.Log("storing future request", "number", request.Proposal.Number(),
		"Hash", request.Proposal.Hash())
	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()
	c.pendingRequests.Push(request, -float32(request.Proposal.Number().Int64()))
}

func (c *core) processPendingRequests() {
	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	for !c.pendingRequests.Empty() {
		m, prio := c.pendingRequests.Pop()
		r, ok := m.(*consensus.Request)

		if !ok {
			c.logger.Log("request malformed")
			continue
		}
		if err := c.checkRequest(r); err != nil {
			if err == errFutureMessage {
				c.logger.Log("future message", "number", r.Proposal.Number(),
					"Hash", r.Proposal.Hash())

				c.pendingRequests.Push(m, prio)
			}
			c.logger.Log("skipping pending request", "number", r.Proposal.Number(),
				"Hash", r.Proposal.Hash())
			continue
		}
		c.logger.Log("processing pending request", "number", r.Proposal.Number(),
			"Hash", r.Proposal.Hash())

		// TODO send request event
	}
}
