package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/Lywel/go-ibft/consensus/backend"
	"github.com/Lywel/go-ibft/consensus/core"
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

	core := core.New(&backend)

	core.Start()
	defer core.Stop()
	time.Sleep(20 * time.Second)
}
