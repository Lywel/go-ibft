package core

import (
	"crypto/ecdsa"
	"flag"
	"io/ioutil"
	"math/big"
	"sync"
	"time"

	"bitbucket.org/ventureslash/go-ibft"
	"bitbucket.org/ventureslash/go-ibft/crypto"
	eth "github.com/ethereum/go-ethereum/crypto"
	"github.com/google/logger"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

var verbose = flag.Bool("verbose-core", false, "print core info level logs")

type core struct {
	address               ibft.Address
	privateKey            *ecdsa.PrivateKey
	backend               backend
	state                 State
	valSet                *ibft.ValidatorSet
	current               *roundState
	eventsIn              chan Event
	eventsOut             chan Event
	pendingRequests       *prque.Prque
	pendingRequestsMu     *sync.Mutex
	backlogs              map[*ibft.Validator]*prque.Prque
	backlogsMu            *sync.Mutex
	logger                *logger.Logger
	waitingForRoundChange bool
	wg                    sync.WaitGroup
	proposalManager       ibft.ProposalManager
	timeouts              map[*ibft.Validator]*time.Timer
	timeoutsMu            *sync.Mutex
	roundChangeSet        *roundChangeSet
	roundChangeTimer      *time.Timer
	networkMap            map[ibft.Address]string
}

// New initialize a new core
func New(b backend, proposalManager ibft.ProposalManager) ibft.Core {
	//networkManager := network.New(backend.Network(), eventHandler)
	address := crypto.PubkeyToAddress(b.PrivateKey().PublicKey)
	view := &ibft.View{
		Round:    big.NewInt(0),
		Sequence: big.NewInt(0),
	}
	c := &core{
		address:           address,
		privateKey:        b.PrivateKey(),
		state:             StateAcceptRequest,
		logger:            logger.Init("Core", *verbose, false, ioutil.Discard),
		backend:           b,
		pendingRequests:   prque.New(),
		pendingRequestsMu: &sync.Mutex{},
		backlogs:          make(map[*ibft.Validator]*prque.Prque),
		backlogsMu:        &sync.Mutex{},
		eventsIn:          b.EventsInChan(),
		eventsOut:         b.EventsOutChan(),
		current:           newRoundState(view, nil, ibft.NewSet([]ibft.Address{address}), nil),
		valSet:            ibft.NewSet([]ibft.Address{address}),
		proposalManager:   proposalManager,
		timeouts:          make(map[*ibft.Validator]*time.Timer),
		timeoutsMu:        &sync.Mutex{},
		networkMap:        make(map[ibft.Address]string),
	}
	c.logger.Infof("Start")
	return c
}

func (c *core) NetworkMap() map[ibft.Address]string {
	for key := range c.networkMap {
		if !c.isValidator(key) {
			delete(c.networkMap, key)
		}
	}
	return c.networkMap
}

// Start implements core.Start
func (c *core) Start() {
	//c.eventsOut <- JoinEvent{Address: c.address}
	c.startNewRound(ibft.Big0)
	c.logger.Info("Core started")
	go c.handleEvents()
}

// Stop implements core.Stop
func (c *core) Stop() {
	c.logger.Info(c.address, ": Stopping the core")
	c.wg.Wait()
	c.logger.Info(c.address, ": Core stopped")
}

func (c *core) isValidator(a ibft.Address) bool {
	i, _ := c.valSet.GetByAddress(a)
	return i != -1
}

func (c *core) isProposer() bool {
	if c.valSet == nil {
		return false
	}
	return c.valSet.IsProposer(c.address)
}

func (c *core) currentView() *ibft.View {
	return &ibft.View{
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
	msg.Signature, err = c.sign(data)
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
		c.logger.Warning(c.address, ": Failed to finalize message ", "msg ", msg, " err ", err)
		return
	}
	// Broadcast
	c.eventsOut <- MessageEvent{Payload: payload}
}

func (c *core) verify(p ibft.Proposal) error {
	return c.proposalManager.Verify(p)
}

func (c *core) checkMessage(msgType uint64, view *ibft.View) error {
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
	proposal := c.current.Preprepare.Proposal
	if proposal != nil {
		if err := c.proposalManager.Commit(c.current.Preprepare.Proposal); err == nil {
			c.logger.Info(c.address, ": Proposal committed")
			c.startNewRound(ibft.Big0)

		} else {
			c.logger.Warning(c.address, ": Commit failed ", "err ", err)
			//TODO trigger change view
		}
	}
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
	c.logger.Info(c.address, ": Start new sequence")
	roundChange := false

	// TODO UpdateSince

	if round.Cmp(ibft.Big0) > 0 {
		roundChange = true
	}

	var view *ibft.View
	if roundChange {
		view = &ibft.View{
			Sequence: new(big.Int).Set(c.current.sequence),
			Round:    new(big.Int).Set(round),
		}
		c.valSet.UpdateProposer()
	} else {
		view = &ibft.View{
			Sequence: new(big.Int).Add(c.current.sequence, ibft.Big1),
			Round:    new(big.Int),
		}
	}
	// TODO update validators with new list

	c.roundChangeSet = newRoundChangeSet(c.valSet)
	c.waitingForRoundChange = false

	c.updateRoundState(view, c.valSet, roundChange)
	c.setState(StateAcceptRequest)

	if roundChange && c.isProposer() && c.current != nil && c.current.pendingRequest != nil {
		c.sendPreprepare(c.current.pendingRequest)
	}
	c.newRoundChangeTimer()

}

func (c *core) updateRoundState(view *ibft.View, valSet *ibft.ValidatorSet,
	roundChange bool) {
	if roundChange {
		c.current = newRoundState(view, c.current.Preprepare, valSet, c.current.pendingRequest)
	} else {
		c.current = newRoundState(view, nil, valSet, nil)
	}
}

func (c *core) newRoundChangeTimer() {
	if c.roundChangeTimer != nil {
		c.roundChangeTimer.Stop()
	}
	c.roundChangeTimer = time.AfterFunc(ibft.RequestTimeout, func() {
		c.logger.Info(c.address, ": Round change triggered")
		c.handleTimeoutMsg()
	})
}

func (c *core) ValidateFn(data []byte, sig []byte) (ibft.Address, error) {
	hashData := eth.Keccak256([]byte(data))
	signer, err := eth.SigToPub(hashData, sig)
	if err != nil {
		return ibft.Address{}, err
	}
	address := crypto.PubkeyToAddress(*signer)
	i, _ := c.valSet.GetByAddress(address)
	if i == -1 {
		return address, errUnauthorized
	}
	return address, nil
}

func (c *core) sign(data []byte) ([]byte, error) {
	return crypto.Sign(data, c.privateKey)
}
