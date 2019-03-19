package ntrees

import (
	"github.com/gokit/npkg/natomic"
)

//*****************************************************
// Event
//*****************************************************

// Event defines a giving underline signal representing an event.
type Event struct {
	TypeName string
	SourceID string
	TargetID string
	Data     interface{}
}

// Type returns the underline typename of the event.
//
// It implements the natomic.Signal interface.
func (e Event) Type() string {
	return e.TypeName
}

// Source returns the source of the giving event.
//
// It implements the natomic.Signal interface.
func (e Event) Source() string {
	return e.SourceID
}

// Target returns the target of the giving event.
//
// It implements the natomic.Signal interface.
func (e Event) Target() string {
	return e.TargetID
}

//*****************************************************
// EventDescriptor
//*****************************************************

// parentEvent defines a string type for defining a
// parent event listen list.
type parentEvent string

// Mount adds giving event into parent event listen list.
func (p parentEvent) Mount(n *Node) error {
	n.crossEvents[string(p)] = true
	if n.parent != nil {
		n.parent.addChildEventListener(string(p), n)
	}
	return nil
}

// OnParentEvent returns a event which will be added to the
// parent of the node it is added to, this will have the
// node's parent nodify the node on the occurrence of that
// event.
func OnParentEvent(eventName string) Mounter {
	return parentEvent(eventName)
}

// EventPreventer wraps a giving event signal returning
// default prevention.
type EventPreventer struct {
	natomic.Signal
	PreventDefault bool
}

// EventDescriptorResponder defines an interface which responds to a
// signal with giving EventDescriptor.
type EventDescriptorResponder interface {
	RespondEvent(natomic.Signal, EventDescriptor)
}

// EventModder defines a function to modify a giving
// EventDescriptor.
type EventModder func(*EventDescriptor)

// PreventDefault returns a EventModder that sets true
// to EventDescriptor.PreventDefault flag..
func PreventDefault() EventModder {
	return func(e *EventDescriptor) {
		e.PreventDefault = true
	}
}

// StopPropagation returns a EventModder that sets true
// to EventDescriptor.StopPropagation flag..
func StopPropagation() EventModder {
	return func(e *EventDescriptor) {
		e.StopPropagation = true
	}
}

// EventDescriptor defines a type representing a event descriptor with
// associated response.
type EventDescriptor struct {
	Name            string
	PreventDefault  bool
	StopPropagation bool
	SignalResponder natomic.SignalResponder
	EventResponder  EventDescriptorResponder

	// rootCrisscross tells us this is an event
	// from a lower node to it's parent, we should
	// avoid doing upward propagation.
	rootCrisscross bool
}

// NewEventDescriptor returns a new instance of an EventDescriptor.
func NewEventDescriptor(event string, responder interface{}, mods ...EventModder) *EventDescriptor {
	var desc EventDescriptor
	desc.Name = event

	for _, mod := range mods {
		mod(&desc)
	}

	if responder != nil {
		switch tm := responder.(type) {
		case EventDescriptorResponder:
			desc.EventResponder = tm
		case natomic.SignalResponder:
			desc.SignalResponder = tm
		default:
			panic("Unknown type for event handler")
		}
	}
	return &desc
}

// Mount implements the RenderMount interface.
func (ed *EventDescriptor) Mount(n *Node) {
	n.Events.Add(*ed)
}

// Remove implements the RenderMount interface.
func (ed *EventDescriptor) Unmount(n *Node) {
	n.Events.Remove(*ed)
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

// Respond delivers giving event to all descriptors of events within hash.
func (na *EventHashList) Respond(s natomic.Signal) {
	if na.nodes == nil || len(na.nodes) == 0 {
		return
	}

	var descSet = na.nodes[s.Type()]
	for _, desc := range descSet {
		// if we are expected to prevent default, then wrap it before sending
		// signal.
		if em, ok := s.(*EventPreventer); ok {
			if desc.PreventDefault {
				em.PreventDefault = true
			} else {
				em.PreventDefault = false
			}
		} else if desc.PreventDefault {
			s = &EventPreventer{
				Signal:         s,
				PreventDefault: desc.PreventDefault,
			}
		}

		if desc.EventResponder != nil {
			desc.EventResponder.RespondEvent(s, desc)
			continue
		}
		if desc.SignalResponder != nil {
			desc.SignalResponder.Respond(s)
			continue
		}
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
func (na *EventHashList) Add(event EventDescriptor) {
	if na.nodes == nil {
		na.nodes = map[string][]EventDescriptor{}
	}

	na.nodes[event.Name] = append(na.nodes[event.Name], event)
}

// RemoveAll removes giving node in list if it has
// giving attribute value.
func (na *EventHashList) RemoveAll(event string) {
	if na.nodes == nil {
		na.nodes = map[string][]EventDescriptor{}
	}
	delete(na.nodes, event)
}

// Remove removes giving node in list if it has
// giving handler.
func (na *EventHashList) Remove(event EventDescriptor) {
	if na.nodes == nil {
		na.nodes = map[string][]EventDescriptor{}
	}

	var set = na.nodes[event.Name]
	for index, desc := range set {
		if desc.SignalResponder != nil && desc.SignalResponder == event.SignalResponder {
			set = append(set[:index], set[index+1:]...)
			na.nodes[event.Name] = set
			return
		}
		if desc.EventResponder != nil && desc.EventResponder == event.EventResponder {
			set = append(set[:index], set[index+1:]...)
			na.nodes[event.Name] = set
			return
		}
	}
}

// RemoveSignalResponder removes giving event descriptor for giving  responder.
func (na *EventHashList) RemoveSignalResponder(event string, r natomic.SignalResponder) {
	if na.nodes == nil {
		na.nodes = map[string][]EventDescriptor{}
	}

	var set = na.nodes[event]
	for index, desc := range set {
		if desc.SignalResponder == r {
			set = append(set[:index], set[index+1:]...)
			na.nodes[event] = set
			return
		}
	}
}

// RemoveEventResponder removes giving event descriptor for giving  responder.
func (na *EventHashList) RemoveEventResponder(event string, r EventDescriptorResponder) {
	if na.nodes == nil {
		na.nodes = map[string][]EventDescriptor{}
	}

	var set = na.nodes[event]
	for index, desc := range set {
		if desc.EventResponder == r {
			set = append(set[:index], set[index+1:]...)
			na.nodes[event] = set
			return
		}
	}
}
