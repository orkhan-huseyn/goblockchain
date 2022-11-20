package main

import (
	"github.com/orkhan-huseyn/goblockchain/blockchain"
	"github.com/orkhan-huseyn/goblockchain/wallet"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache map[string]*blockchain.Blockchain = make(map[string]*blockchain.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (s *BlockchainServer) Port() uint16 {
	return s.port
}

func (s *BlockchainServer) GetBlockchain() *blockchain.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = blockchain.NewBlockchain(minersWallet.Address(), s.port)
		cache["blockchain"] = bc

		log.Printf("miner's private key: %s\n", minersWallet.PrivateKeyStr())
		log.Printf("miner's public key: %s\n", minersWallet.PublicKeyStr())
		log.Printf("miner's address: %s\n", minersWallet.Address())
	}
	return bc
}

func (s *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bc := s.GetBlockchain()
		encoded, _ := bc.MarshalJSON()
		_, _ = io.WriteString(w, string(encoded[:]))
	default:
		log.Printf("Error: unsupported method\n")
	}
}

func (s *BlockchainServer) Run() {
	http.HandleFunc("/", s.GetChain)
	port := strconv.Itoa(int(s.port))
	log.Println("the server is running on port " + port)
	_ = http.ListenAndServe("0.0.0.0:"+port, nil)
}
