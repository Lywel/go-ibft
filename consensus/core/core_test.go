package core

import (
	"math/big"

	"github.com/Lywel/ibft-go/types"
)

func newBlock(number int64) *types.Block {
	return types.NewBlock(big.NewInt(number), "test")
}
