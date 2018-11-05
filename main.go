package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/Lywel/go-ibft/consensus/backend"
	eth "github.com/ethereum/go-ethereum/crypto"
	"log"
	"os"
	"time"
)

func main() {
	privkey, err := ecdsa.GenerateKey(eth.S256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	backend := backend.New(&backend.Config{
		LocalAddr:   os.Args[1],
		RemoteAddrs: os.Args[2:],
	}, privkey)

	backend.Start()
	defer backend.Stop()
	time.Sleep(240 * time.Second)
}
