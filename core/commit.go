package core

import (
	"reflect"

	"bitbucket.org/ventureslash/go-ibft"
)

func (c *core) sendCommit() {
	c.logger.Info(c.address, ": Send commit")
	subject := c.current.Subject()
	commit, err := Encode(subject)
	if err != nil {
		c.logger.Warning(c.address, ": Encode subject failed")
		return
	}
	c.broadcast(&message{
		Type: typeCommit,
		Msg:  commit,
	})
}

func (c *core) handleCommit(msg *message, src *ibft.Validator) error {
	c.logger.Info(c.address, ": Handle commit from ", src)
	var commit *ibft.Subject
	err := msg.Decode(&commit)
	if err != nil {
		return errFailedDecodeCommit
	}
	if err := c.checkMessage(typeCommit, commit.View); err != nil {
		return err
	}
	if err := c.verifyCommit(commit); err != nil {
		return err
	}
	if err := c.current.Commits.Add(msg); err != nil {
		c.logger.Warning(c.address, ": Failed to add COMMIT message ", "msg ", msg, " err ", err)
	}
	if c.current.Commits.Size() > 2*c.valSet.F() && c.state.Cmp(StateCommitted) < 0 {
		c.logger.Info("got ", c.current.Commits.Size(), " signatures, needed more than ", 2*c.valSet.F())
		c.commit()
	}
	return nil
}

func (c *core) verifyCommit(commit *ibft.Subject) error {
	subject := c.current.Subject()
	if !reflect.DeepEqual(commit, subject) {
		c.logger.Warning(c.address, ": Subjects do not match: expected ", subject, " got ", commit)
		return errSubjectsDoNotMatch
	}
	return nil
}
