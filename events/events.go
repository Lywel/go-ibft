package events

const (
	eventChannelBufferSize = 256
)

// Event is an ibft event
type Event interface{}

// Handler is an event handler using a channel
type Handler struct {
	eventChan chan Event
}

// New returns a new events.Handler
func New() Handler {
	return Handler{
		eventChan: make(chan Event, eventChannelBufferSize),
	}
}

// Push implements Handler.Push. It adds an event in the processing queue.
func (h Handler) Push(event Event) {
	h.eventChan <- event
}

// Close implements Handle.Close. It closes the event handler channel
func (h Handler) Close() {
	close(h.eventChan)
}

// EventChan implements Handle.EventChan
func (h Handler) EventChan() <-chan Event {
	return h.eventChan
}
