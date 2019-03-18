package ntrees

import (
	"github.com/gokit/npkg/natomic"
)

type EventHandler func(natomic.Signal)

type EventResponder struct {
	Name            string
	Immediate       bool
	PreventDefault  bool
	StopPropagation bool
	Handlers        []EventHandler
}

// Add adds giving handlers to the list of EventResponder handlers.
func (er *EventResponder) Add(h EventHandler) {
	if capacity := cap(er.Handlers); capacity == len(er.Handlers) {
		if capacity == 0 {
			capacity = 1
		}
		var newHandlers = make([]EventHandler, capacity*2)
		var copied = copy(newHandlers, er.Handlers)
		er.Handlers = newHandlers[:copied]
	}
	er.Handlers = append(er.Handlers, h)
}

// Respond delivers giving event signal to all handlers.
func (er *EventResponder) Respond(s natomic.Signal) {
	for _, h := range er.Handlers {
		h(s)
	}
}

// EventHashList implements the a set list for Nodes using
// their Node.RefID() value as unique keys.
type EventHashList struct {
	nodes map[string]EventResponder
}

// Reset resets the internal hashmap used for storing
// nodes. There by removing all registered nodes.
func (na *EventHashList) Reset() {
	na.nodes = map[string]EventResponder{}
}

// Count returns the total content count of map
func (na *EventHashList) Count() int {
	if na.nodes == nil {
		return 0
	}
	return len(na.nodes)
}

// Add adds giving node into giving list if it has
// giving attribute value.
func (na *EventHashList) Add(event string, handler EventHandler) {
	if na.nodes == nil {
		na.nodes = map[string]EventResponder{}
	}
}

// Remove removes giving node into giving list if it has
// giving attribute value.
func (na *EventHashList) Remove(event string, handler EventHandler) {
	if na.nodes == nil {
		na.nodes = map[string]EventResponder{}
	}
}

// Event defines a giving underline signal representing an event.
type Event struct {
	TypeName string
	SourceID string
	TargetID string
}

// Type returns the underline typename of the event.
//
// It implements the natomic.Signal interface.
func (e *Event) Type() string {
	return e.TypeName
}

// Source returns the source of the giving event.
//
// It implements the natomic.Signal interface.
func (e *Event) Source() string {
	return e.SourceID
}

// Target returns the target of the giving event.
//
// It implements the natomic.Signal interface.
func (e *Event) Target() string {
	return e.TargetID
}
