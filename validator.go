package ibft

import (
	"math"
	"reflect"
	"sort"
	"strings"
	"sync"
)

// Validator is a node of the consensus
type Validator struct {
	address Address
}

// NewValidator initialize a validator
func NewValidator(addr Address) *Validator {
	return &Validator{
		address: addr,
	}
}

// Address return address from validator
func (val *Validator) Address() Address {
	return val.address
}

// String convert validator to string representation
func (val *Validator) String() string {
	return val.Address().String()
}

// Validators is an array of validator
type Validators []*Validator

func (slice Validators) Len() int {
	return len(slice)
}

func (slice Validators) Less(i, j int) bool {
	return strings.Compare(slice[i].String(), slice[j].String()) < 0
}

func (slice Validators) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ValidatorSet contains a list of validators
type ValidatorSet struct {
	validators Validators

	proposer    *Validator
	validatorMu sync.RWMutex
}

// NewSet initialize a validatorSet
func NewSet(addrs []Address) *ValidatorSet {
	valSet := &ValidatorSet{}
	valSet.validators = make([]*Validator, len(addrs))
	for i, addr := range addrs {
		valSet.validators[i] = NewValidator(addr)
	}
	sort.Sort(valSet.validators)

	if valSet.Size() > 0 {
		valSet.proposer = valSet.validators[0]
	}
	return valSet

}

// Size returns the validator set size
func (valSet *ValidatorSet) Size() int {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return len(valSet.validators)
}

// List returns validator list
func (valSet *ValidatorSet) List() []*Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return valSet.validators
}

// AddValidator add a new validator to the set if he is not already present in the set
func (valSet *ValidatorSet) AddValidator(address Address) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()
	for _, v := range valSet.validators {
		if v.String() == address.String() {
			return false
		}
	}
	valSet.validators = append(valSet.validators, NewValidator(address))
	sort.Sort(valSet.validators)
	return true
}

// RemoveValidator remove a validator is he is in the set
func (valSet *ValidatorSet) RemoveValidator(address Address) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()
	for i, v := range valSet.validators {
		if v.String() == address.String() {
			valSet.validators = append(valSet.validators[:i], valSet.validators[i+1:]...)
			return true
		}
	}
	return false
}

// GetByAddress returns a validator and it's position corresponding to an address
func (valSet *ValidatorSet) GetByAddress(address Address) (int,
	*Validator) {
	for i, v := range valSet.validators {
		if v.String() == address.String() {
			return i, v
		}
	}
	return -1, nil
}

// GetProposer returns the current proposer
func (valSet *ValidatorSet) GetProposer() *Validator {
	return valSet.proposer
}

// IsProposer checks whether the validator with given address is a proposer
func (valSet *ValidatorSet) IsProposer(addr Address) bool {
	_, val := valSet.GetByAddress(addr)
	return reflect.DeepEqual(val, valSet.proposer)
}

// F gets the maximum number of faulty nodes
func (valSet *ValidatorSet) F() int { return int(math.Ceil(float64(valSet.Size())/3)) - 1 }
