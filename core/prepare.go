package core

import (
	"reflect"

	"bitbucket.org/ventureslash/go-ibft"
)

func (c *core) sendPrepare() {
	c.logger.Info(c.address, ": Send prepare")
	subject := c.current.Subject()
	prepare, err := Encode(subject)
	if err != nil {
		c.logger.Warning(c.address, ": Encode subject failed")
		return
	}
	c.broadcast(&message{
		Type: typePrepare,
		Msg:  prepare,
	})
}

func (c *core) handlePrepare(msg *message, src *ibft.Validator) error {
	c.logger.Info(c.address, ": Handle prepare from ", src.String())
	var prepare *ibft.Subject
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
		c.logger.Warning(c.address, ": Failed to add PREPARE message ", "msg ", msg, " err ", err)
	}
	if c.current.GetPrepareOrCommitSize() > 2*c.valSet.F() && c.state.Cmp(StatePrepared) < 0 {
		c.logger.Info("got ", c.current.Prepares.Size(), " signatures, needed more than ", 2*c.valSet.F())

		c.setState(StatePrepared)
		c.sendCommit()
	} else {
		c.logger.Info("got ", c.current.Prepares.Size(), " signatures, need more than ", 2*c.valSet.F())

	}

	return nil
}

func (c *core) verifyPrepare(prepare *ibft.Subject) error {
	subject := c.current.Subject()
	if !reflect.DeepEqual(prepare, subject) {
		c.logger.Info(c.address, ": Subjects do not match: expected ", subject, " got ", prepare)
		return errSubjectsDoNotMatch
	}
	return nil
}
