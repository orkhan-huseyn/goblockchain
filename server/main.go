package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", 5000, "TCP port number for blockchain server")
	flag.Parse()
	server := NewBlockchainServer(uint16(*port))
	server.Run()
}
