package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/orkhan-huseyn/goblockchain/utils"
	"log"
	"strings"
)

const (
	MiningDifficulty = 3
	MiningSender     = "[TheBlockchain]"
	MiningReward     = 1.0
)

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
	address         string
	port            uint16
}

func NewBlockchain(address string, port uint16) *Blockchain {
	bc := new(Blockchain)
	bc.address = address
	bc.CreateBlock(0, [32]byte{})
	bc.port = port
	return bc
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Chain []*Block `json:"chain"`
	}{
		Chain: bc.chain,
	})
}

func (bc *Blockchain) CreateBlock(nonce int, prevHash [32]byte) *Block {
	b := NewBlock(nonce, prevHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) AddTransaction(from string, to string, amount float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(from, to, amount)

	if from == MiningSender {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		//if bc.CalculateTotalAmount(from) < amount {
		//	log.Println("ERROR: Not enough balance in a wallet")
		//	return false
		//}

		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	log.Println("ERROR: Verify transaction")
	return false
}

func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	encoded, _ := json.Marshal(t)
	hash := sha256.Sum256(encoded)
	return ecdsa.Verify(senderPublicKey, hash[:], s.R, s.S)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, len(bc.transactionPool))
	for _, t := range bc.transactionPool {
		transactions = append(transactions, &Transaction{
			SenderAddress:    t.SenderAddress,
			RecipientAddress: t.RecipientAddress,
			Amount:           t.Amount,
		})
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, prevHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeroes := strings.Repeat("0", difficulty)
	guessBlock := Block{nonce, prevHash, 0, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeroes
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	prevHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, prevHash, transactions, MiningDifficulty) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mine() {
	bc.AddTransaction(MiningSender, bc.address, MiningReward, nil, nil)
	nonce := bc.ProofOfWork()
	prevHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, prevHash)
	log.Println("action=mining, status=success")
}

func (bc *Blockchain) CalculateTotalAmount(address string) float32 {
	var total float32
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.Amount
			if address == t.SenderAddress {
				total -= value
			}
			if address == t.RecipientAddress {
				total += value
			}
		}
	}
	return total
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Block %04d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("=", 62))
}
