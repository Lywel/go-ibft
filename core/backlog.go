package core

import (
	"bitbucket.org/ventureslash/go-ibft"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

var (
	typePriority = map[uint64]int{
		typePreprepare: 1,
		typeCommit:     2,
		typePrepare:    3,
	}
)

func (c *core) storeBacklog(msg *message, src *ibft.Validator) {
	if src.Address() == c.address {
		return
	}
	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	backlog := c.backlogs[src]
	if backlog == nil {
		backlog = prque.New()
	}
	switch msg.Type {
	case typePreprepare:

		var encodedPreprepare *ibft.EncodedPreprepare
		msg.Decode(&encodedPreprepare)

		proposal, err := c.proposalManager.DecodeProposal(encodedPreprepare.Prop)
		if err != nil {
			return
		}

		var preprepare = &ibft.Preprepare{
			View:     encodedPreprepare.View,
			Proposal: proposal,
		}
		if err == nil {
			backlog.Push(msg, toPriority(typePrepare, preprepare.View))
		}
	default:
		var subject *ibft.Subject
		err := msg.Decode(&subject)
		if err == nil {
			backlog.Push(msg, toPriority(typePrepare, subject.View))
		}
	}
	c.backlogs[src] = backlog
}

func (c *core) processBacklogs() {
	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	// Item on each validator backlog
	for _, backlog := range c.backlogs {
		if backlog == nil {
			continue
		}
		isFuture := false

		// Iter on a backlog
		for !backlog.Empty() && !isFuture {
			m, prio := backlog.Pop()
			msg, ok := m.(*message)
			if !ok {
				c.logger.Log("failed to cast message from backlog")
				break
			}
			var view *ibft.View
			switch msg.Type {
			case typePrepare:
				var preprepare *ibft.Preprepare
				err := msg.Decode(&preprepare)
				if err == nil {
					view = preprepare.View
				}
			default:
				var subject ibft.Subject
				err := msg.Decode(&subject)
				if err == nil {
					view = subject.View
				}
			}
			if err := c.checkMessage(msg.Type, view); err != nil {
				if err == errFutureMessage {
					c.logger.Log("stop processing future backlog", "msg", msg)
					backlog.Push(msg, prio)
					isFuture = true
					break
				}
				c.logger.Log("stop processing invalid backlog", "msg", msg)
				continue
			}
			c.eventsIn <- BacklogEvent{
				Message: msg,
			}
		}

	}
}

func toPriority(msgType uint64, view *ibft.View) float32 {
	if msgType == typeRoundChange {
		// For msgRoundChange, set the message priority based on its sequence
		return -float32(view.Sequence.Uint64() * 1000)
	}
	// 10 * Round limits the range of message code is from 0 to 9
	// 1000 * Sequence limits the range of round is from 0 to 99
	return -float32(view.Sequence.Uint64()*1000 + view.Round.Uint64()*10 + uint64(typePriority[msgType]))
}
