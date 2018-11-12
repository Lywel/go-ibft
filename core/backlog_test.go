package core

import (
	"math/big"
	"testing"

	"bitbucket.org/ventureslash/go-ibft"
)

func TestStoreBacklog(t *testing.T) {
	/*privkey, _ := ecdsa.GenerateKey(eth.S256(), rand.Reader)
	var a, b, c ibft.Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}
	backend := &backendTest{}
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
		backend:    backend,
		privateKey: privkey,
	}

	block := newBlockTest(ibft.Big1, "test")

	preprepare, err := Encode(&ibft.Preprepare{
		View:     core.currentView(),
		Proposal: block,
	})
	if err != nil {
		t.Errorf("failed encode preprepare")
		t.Log(err)
	}
	ms := &message{
		Type:    typePreprepare,
		Msg:     preprepare,
		Address: a,
	}
	_, val := core.valSet.GetByAddress(a)
	core.storeBacklog(ms, val)

	backlog := core.backlogs[val]

	if backlog.Empty() {
		t.Errorf("failed store backlog")
	}*/
}

func TestToPriority(t *testing.T) {
	view1 := &ibft.View{
		Sequence: big.NewInt(0),
		Round:    big.NewInt(0),
	}
	view2 := &ibft.View{
		Sequence: big.NewInt(1),
		Round:    big.NewInt(0),
	}

	prioView1Preprepare := toPriority(typePreprepare, view1)
	prioView1Prepare := toPriority(typePrepare, view1)
	prioView2Preprepare := toPriority(typePreprepare, view2)
	prioView2Prepare := toPriority(typePrepare, view2)

	if prioView1Preprepare <= prioView1Prepare {
		t.Errorf("view 1 preprepare <= view 1 prepare")
		t.Log(prioView1Preprepare, "<=", prioView1Prepare)
	}
	if prioView1Prepare <= prioView2Preprepare {
		t.Errorf("view 1 prepare <= view 2 preprepare")
		t.Log(prioView2Prepare, "<=", prioView2Preprepare)
	}
	if prioView2Preprepare <= prioView2Prepare {
		t.Errorf("view 2 preprepare <= view 2 prepare")
		t.Log(prioView2Preprepare, "<=", prioView2Prepare)
	}
}
