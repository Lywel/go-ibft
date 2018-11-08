package core

import (
	"testing"

	"bitbucket.org/ventureslash/go-ibft"
)

func TestEncode(t *testing.T) {

	block := newBlockTest(ibft.Big1, "test")

	preprepare, err := Encode(&ibft.Preprepare{
		View: &ibft.View{
			Sequence: ibft.Big1,
			Round:    ibft.Big0,
		},
		Proposal: block,
	})
	if err != nil {
		t.Errorf("failed encode preprepare")
	}
	t.Log(preprepare)
}

func TestDecodePrePrepare(t *testing.T) {
	var a ibft.Address = [20]byte{0, 1, 2}

	block := newBlockTest(ibft.Big1, "test")

	pre, err := Encode(&ibft.Preprepare{
		View: &ibft.View{
			Sequence: ibft.Big1,
			Round:    ibft.Big0,
		},
		Proposal: block,
	})
	if err != nil {
		t.Errorf("failed encode preprepare")
	}
	ms := &message{
		Type:    typePreprepare,
		Msg:     pre,
		Address: a,
	}
	var preprepare *ibft.Preprepare

	err = ms.Decode(&preprepare)
	if err != nil {
		t.Errorf("decode failed")
		t.Log(err)
	}
}

func TestDecodeSubject(t *testing.T) {
	var a ibft.Address = [20]byte{0, 1, 2}
	subject := &ibft.Subject{
		View: &ibft.View{
			Sequence: ibft.Big1,
			Round:    ibft.Big0,
		},
		Digest: []byte{9},
	}
	pre, err := Encode(subject)
	if err != nil {
		t.Errorf("Encode subject failed")
		return
	}
	ms := &message{
		Type:    typePrepare,
		Msg:     pre,
		Address: a,
	}

	var prepare *ibft.Subject
	err = ms.Decode(&prepare)
	if err != nil {
		t.Errorf("decode failed")
		t.Log(err)
	}
}
