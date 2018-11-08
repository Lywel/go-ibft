package core


func TestHandlePreprepare(t *testing.T) {
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
	}
}

