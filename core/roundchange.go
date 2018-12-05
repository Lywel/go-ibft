package core

import (
	"math/big"
	"sync"

	"bitbucket.org/ventureslash/go-ibft"
)

func newRoundChangeSet(valSet *ibft.ValidatorSet) *roundChangeSet {
	return &roundChangeSet{
		validatorSet: valSet,
		roundChanges: make(map[uint64]*messageSet),
		mu:           new(sync.Mutex),
	}
}

type roundChangeSet struct {
	validatorSet *ibft.ValidatorSet
	roundChanges map[uint64]*messageSet
	mu           *sync.Mutex
}

// Add adds the round and message into round change set
func (rcs *roundChangeSet) Add(r *big.Int, msg *message) (int, error) {
	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	round := r.Uint64()
	if rcs.roundChanges[round] == nil {
		rcs.roundChanges[round] = newMessageSet(rcs.validatorSet)
	}
	err := rcs.roundChanges[round].Add(msg)
	if err != nil {
		return 0, err
	}
	return rcs.roundChanges[round].Size(), nil
}

// Clear deletes the messages with smaller round
func (rcs *roundChangeSet) Clear(round *big.Int) {
	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	for k, rms := range rcs.roundChanges {
		if len(rms.Values()) == 0 || k < round.Uint64() {
			delete(rcs.roundChanges, k)
		}
	}
}

// MaxRound returns the max round which the number of messages is equal or larger than num
func (rcs *roundChangeSet) MaxRound(num int) *big.Int {
	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	var maxRound *big.Int
	for k, rms := range rcs.roundChanges {
		if rms.Size() < num {
			continue
		}
		r := big.NewInt(int64(k))
		if maxRound == nil || maxRound.Cmp(r) < 0 {
			maxRound = r
		}
	}
	return maxRound
}

// ----------------------------------------------------------

func (c *core) sendNextRoundChange() {
	cv := c.currentView()
	c.sendRoundChange(new(big.Int).Add(cv.Round, ibft.Big1))
}

func (c *core) sendRoundChange(round *big.Int) {
	c.logger.Info(c.address, ": Send round change for round ", round)
	cv := c.currentView()

	// TODO catchup round
	c.waitingForRoundChange = true

	newView := &ibft.View{
		// The round number we'd like to transfer to.
		Round:    new(big.Int).Set(round),
		Sequence: new(big.Int).Set(cv.Sequence),
	}
	if c.current != nil {
		c.current = newRoundState(newView, c.current.Preprepare, c.valSet, c.current.pendingRequest)
	}

	cv = c.currentView()
	rc := &ibft.Subject{
		View:   cv,
		Digest: []byte{},
	}

	payload, err := Encode(rc)
	if err != nil {
		c.logger.Warning(c.address, ": Failed to encode ROUND CHANGE ", "rc ", rc, " err ", err)
		return
	}

	c.broadcast(&message{
		Type: typeRoundChange,
		Msg:  payload,
	})
}

func (c *core) handleRoundChange(msg *message, src *ibft.Validator) error {
	c.logger.Info(c.address, ": Handle round change from ", src.String())
	var roundchange *ibft.Subject
	err := msg.Decode(&roundchange)
	if err != nil {
		return errFailedDecodeRoundChange
	}

	if err = c.checkMessage(typeRoundChange, roundchange.View); err != nil {
		return err
	}

	cv := c.currentView()
	roundView := roundchange.View

	// Add the ROUND CHANGE message to its message set and return how many
	// messages we've got with the same round number and sequence number.
	num, err := c.roundChangeSet.Add(roundView.Round, msg)
	if err != nil {
		c.logger.Warning(c.address, ": Failed to add round change message ", "from ", src, " msg ", msg, " err ", err)
		return err
	}

	// Once we received f+1 ROUND CHANGE messages, those messages form a weak certificate.
	// If our round number is smaller than the certificate's round number, we would
	// try to catch up the round number.
	if c.waitingForRoundChange && num == int(c.valSet.F()+1) {
		if cv.Round.Cmp(roundView.Round) < 0 {
			c.sendRoundChange(roundView.Round)
		}
		return nil
	} else if num >= int(2*c.valSet.F()+1) && (c.waitingForRoundChange || cv.Round.Cmp(roundView.Round) < 0) {
		c.logger.Info("got ", num, " signatures, needed more than ", 2*c.valSet.F())

		// We've received 2f+1 ROUND CHANGE messages, start a new round immediately.
		c.startNewRound(roundView.Round)
		return nil
	} else {
		c.logger.Info("got ", num, " signatures, expected more than ", 2*c.valSet.F())

	}
	return nil
}
