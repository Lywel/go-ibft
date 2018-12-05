package ibft

import (
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)

func TestValidator(t *testing.T) {
	var a Address = [20]byte{0, 1, 2}
	v := NewValidator(a)
	if v.Address() != a {
		t.Errorf("address mismatch: expected %v, got %v", a.String(), v.String())
		t.FailNow()
	}
}

func newTestValidatorSet() *ValidatorSet {
	var a, b, c Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}
	return NewSet([]Address{a, b, c})
}

func TestAddValidator(t *testing.T) {
	var a, b, c, d Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}, [20]byte{0, 1, 3}
	set := NewSet([]Address{a, b, c})
	if !set.AddValidator(d) {
		t.Errorf("AddValidator failed")
	}
	if set.AddValidator(a) {
		t.Errorf("AddValidator should have failed but succeed")
	}
}

func TestRemoveValidator(t *testing.T) {
	var a, b, c Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}
	set := NewSet([]Address{a, b, c})

	if b := set.RemoveValidator(a); !b {
		t.Errorf("RemoveValidator failed")
	}
	if b := set.RemoveValidator(a); b {
		t.Errorf("RemoveValidator should have failed but succeed")
	}
}

func TestSize(t *testing.T) {
	var a, b, c, d Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}, [20]byte{0, 1, 3}
	set := NewSet([]Address{a, b, c})
	if set.Size() != 3 {
		t.Errorf("size mismatch: expected %d, got %d", 3, set.Size())
	}
	if !set.AddValidator(d) {
		t.Errorf("AddValidator failed")
	}
	if set.Size() != 4 {
		t.Errorf("size mismatch: expected %d, got %d", 4, set.Size())
	}
}

func TestGetByAddress(t *testing.T) {
	var a, b, c, d Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}, [20]byte{0, 1, 3}
	set := NewSet([]Address{a, b, c})
	if i, _ := set.GetByAddress(d); i != -1 {
		t.Errorf("GetByAddress failed: expected %d, got %d", -1, i)
	}
	if set.AddValidator(a) {
		t.Errorf("AddValidator should have failed but succeed")
	}
	if i, _ := set.GetByAddress(a); i == -1 {
		t.Errorf("GetByAddress failed: expected %d, got %d", -1, i)
	}
}

func TestIsProposer(t *testing.T) {
	var a, b, c, e Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}, [20]byte{0, 1, 2}
	set := NewSet([]Address{a, b, c})
	if res := set.IsProposer(e); res == false {
		t.Errorf("Bad proposer, expected %v, got %v", e.String(),
			set.GetProposer().String())
		t.Errorf("%v", set.validators)
	}
}

func TestRlp(t *testing.T) {
	var a, b, c, e Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}, [20]byte{0, 1, 2}
	set := NewSet([]Address{a, b, c, e})
	encodedSet, err := rlp.EncodeToBytes(set)
	if err != nil {
		t.Errorf("encode set failed")
		t.Log(err)
	}
	var decodedSet *ValidatorSet
	err = rlp.DecodeBytes(encodedSet, &decodedSet)
	if err != nil {
		t.Errorf("decode set failed")
		t.Log(err)
	}
	if set.Size() != decodedSet.Size() {
		t.Errorf("sets dont match")
	}
}

func TestUpdateValidator(t *testing.T) {
	var a, b, c Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}
	set := NewSet([]Address{a, b, c})
	if !set.IsProposer(a) {
		t.Errorf("wrong proposer, expected %v, got %v", a.String(), set.GetProposer().String())
	}
	set.UpdateProposer()
	if !set.IsProposer(b) {
		t.Errorf("wrong proposer, expected %v, got %v", b.String(), set.GetProposer().String())
	}
	set.UpdateProposer()
	if !set.IsProposer(c) {
		t.Errorf("wrong proposer, expected %v, got %v", c.String(), set.GetProposer().String())
	}
	set.UpdateProposer()
	if !set.IsProposer(a) {
		t.Errorf("wrong proposer, expected %v, got %v", a.String(), set.GetProposer().String())
	}
}
