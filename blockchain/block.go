package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	nonce        int
	prevHash     [32]byte
	timestamp    int64
	transactions []*Transaction
}

func NewBlock(nonce int, prevHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		timestamp:    time.Now().UnixNano(),
		prevHash:     prevHash,
		nonce:        nonce,
		transactions: transactions,
	}
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PrevHash     string         `json:"previous_hash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Nonce:        b.nonce,
		PrevHash:     fmt.Sprintf("%x", b.prevHash),
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
	})
}

func (b *Block) Print() {
	fmt.Printf("timestamp      %d\n", b.timestamp)
	fmt.Printf("nonce          %d\n", b.nonce)
	fmt.Printf("previous_hash  %x\n", b.prevHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	encoded, _ := json.Marshal(b)
	return sha256.Sum256(encoded)
}
