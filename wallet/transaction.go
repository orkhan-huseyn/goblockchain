package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"github.com/orkhan-huseyn/goblockchain/utils"
)

type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	amount                     float32
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, from string, to string, amount float32) *Transaction {
	return &Transaction{privateKey, publicKey, from, to, amount}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	encoded, _ := t.MarshalJSON()
	hash := sha256.Sum256(encoded)
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, hash[:])
	return &utils.Signature{R: r, S: s}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_address"`
		Recipient string  `json:"recipient_address"`
		Amount    float32 `json:"amount"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Amount:    t.amount,
	})
}
