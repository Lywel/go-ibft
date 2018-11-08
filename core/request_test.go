package core

import (
	"math/big"
	"sync"
	"testing"

	"bitbucket.org/ventureslash/go-ibft"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

func TestHandleRequest(t *testing.T) {
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
	}

	block := newBlockTest(ibft.Big1, "test")
	request := &ibft.Request{
		Proposal: block,
	}

	err := core.handleRequest(request)
	if err != nil {
		t.Errorf("handle request returned error")
		t.Log(err)
	}
}

func TestHandleOldRequest(t *testing.T) {
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
	}

	oldBlock := newBlockTest(ibft.Big0, "test")
	request := &ibft.Request{
		Proposal: oldBlock,
	}
	err := core.handleRequest(request)
	if err != errOldMessage {
		t.Errorf("handle request did not returned error")
		t.Log(err)
	}
}

func TestHandleFutureRequest(t *testing.T) {
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
	}

	futureBlock := newBlockTest(big.NewInt(2), "test")
	request := &ibft.Request{
		Proposal: futureBlock,
	}
	err := core.handleRequest(request)
	if err != errFutureMessage {
		t.Errorf("handle request did not returned error")
		t.Log(err)
	}
}
