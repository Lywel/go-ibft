package core

import (
	"bitbucket.org/ventureslash/go-ibft/consensus"
	"bitbucket.org/ventureslash/go-ibft/consensus/backend/crypto"
	"bitbucket.org/ventureslash/go-ibft/consensus/backend/network"
	"bitbucket.org/ventureslash/go-ibft/events"
	eth "github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
	"math/big"
	"sync"
)

type core struct {
	address               consensus.Address
	backend               consensus.Backend
	state                 State
	valSet                *consensus.ValidatorSet
	current               *roundState
	events                events.Handler
	networkManager        network.Manager
	pendingRequests       *prque.Prque
	pendingRequestsMu     *sync.Mutex
	backlogs              map[*consensus.Validator]*prque.Prque
	backlogsMu            *sync.Mutex
	logger                *Logger
	waitingForRoundChange bool
	wg                    sync.WaitGroup
}

// New initialize a new core
func New(backend consensus.Backend) consensus.Engine {
	eventHandler := events.New()
	networkManager := network.New(backend.Network(), eventHandler)
	view := &consensus.View{
		Round:    big.NewInt(0),
		Sequence: big.NewInt(0),
	}
	return &core{
		state: StateAcceptRequest,
		logger: &Logger{
			address: backend.Address(),
		},
		backend:           backend,
		pendingRequests:   prque.New(),
		pendingRequestsMu: &sync.Mutex{},
		backlogsMu:        &sync.Mutex{},
		backlogs:          make(map[*consensus.Validator]*prque.Prque),
		events:            eventHandler,
		networkManager:    networkManager,
		current:           newRoundState(view, nil, consensus.NewSet([]consensus.Address{backend.Address()}), nil),
		valSet:            consensus.NewSet([]consensus.Address{backend.Address()}),
	}
}

// Start implements core.Start
func (c *core) Start() {
	c.networkManager.Start(c.backend.Address())
	c.startNewRound(consensus.Big0)
	c.logger.Log("Core started")
	go c.handleEvents()
}

// Stop implements core.Stop
func (c *core) Stop() {
	c.logger.Log("Stopping the core")
	c.wg.Wait()
	c.logger.Log("Core stopped")
}

func (c *core) isValidator(a consensus.Address) bool {
	i, _ := c.valSet.GetByAddress(a)
	return i != -1
}

func (c *core) isProposer() bool {
	if c.valSet == nil {
		return false
	}
	return c.valSet.IsProposer(c.address)
}

func (c *core) currentView() *consensus.View {
	return &consensus.View{
		Round:    new(big.Int).Set(c.current.round),
		Sequence: new(big.Int).Set(c.current.sequence),
	}
}

func (c *core) finalizeMessage(msg *message) ([]byte, error) {
	msg.Address = c.address
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = c.backend.Sign(data)
	if err != nil {
		return nil, err
	}

	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (c *core) broadcast(msg *message) {
	payload, err := c.finalizeMessage(msg)
	if err != nil {
		c.logger.Log("failed to finalize message", "msg", msg, "err", err)
		return
	}
	if err := c.networkManager.Broadcast(payload); err != nil {
		c.logger.Log("failed to broadcast message", "msg", msg, "err", err)
	}
}

func (c *core) verify(p consensus.Proposal) error {
	return nil
}

func (c *core) checkMessage(msgType int, view *consensus.View) error {
	if view == nil || view.Sequence == nil || view.Round == nil {
		return errInvalidMessage
	}

	if msgType == typeRoundChange {
		if view.Sequence.Cmp(c.currentView().Sequence) > 0 {
			return errFutureMessage
		} else if view.Cmp(c.currentView()) < 0 {
			return errOldMessage
		}
		return nil
	}

	if view.Cmp(c.currentView()) > 0 {
		return errFutureMessage
	}

	if view.Cmp(c.currentView()) < 0 {
		return errOldMessage
	}

	if c.waitingForRoundChange {
		return errFutureMessage
	}

	// StateAcceptRequest only accepts msgPreprepare
	// other messages are future messages
	if c.state == StateAcceptRequest {
		if msgType > typePreprepare {
			return errFutureMessage
		}
		return nil
	}

	// For states(StatePreprepared, StatePrepared, StateCommitted),
	// can accept all message types if processing with same view
	return nil
}

func (c *core) commit() {
	c.setState(StateCommitted)
	c.logger.Log("committed")
	c.startNewRound(consensus.Big0)
}

func (c *core) setState(state State) {
	c.state = state
	if state == StateAcceptRequest {
		c.processPendingRequests()
	}
	c.processBacklogs()
}

// startNewRound starts a new round. if round equals to 0, it means to starts a
// new sequence
func (c *core) startNewRound(round *big.Int) {
	c.logger.Log("start new sequence")
	roundChange := false
	// TODO check if there is a round change

	var view *consensus.View
	if roundChange {
		view = &consensus.View{
			Sequence: new(big.Int).Set(c.current.sequence),
			Round:    new(big.Int).Set(round),
		}
	} else {
		view = &consensus.View{
			Sequence: new(big.Int).Add(c.current.sequence, consensus.Big1),
			Round:    new(big.Int),
		}
		// TODO update validators with new list

		c.waitingForRoundChange = false

		c.updateRoundState(view, c.valSet, roundChange)
		c.setState(StateAcceptRequest)
	}
}

func (c *core) updateRoundState(view *consensus.View, valSet *consensus.ValidatorSet,
	roundChange bool) {
	if roundChange {
		c.logger.Log("update round")
		// TODO round change
	} else {
		c.current = newRoundState(view, nil, valSet, nil)
	}
}

func (c *core) ValidateFn(data []byte, sig []byte) (consensus.Address, error) {
	hashData := eth.Keccak256([]byte(data))
	signer, err := eth.SigToPub(hashData, sig)
	if err != nil {
		return consensus.Address{}, err
	}
	address := crypto.PubkeyToAddress(*signer)
	i, _ := c.valSet.GetByAddress(address)
	if i == -1 {
		return address, errUnauthorized
	}
	return address, nil
}
