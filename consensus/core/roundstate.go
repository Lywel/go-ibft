package core

import (
	"math/big"
	"sync"

	"github.com/Lywel/ibft-go/consensus"
)

// roundState stores the consensus state
type roundState struct {
	round          *big.Int
	sequence       *big.Int
	Preprepare     *consensus.Preprepare
	Prepares       *messageSet
	Commits        *messageSet
	pendingRequest *consensus.Request
	mu             *sync.RWMutex
}

// newRoundState creates a new roundState instance with the given view and validatorSet
func newRoundState(view *consensus.View, preprepare *consensus.Preprepare,
	valSet *consensus.ValidatorSet, request *consensus.Request) *roundState {
	return &roundState{
		round:          view.Round,
		sequence:       view.Sequence,
		Preprepare:     preprepare,
		Prepares:       newMessageSet(valSet),
		Commits:        newMessageSet(valSet),
		pendingRequest: request,
		mu:             new(sync.RWMutex),
	}
}

// Subject returns the subject of the current round
func (s *roundState) Subject() *consensus.Subject {
	return &consensus.Subject{
		View: &consensus.View{
			Round:    new(big.Int).Set(s.round),
			Sequence: new(big.Int).Set(s.sequence),
		},
		Digest: s.Preprepare.Proposal.Hash(),
	}
}
