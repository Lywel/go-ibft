package core

import (
	"bitbucket.org/ventureslash/go-ibft/consensus"
)

func newPreprepare(v *consensus.View) *consensus.Preprepare {
	return &consensus.Preprepare{
		View:     v,
		Proposal: newBlock(1),
	}
}
