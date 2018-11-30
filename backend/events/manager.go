package events

import (
	"bitbucket.org/ventureslash/go-gossipnet"
	"bitbucket.org/ventureslash/go-ibft"
	"bitbucket.org/ventureslash/go-ibft/core"
	"flag"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/logger"
	"io/ioutil"
)

var verbose = flag.Bool("verbose-manager", false, "print manager info level logs")

// Manager handles data from the network
type Manager struct {
	node           *gossipnet.Node
	eventsOut      chan core.Event
	eventsIn       chan core.Event
	eventsCustom   chan []byte
	nodesToConnect int
	debug          *logger.Logger
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
	mngr.debug = logger.Init("Manager", *verbose, false, ioutil.Discard)
	mngr.debug.Infof("Start")

	addrBytes := addr.GetBytes()
	joinBytes, err := rlp.EncodeToBytes(networkMessage{
		Type: joinEvent,
		Data: addrBytes[:],
	})
	if err != nil {
		mngr.debug.Warningf("encode error: ", err)
	}

	// Dispatch network events to IBFT
	go func() {
		for event := range mngr.node.EventChan() {
			switch ev := event.(type) {
			case gossipnet.ConnOpenEvent:
				mngr.debug.Info("ConnOpenEvent")
				// TODO: dont gossip to everyone, just the new connection
				if mngr.nodesToConnect > 0 {
					mngr.node.Gossip(joinBytes)
					mngr.nodesToConnect--
				}
			case gossipnet.ConnCloseEvent:
				mngr.debug.Info("ConnCloseEvent")
			case gossipnet.DataEvent:
				mngr.debug.Info("DataEvent")
				var msg networkMessage
				err := rlp.DecodeBytes(ev.Data, &msg)
				if err != nil {
					mngr.debug.Warningf("Error parsing msg:", string(ev.Data))
					continue
				}
				switch msg.Type {
				case messageEvent:
					mngr.debug.Info(" -MsgEvent")
					mngr.eventsIn <- core.MessageEvent{
						Payload: msg.Data,
					}
				case joinEvent:
					mngr.debug.Info(" -JoinEvent")
					evt := core.JoinEvent{
						NetworkAddr: ev.Addr,
					}
					evt.Address.FromBytes(msg.Data)
					if err != nil {
						mngr.debug.Warning(err)
						continue
					}
					mngr.eventsIn <- evt
				case stateEvent:
					mngr.debug.Info(" -StateEvent")
					evt := core.StateEvent{}
					rlp.DecodeBytes(msg.Data, &evt)
					if err != nil {
						mngr.debug.Warning(err)
						continue
					}
					mngr.eventsIn <- evt
				case requestEvent:
					mngr.debug.Info(" -RequestEvent")
					evt := core.EncodedRequestEvent{}
					rlp.DecodeBytes(msg.Data, &evt)
					if err != nil {
						mngr.debug.Warning(err)
						continue
					}

					mngr.eventsIn <- evt
				case addValidatorEvent:
					mngr.debug.Info(" -AddValidatorEvent")
					evt := core.AddValidatorEvent{}
					rlp.DecodeBytes(msg.Data, &evt)
					if err != nil {
						mngr.debug.Warning(err)
						continue
					}
					mngr.eventsIn <- evt
				case customEvent:
					mngr.debug.Info(" -customEvent")
					mngr.eventsCustom <- msg.Data
				}
			case gossipnet.ListenEvent:
				mngr.debug.Info("ListenEvent")
			case gossipnet.CloseEvent:
				mngr.debug.Info("CloseEvent")
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
					mngr.debug.Warning(err)
					return
				}
				mngr.broadcast(evBytes, addValidatorEvent)
			case core.StateEvent:
				mngr.debug.Infof("encoded view: %s", ev.View)
				evBytes, err := rlp.EncodeToBytes(ev)
				if err != nil {
					mngr.debug.Warning(err)
					return
				}
				mngr.broadcast(evBytes, stateEvent)
			case core.EncodedRequestEvent:
				evBytes, err := rlp.EncodeToBytes(ev)
				if err != nil {
					mngr.debug.Warning(err)
					return
				}
				mngr.broadcast(evBytes, requestEvent)

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
