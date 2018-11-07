package ibft

import (
	"testing"
)

var (
	a, b, c, d, e Address = [20]byte{0, 1, 2}, [20]byte{0, 1, 2, 4}, [20]byte{0, 1, 2, 5}, [20]byte{0, 1, 3}, [20]byte{0, 1, 2}
)

func TestValidator(t *testing.T) {
	v := NewValidator(a)
	if v.Address() != a {
		t.Errorf("address mismatch: expected %v, got %v", a.String(), v.String())
		t.FailNow()
	}
}

func TestValidatorSet(t *testing.T) {
	set := NewSet([]Address{a, b, c})
	if set.Size() != 3 {
		t.Errorf("size mismatch: expected %d, got %d", 3, set.Size())
	}
	if i, _ := set.GetByAddress(d); i != -1 {
		t.Errorf("GetByAddress failed: expected %d, got %d", -1, i)
	}
	if !set.AddValidator(d) {
		t.Errorf("AddValidator failed")
	}
	if set.Size() != 4 {
		t.Errorf("size mismatch: expected %d, got %d", 4, set.Size())
	}
	if set.AddValidator(a) {
		t.Errorf("AddValidator should have failed but succeed")
	}
	if i, _ := set.GetByAddress(a); i == -1 {
		t.Errorf("GetByAddress failed: expected %d, got %d", -1, i)
	}
	if res := set.IsProposer(e); res == false {
		t.Errorf("Bad proposer, expected %v, got %v", e.String(),
			set.GetProposer().String())
		t.Errorf("%v", set.validators)
	}
	if b := set.RemoveValidator(a); !b {
		t.Errorf("RemoveValidator failed")
	}
	if b := set.RemoveValidator(a); b {
		t.Errorf("RemoveValidator should have failed but succeed")
	}

}
