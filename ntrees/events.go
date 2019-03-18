package ntrees

import (
	"github.com/gokit/npkg/natomic"
)

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

// EventDescriptor defines a type representing a event descriptor with
// associated response.
type EventDescriptor struct {
	Name            string
	Immediate       bool
	PreventDefault  bool
	StopPropagation bool
	Handler         natomic.SignalResponder
}

// Respond delivers giving event signal to all handlers.
func (er *EventDescriptor) Respond(s natomic.Signal) {
	if er.Handler == nil {
		return
	}
	er.Handler.Respond(s)
}

// EventHashList implements the a set list for Nodes using
// their Node.RefID() value as unique keys.
type EventHashList struct {
	nodes map[string][]EventDescriptor
}

// NewEventHashList returns a new instance EventHashList.
func NewEventHashList() *EventHashList {
	return &EventHashList{
		nodes: map[string][]EventDescriptor{},
	}
}

// Reset resets the internal hashmap used for storing
// nodes. There by removing all registered nodes.
func (na *EventHashList) Reset() {
	na.nodes = map[string][]EventDescriptor{}
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
func (na *EventHashList) Add(event string, preventDef bool, stopPropagate bool, immediate bool, responder natomic.SignalResponder) {
	if na.nodes == nil {
		na.nodes = map[string][]EventDescriptor{}
	}

	var desc EventDescriptor
	desc.Name = event
	desc.Handler = responder
	desc.Immediate = immediate
	desc.PreventDefault = preventDef
	desc.StopPropagation = stopPropagate
	na.nodes[event] = append(na.nodes[event], desc)
}

// RemoveResponder removes giving event descriptor for giving  responder.
func (na *EventHashList) RemoveResponder(event string, r natomic.SignalResponder) {
	if na.nodes == nil {
		na.nodes = map[string][]EventDescriptor{}
	}

	var set = na.nodes[event]
	for index, desc := range set {
		if desc.Handler == r {
			set = append(set[:index], set[index+1:]...)
			na.nodes[event] = set
			return
		}
	}
}

// Remove removes giving node into giving list if it has
// giving attribute value.
func (na *EventHashList) Remove(event string) {
	if na.nodes == nil {
		na.nodes = map[string][]EventDescriptor{}
	}
	delete(na.nodes, event)
}
