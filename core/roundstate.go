package core

import (
	"math/big"
	"sync"

	"bitbucket.org/ventureslash/go-ibft"
)

// roundState stores the consensus state
type roundState struct {
	round          *big.Int
	sequence       *big.Int
	Preprepare     *ibft.Preprepare
	Prepares       *messageSet
	Commits        *messageSet
	pendingRequest *ibft.Request
	mu             *sync.RWMutex
}

// newRoundState creates a new roundState instance with the given view and validatorSet
func newRoundState(view *ibft.View, preprepare *ibft.Preprepare,
	valSet *ibft.ValidatorSet, request *ibft.Request) *roundState {
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
func (s *roundState) Subject() *ibft.Subject {
	return &ibft.Subject{
		View: &ibft.View{
			Round:    new(big.Int).Set(s.round),
			Sequence: new(big.Int).Set(s.sequence),
		},
		Digest: s.Preprepare.Proposal.Hash(),
	}
}

func (s *roundState) GetPrepareOrCommitSize() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := s.Prepares.Size() + s.Commits.Size()

	// find duplicate one
	for _, m := range s.Prepares.Values() {
		if s.Commits.Get(m.Address) != nil {
			result--
		}
	}
	return result
}
