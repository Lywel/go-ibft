package event

const (
	eventChannelBufferSize = 256
)

type Event interface{}

type RequestEvent struct {
	Proposal Proposal
}

type MessageEvent struct {
	Payload []byte
}

type BacklogEvent struct {
	msg *Message
}

type Handler struct {
	eventChan chan Event
}

func New() Handler {
	h := Handler{
		eventChan: make(chan Event, eventChannelBufferSize),
	}
	go func() {
		for event := range h.eventChan {
			switch event.(type) {
			case RequestEvent:
			case MessageEvent:
			case BacklogEvent:
			default:
			}
		}
	}()
	return h
}

func (h Handler) Push(event Event) {
	h.eventChan <- event
}

func (h Handler) Close() {
	close(h.eventChan)
}
