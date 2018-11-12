package types

import (
	"bitbucket.org/ventureslash/go-ibft"
)

// ValidatorState is sent when a new validator wants to join the network
type ValidatorState struct {
	ValSet *ibft.ValidatorSet
	View   *ibft.View
}
