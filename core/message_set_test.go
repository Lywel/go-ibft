package core

import (
	"testing"

	"bitbucket.org/ventureslash/go-ibft"
)

func TestAdd(t *testing.T) {
	var a, b, c, d ibft.Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}, [20]byte{0, 1, 3}
	valSet := ibft.NewSet([]ibft.Address{a, b, c})
	msgSet := newMessageSet(valSet)

	ms := &message{
		Address: a,
	}
	invalidMs := &message{
		Address: d,
	}
	if err := msgSet.Add(ms); err != nil {
		t.Errorf("msg insertion failed")
	}
	if err := msgSet.Add(invalidMs); err == nil {
		t.Errorf("msg insertion should have failed but succeed")
	}
}

func TestGet(t *testing.T) {
	var a, b, c ibft.Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}
	valSet := ibft.NewSet([]ibft.Address{a, b, c})
	msgSet := newMessageSet(valSet)

	ms := &message{
		Address: a,
	}
	if err := msgSet.Add(ms); err != nil {
		t.Errorf("msg insertion failed")
	}
	if ms := msgSet.Get(a); ms == nil {
		t.Errorf("Get method failed")
	}
	if ms := msgSet.Get(b); ms != nil {
		t.Errorf("Get method failed")
	}
}

func TestSize(t *testing.T) {
	var a, b, c, d ibft.Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}, [20]byte{0, 1, 3}
	valSet := ibft.NewSet([]ibft.Address{a, b, c})
	msgSet := newMessageSet(valSet)
	ms := &message{
		Address: a,
	}
	invalidMs := &message{
		Address: d,
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
}
