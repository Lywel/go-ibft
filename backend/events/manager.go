package events

import (
	"flag"
	"io/ioutil"
	"time"

	"bitbucket.org/ventureslash/go-gossipnet"
	"bitbucket.org/ventureslash/go-ibft"
	"bitbucket.org/ventureslash/go-ibft/core"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/logger"
)

var verbose = flag.Bool("verbose-manager", false, "print manager info level logs")

// Manager handles data from the network
type Manager struct {
	node           *gossipnet.Node
	eventsOut      chan core.Event
	eventsIn       chan core.Event
	eventsCustom   chan core.CustomEvent
	nodesToConnect int
	debug          *logger.Logger
}

type networkMessage struct {
	Type      uint
	Data      []byte
	Timestamp uint64
}

const (
	messageEvent uint = iota
	requestEvent
	backlogEvent
	joinEvent
	stateEvent
	addValidatorEvent
	customEvent
	validatorSetEvent
)

// New returns a new network.Manager
func New(node *gossipnet.Node, eventsIn, eventsOut chan core.Event, eventsCustom chan core.CustomEvent, nodesToConnect int) Manager {
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
		Type:      joinEvent,
		Data:      addrBytes[:],
		Timestamp: uint64(time.Now().Unix()),
	})
	if err != nil {
		mngr.debug.Warningf("encode error: %v", err)
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
				mngr.debug.Infof("Parsing %d bytes", len(ev.Data))
				var msg networkMessage
				err := rlp.DecodeBytes(ev.Data, &msg)
				if err != nil {
					mngr.debug.Warningf("Error parsing msg: %v", err)
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
					if evt.Address == addr {
						mngr.debug.Errorf("Don't add someone using your identity. It might even be yourself: disconnecting...")
						err := mngr.node.RemovePeer(evt.NetworkAddr)
						if err != nil {
							mngr.debug.Warningf("RemovePeer failed: %v", err)
						}
						continue
					}

					// If it's not youself => forwrd to POA
					mngr.eventsIn <- evt
					mngr.eventsCustom <- core.CustomEvent{
						Type: ibft.TypeJoinEvent,
						Msg:  msg.Data,
					}

				case requestEvent:
					mngr.debug.Info(" -RequestEvent")
					evt := core.EncodedRequestEvent{}
					rlp.DecodeBytes(msg.Data, &evt)
					if err != nil {
						mngr.debug.Warning(err)
						continue
					}

					mngr.eventsIn <- evt

				case validatorSetEvent:
					mngr.debug.Info(" -ValidatorSetEvent")
					mngr.eventsCustom <- core.CustomEvent{
						Type: ibft.TypeValidatorSetEvent,
						Msg:  msg.Data,
					}
				case customEvent:
					mngr.debug.Info(" -customEvent")
					mngr.eventsCustom <- core.CustomEvent{
						Type: ibft.TypeCustomEvents,
						Msg:  msg.Data,
					}
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
				mngr.debug.Infof("Broadcasting core.MessageEvent")
				mngr.broadcast(ev.Payload, messageEvent)
			case core.AddValidatorEvent:
				mngr.debug.Infof("Broadcasting core.AddValidatorEvent: %s", ev)
				evBytes, err := rlp.EncodeToBytes(ev)
				if err != nil {
					mngr.debug.Warning(err)
					return
				}
				mngr.broadcast(evBytes, addValidatorEvent)

			case core.EncodedRequestEvent:
				mngr.debug.Infof("Broadcasting core.EncodedRequestEvent: %s", ev)
				evBytes, err := rlp.EncodeToBytes(ev)
				if err != nil {
					mngr.debug.Warning(err)
					return
				}
				mngr.broadcast(evBytes, requestEvent)
			case core.ValidatorSetEvent:
				mngr.debug.Infof("Broadcasting core.ValidatorSetEvent: %s", ev)
				evBytes, err := rlp.EncodeToBytes(ev)
				if err != nil {
					mngr.debug.Warning(err)
					return
				}
				mngr.broadcast(evBytes, validatorSetEvent)
			}
		}
	}()
}

// Broadcast implements network.Manager.Broadcast. It will tag the payload
// forward it to the network node
func (mngr Manager) broadcast(payload []byte, msgType uint) (err error) {
	data, err := rlp.EncodeToBytes(networkMessage{
		Type:      msgType,
		Data:      payload,
		Timestamp: uint64(time.Now().Unix()),
	})
	if err != nil {
		return
	}
	mngr.node.Broadcast(data)
	return
}
