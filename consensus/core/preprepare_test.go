package core

import (
	"github.com/Lywel/ibft-go/consensus"
)

func newPreprepare(v *consensus.View) *consensus.Preprepare {
	return &consensus.Preprepare{
		View:     v,
		Proposal: newBlock(1),
	}
}
