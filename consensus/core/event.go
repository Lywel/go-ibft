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
	EventChan chan Event
}

func New() Handler {
	return Handler{
		eventChan: make(chan Event, eventChannelBufferSize),
	}
}

func (h Handler) Push(event Event) {
	h.eventChan <- event
}

func (h Handler) Close() {
	close(h.eventChan)
}
