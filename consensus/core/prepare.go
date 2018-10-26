package core

import (
	"reflect"

	"github.com/Lywel/ibft-go/consensus"
)

func (c *core) sendPrepare() {
	c.logger.Log("send prepare")
	subject := c.current.Subject()
	prepare, err := Encode(subject)
	if err != nil {
		c.logger.Log("Encode subject failed")
		return
	}
	c.broadcast(&message{
		Type: typePrepare,
		Msg:  prepare,
	})
}

func (c *core) handlePrepare(msg *message, src *consensus.Validator) error {
	c.logger.Log("handle prepare from", src.String())
	var prepare *consensus.Subject
	err := msg.Decode(&prepare)
	if err != nil {
		return errFailedDecodePrepare
	}

	if err := c.checkMessage(typePrepare, prepare.View); err != nil {
		return err
	}

	if err := c.verifyPrepare(prepare); err != nil {
		return err
	}
	if err := c.current.Prepares.Add(msg); err != nil {
		c.logger.Log("Failed to add PREPARE message", "msg", msg, "err", err)
	}
	if c.current.Prepares.Size() > 2*c.valSet.F() && c.state.Cmp(StatePrepared) < 0 {
		c.setState(StatePrepared)
		c.sendCommit()
	}

	return nil
}

func (c *core) verifyPrepare(prepare *consensus.Subject) error {
	subject := c.current.Subject()
	if !reflect.DeepEqual(prepare, subject) {
		c.logger.Log("subjects do not match: expected", subject, "got", prepare)
		return errSubjectsDoNotMatch
	}
	return nil
}
