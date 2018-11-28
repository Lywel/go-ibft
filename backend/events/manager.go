package events

import (
	"log"

	"bitbucket.org/ventureslash/go-gossipnet"
	"bitbucket.org/ventureslash/go-ibft"
	"bitbucket.org/ventureslash/go-ibft/core"
	"github.com/ethereum/go-ethereum/rlp"
)

// Manager handles data from the network
type Manager struct {
	node           *gossipnet.Node
	eventsOut      chan core.Event
	eventsIn       chan core.Event
	eventsCustom   chan []byte
	nodesToConnect int
}

type networkMessage struct {
	Type uint
	Data []byte
}

const (
	messageEvent uint = iota
	requestEvent
	backlogEvent
	joinEvent
	stateEvent
	addValidatorEvent
	customEvent
)

// New returns a new network.Manager
func New(node *gossipnet.Node, eventsIn, eventsOut chan core.Event, eventsCustom chan []byte, nodesToConnect int) Manager {
	return Manager{
		node:           node,
		eventsIn:       eventsIn,
		eventsOut:      eventsOut,
		eventsCustom:   eventsCustom,
		nodesToConnect: nodesToConnect,
	}
}

// Start Broadcast core address and starts to listen on node.EventChan()
func (mngr Manager) Start(addr ibft.Address) {
	addrBytes := addr.GetBytes()
	joinBytes, err := rlp.EncodeToBytes(networkMessage{
		Type: joinEvent,
		Data: addrBytes[:],
	})
	if err != nil {
		log.Print("encode error: ", err)
	}

	// Dispatch network events to IBFT
	go func() {
		for event := range mngr.node.EventChan() {
			switch ev := event.(type) {
			case gossipnet.ConnOpenEvent:
				log.Print("ConnOpenEvent")
				// TODO: dont gossip to everyone, just the new connection
				if mngr.nodesToConnect > 0 {
					mngr.node.Gossip(joinBytes)
					mngr.nodesToConnect--
				}
			case gossipnet.ConnCloseEvent:
				log.Print("ConnCloseEvent")
			case gossipnet.DataEvent:
				log.Print("DataEvent")
				var msg networkMessage
				err := rlp.DecodeBytes(ev.Data, &msg)
				if err != nil {
					log.Print("Error parsing msg:", string(ev.Data))
					continue
				}
				switch msg.Type {
				case messageEvent:
					log.Print(" -MsgEvent")
					mngr.eventsIn <- core.MessageEvent{
						Payload: msg.Data,
					}
				case joinEvent:
					log.Print(" -JoinEvent")
					evt := core.JoinEvent{
						NetworkAddr: ev.Addr,
					}
					evt.Address.FromBytes(msg.Data)
					if err != nil {
						log.Print(err)
						continue
					}
					mngr.eventsIn <- evt
				case stateEvent:
					log.Print(" -StateEvent")
					evt := core.StateEvent{}
					rlp.DecodeBytes(msg.Data, &evt)
					if err != nil {
						log.Print(err)
						continue
					}
					mngr.eventsIn <- evt
				case addValidatorEvent:
					log.Print(" -AddValidatorEvent")
					evt := core.AddValidatorEvent{}
					rlp.DecodeBytes(msg.Data, &evt)
					if err != nil {
						log.Print(err)
						continue
					}
					mngr.eventsIn <- evt
				case customEvent:
					log.Print(" -customEvent")
					mngr.eventsCustom <- msg.Data
				}
			case gossipnet.ListenEvent:
				log.Print("ListenEvent")
			case gossipnet.CloseEvent:
				log.Print("CloseEvent")
				close(mngr.eventsIn)
				break
			}
		}
	}()

	// Dispatch IBFT events to the network
	go func() {
		for event := range mngr.eventsOut {
			switch ev := event.(type) {
			case core.MessageEvent:
				mngr.broadcast(ev.Payload, messageEvent)
			case core.AddValidatorEvent:
				evBytes, err := rlp.EncodeToBytes(ev)
				if err != nil {
					log.Print(err)
					return
				}
				mngr.broadcast(evBytes, addValidatorEvent)
			case core.StateEvent:
				log.Print("encode view", ev.View)
				evBytes, err := rlp.EncodeToBytes(ev)
				if err != nil {
					log.Print(err)
					return
				}
				mngr.broadcast(evBytes, stateEvent)
			}

		}
	}()
}

// Broadcast implements network.Manager.Broadcast. It will tag the payload
// forward it to the network node
func (mngr Manager) broadcast(payload []byte, msgType uint) (err error) {
	data, err := rlp.EncodeToBytes(networkMessage{
		Type: msgType,
		Data: payload,
	})
	if err != nil {
		return
	}
	mngr.node.Broadcast(data)
	return
}
