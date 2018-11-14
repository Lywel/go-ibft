package core

import (
	"crypto/ecdsa"
	"errors"
	"reflect"
	"testing"

	"bitbucket.org/ventureslash/go-ibft"
	"github.com/ethereum/go-ethereum/rlp"
)

type backendTest struct{}

type proposalManagerTest struct{}

func (b *backendTest) PrivateKey() *ecdsa.PrivateKey {
	return nil
}

func (b *backendTest) Start() {}

func (b *backendTest) Stop() {}

// EventsInChan returns a channel receiving network events
func (b *backendTest) EventsInChan() chan Event {
	return nil
}

// EventsOutChan returns a channel used to emit events to the network
func (b *backendTest) EventsOutChan() chan Event {
	return nil
}

// DecodeProposal parses a payload and return a Proposal interface
func (p *proposalManagerTest) DecodeProposal(prop *ibft.EncodedProposal) (ibft.Proposal, error) {
	switch prop.Type {
	case 1:
		var b *blockTest
		err := rlp.DecodeBytes(prop.Prop, &b)
		if err != nil {
			return nil, err
		}
		return b, nil
		/*
			case type.Transactiono:
				return proposal
		*/
	default:
		return nil, errors.New("Unknown proposal type " + reflect.ValueOf(prop).Elem().String())
	}
}

// Verify returns an error is a proposal should be rejected
func (p *proposalManagerTest) Verify(proposal ibft.Proposal) error { return nil }

// Commit is called by an IBFT algorythm when a Proposal is accepted
func (p *proposalManagerTest) Commit(proposal ibft.Proposal) error { return nil }

func TestHandlePreprepare(t *testing.T) {
	/*privkey, _ := ecdsa.GenerateKey(eth.S256(), rand.Reader)
	proposalManager := &proposalManagerTest{}

	var a, b, c ibft.Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}
	valSet := ibft.NewSet([]ibft.Address{a, b, c})
	core := &core{
		address: b,
		state:   StateAcceptRequest,
		current: newRoundState(&ibft.View{
			Sequence: big.NewInt(1),
			Round:    big.NewInt(0),
		}, nil, valSet, nil),
		backlogsMu: &sync.Mutex{},
		backlogs:   make(map[*ibft.Validator]*prque.Prque),
		valSet:     valSet,
		logger: &Logger{
			address: ibft.Address{0},
		},
		proposalManager: proposalManager,
		privateKey: privkey,
	}

	block := newBlockTest(ibft.Big1, "test")

	preprepare, err := Encode(&ibft.Preprepare{
		View:     core.currentView(),
		Proposal: block,
	})
	if err != nil {
		t.Errorf("failed encode preprepare")
	}
	ms := &message{
		Type:    typePreprepare,
		Msg:     preprepare,
		Address: a,
	}
	_, val := core.valSet.GetByAddress(a)

	err = core.handlePreprepare(ms, val)
	if err != nil {
		t.Errorf("handle preprepare failed")
		t.Log(err)
	}*/
}
