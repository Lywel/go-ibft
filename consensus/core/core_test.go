package core

import (
	"math/big"

	"github.com/Lywel/go-ibft/types"
)

func newBlock(number int64) *types.Block {
	return types.NewBlock(big.NewInt(number), "test")
}
