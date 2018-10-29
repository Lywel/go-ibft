package network

import (
	"github.com/Lywel/go-gossipnet"
	"github.com/Lywel/go-ibft/events"
	"github.com/ethereum/go-ethereum/rlp"
	"log"
)

// Manager handles data from the network
type Manager struct {
	node   *gossipnet.Node
	events events.Handler
}

type networkMessage struct {
	Type int
	Data []byte
}

// MessageEvent is emmitted during the IBFT consensus algo
type MessageEvent struct {
	Payload []byte
}

const (
	messageEvent = iota
	requestEvent
	backlogEvent
)

// Start starts to listen on node.EventChan()
func (mngr Manager) Start() {
	for event := range mngr.node.EventChan() {
		switch ev := event.(type) {
		case gossipnet.ConnOpenEvent:
			log.Print("ConnOpenEvent")
		case gossipnet.ConnCloseEvent:
			log.Print("ConnCloseEvent")
		case gossipnet.DataEvent:
			log.Print("DataEvent")
			var msg networkMessage
			rlp.DecodeBytes(ev.Data, &msg)
			mngr.events.Push(MessageEvent{
				Payload: msg.Data,
			})
		case gossipnet.ListenEvent:
			log.Print("ListenEvent")
		case gossipnet.CloseEvent:
			log.Print("CloseEvent")
			break
		}
	}
}

// Broadcast implements network.Manager.Broadcast. It will tag the payload
// forward it to the network node
func (mngr Manager) Broadcast(payload []byte) (err error) {
	data, err := rlp.EncodeToBytes(networkMessage{
		Type: messageEvent,
		Data: payload,
	})
	if err != nil {
		return
	}
	mngr.node.Broadcast(data)
	return
}
