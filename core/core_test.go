package core

import (
	"fmt"
	"io"
	"math/big"
	"testing"

	"bitbucket.org/ventureslash/go-ibft"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

type blockTest struct {
	number *big.Int
	data   string
}

func newBlockTest(number *big.Int, data string) *blockTest {
	return &blockTest{
		number: number,
		data:   data,
	}
}

// "external" block encoding. used for eth protocol, etc.
type extblockTest struct {
	Number *big.Int
	Data   string
}

type ProposalWrapper struct {
	Type     uint
	Proposal interface{}
}

const (
	blockType uint = iota
)

func RlpHash(x interface{}) []byte {
	var h common.Hash
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h.Bytes()
}

// Hash compute the hash of a block
func (b *blockTest) Hash() []byte {
	return RlpHash(b)
}

// Number return the number of a block
func (b *blockTest) Number() *big.Int {
	return b.number
}

func (b *blockTest) String() string {
	return fmt.Sprintf("number %d, data %s", b.number, b.data)
}

func (b *blockTest) EncodeRLP(w io.Writer) error {
	ext := extblockTest{
		Number: b.number,
		Data:   b.data,
	}
	propToBytes, err := rlp.EncodeToBytes(ext)
	if err != nil {
		return err
	}
	return rlp.Encode(w, ibft.EncodedProposal{
		Type: 1,
		Prop: propToBytes,
	})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *blockTest) DecodeRLP(s *rlp.Stream) error {
	var eb extblockTest

	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.number, b.data = eb.Number, eb.Data

	return nil
}

func TestConnection(t *testing.T) {
	
}

// TODO: integration tests
