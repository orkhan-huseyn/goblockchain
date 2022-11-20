package blockchain

import (
	"fmt"
	"strings"
)

type Transaction struct {
	SenderAddress    string  `json:"sender_address"`
	RecipientAddress string  `json:"recipient_address"`
	Amount           float32 `json:"amount"`
}

func NewTransaction(from string, to string, amount float32) *Transaction {
	return &Transaction{
		SenderAddress:    from,
		RecipientAddress: to,
		Amount:           amount,
	}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 62))
	fmt.Printf("from:    %s\n", t.SenderAddress)
	fmt.Printf("to:      %s\n", t.RecipientAddress)
	fmt.Printf("amount:  %.1f\n", t.Amount)
}
