package core

import (
	"fmt"
	"io"
	"math/big"

	"bitbucket.org/ventureslash/go-ibft"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

// Transaction represents a transaction sent over the network
type transactionTest struct {
	From      ibft.Address
	To        ibft.Address
	Amount    *big.Int
	Signature []byte
}

type transactionsTest []*transactionTest

type headerTest struct {
	Number     *big.Int
	ParentHash []byte
}

// Block is used to build the blockchain
type blockTest struct {
	Header       *headerTest
	Transactions transactionsTest
}

// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	Header       *headerTest
	Transactions transactionsTest
}

// NewBlock create a new bock
func newBlockTest(header *headerTest, transactions []*transactionTest) *blockTest {
	return &blockTest{
		Header:       header,
		Transactions: transactions,
	}
}

// Hash compute the hash of a block
func (b *blockTest) Hash() []byte {
	return RlpHash(b.Header)
}

// Number return the number of a block
func (b *blockTest) Number() *big.Int {
	return new(big.Int).Set(b.Header.Number)
}

func (b *blockTest) String() string {
	return fmt.Sprintf("number %d", b.Number())
}

// EncodeRLP TODO
func (b *blockTest) EncodeRLP(w io.Writer) error {
	ext := extblock{
		Header:       b.Header,
		Transactions: b.Transactions,
	}
	propToBytes, err := rlp.EncodeToBytes(ext)
	if err != nil {
		return err
	}
	return rlp.Encode(w, ibft.EncodedProposal{
		Type: 0,
		Prop: propToBytes,
	})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *blockTest) DecodeRLP(s *rlp.Stream) error {
	var eb extblock

	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.Header, b.Transactions = eb.Header, eb.Transactions

	return nil
}

func RlpHash(x interface{}) []byte {
	var h common.Hash
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h.Bytes()
}
