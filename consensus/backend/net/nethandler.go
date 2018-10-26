package net

import (
	"log"

	"github.com/Lywel/go-gossipnet"
	"github.com/ethereum/go-ethereum/rlp"
)

// Handler handles data from the network
type Handler struct {
	node    *gossipnet.Node
	handler event.Handler
}

const (
	messageEvent = iota
	requestEvent
	backlogEvent
)

// Start starts to listen on node.EventChan()
func (mngr EventManager) Start() {
	for event := range mngr.node.EventChan() {
		switch event.(type) {
		case gossipnet.ConnOpenEvent:
			log.Print("ConnOpenEvent")
		case gossipnet.ConnCloseEvent:
			log.Print("ConnCloseEvent")
		case gossipnet.DataEvent:
			log.Print("DataEvent")
			rawMsg := event.(gossipnet.DataEvent).Data
			var msg networkMessage
			rlp.DecodeBytes(rawMsg, &msg)
			mngr.handleEvent()
		case gossipnet.ListenEvent:
			log.Print("ListenEvent")
		case gossipnet.CloseEvent:
			log.Print("CloseEvent")
			break
		}
	}
}

func (mngr EventManager) Broadcast(payload []byte) error {
	data, err := rlp.EncodeToBytes(networkMessage{
		Type: consensusMessage,
		Data: payload,
	})
	if err != nil {
		return err
	}
	mngr.Broadcast(data)
	return nil
}

func handleEvent(event Event) {
	switch event.(type) {

	}
}
