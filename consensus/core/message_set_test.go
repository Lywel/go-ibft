package core

import (
	"testing"

	"github.com/Lywel/go-ibft/consensus"
)

var (
	a, b, c, d consensus.Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}, [20]byte{0, 1, 3}
)

func TestNewMessageSet(t *testing.T) {
	valSet := consensus.NewSet([]consensus.Address{a, b, c})
	msgSet := newMessageSet(valSet)
	ms := &message{
		Address: a,
	}
	invalidMs := &message{
		Address: d,
	}
	if ms == nil {
		t.Errorf("message initialization failed")
	}
	if msgSet == nil {
		t.Errorf("message set initialization failed")
	}
	if err := msgSet.Add(ms); err != nil {
		t.Errorf("msg insertion failed")
	}
	if msgSet.Size() != 1 {
		t.Errorf("invalid message set size: expected %d, got %d", 1, msgSet.Size())
	}

	if err := msgSet.Add(invalidMs); err == nil {
		t.Errorf("msg insertion should have failed but succeed")
	}
	if msgSet.Size() != 1 {
		t.Errorf("invalid message set size: expected %d, got %d", 1, msgSet.Size())
	}
	if ms := msgSet.Get(a); ms == nil {
		t.Errorf("Get method failed")
	}

}
