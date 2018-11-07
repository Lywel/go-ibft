package core

import (
	"bitbucket.org/ventureslash/go-ibft"
)

func newPreprepare(v *ibft.View) *ibft.Preprepare {
	return &ibft.Preprepare{
		View:     v,
		Proposal: newBlock(1),
	}
}
