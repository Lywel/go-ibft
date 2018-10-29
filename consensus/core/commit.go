package core

import (
	"github.com/Lywel/go-ibft/consensus"
	"reflect"
)

func (c *core) sendCommit() {
	c.logger.Log("send commit")
	subject := c.current.Subject()
	commit, err := Encode(subject)
	if err != nil {
		c.logger.Log("Encode subject failed")
		return
	}
	c.broadcast(&message{
		Type: typeCommit,
		Msg:  commit,
	})
}

func (c *core) handleCommit(msg *message, src *consensus.Validator) error {
	c.logger.Log("Handle commit from", src)
	var commit *consensus.Subject
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
		c.logger.Log("Failed to add COMMIT message", "msg", msg, "err", err)
	}
	if c.current.Commits.Size() > 2*c.valSet.F() && c.state.Cmp(StateCommitted) < 0 {
		c.commit()
	}
	return nil
}

func (c *core) verifyCommit(commit *consensus.Subject) error {
	subject := c.current.Subject()
	if !reflect.DeepEqual(commit, subject) {
		c.logger.Log("subjects do not match: expected", subject, "got", commit)
		return errSubjectsDoNotMatch
	}
	return nil
}
