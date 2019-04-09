package ntrees

import (
	"strings"

	"github.com/gokit/npkg"
	"github.com/gokit/npkg/natomic"
	"github.com/gokit/npkg/nerror"
)

//*****************************************************
// Event
//*****************************************************

// EventModder defines a function to modify a giving
// EventDescriptor.
type EventModder func(*Event)

// Event defines a giving underline signal representing an event.
type Event struct {
	TypeName string
	SourceID string
	TargetID string
	Data     interface{}
}

// EventTarget returns an new EventModder setting the
// value of the SourceID to giving value.
func EventTarget(s string) EventModder {
	return func(e *Event) {
		e.TargetID = s
	}
}

// EventSource returns an new EventModder setting the
// value of the SourceID to giving value.
func EventSource(s string) EventModder {
	return func(e *Event) {
		e.SourceID = s
	}
}

// NewEvent returns a new instance of a Event.
func NewEvent(eventName string, data interface{}, mods ...EventModder) *Event {
	var ev = &Event{
		TypeName: eventName,
		Data:     data,
	}
	for _, mod := range mods {
		mod(ev)
	}
	return ev
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

// EventDescriptorRespondHandler defines a function matching the EventDescriptorResponder.
type EventDescriptorRespondHandler func(natomic.Signal, EventDescriptor)

// RespondEvent implements the EventDescriptorResponder interface.
func (eh EventDescriptorRespondHandler) RespondEvent(ns natomic.Signal, ed EventDescriptor) {
	eh(ns, ed)
}

// EventDescriptorModder defines a function to modify a giving
// EventDescriptor.
type EventDescriptorModder func(*EventDescriptor)

// PreventDefault returns a EventModder that sets true
// to EventDescriptor.PreventDefault flag..
func PreventDefault() EventDescriptorModder {
	return func(e *EventDescriptor) {
		e.PreventDefault = true
	}
}

// StopPropagation returns a EventModder that sets true
// to EventDescriptor.StopPropagation flag..
func StopPropagation() EventDescriptorModder {
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
}

// NewEventDescriptor returns a new instance of an EventDescriptor.
func NewEventDescriptor(event string, responder interface{}, mods ...EventDescriptorModder) *EventDescriptor {
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

const eventName = "events"

// Key returns giving key or name of attribute.
func (ed *EventDescriptor) Key() string {
	return eventName
}

// Value returns giving value of attribute.
func (ed *EventDescriptor) Value() interface{} {
	return ed.Name
}

// Text returns giving value of attribute as text.
func (ed *EventDescriptor) Text() string {
	return ed.Name
}

// Namespace returns the underline namespace format for giving EventDescriptor.
//
// Namespace uses a suffix format where the PreventDefault and StopPropagation are
// arrange in the 00 order, where the first is PreventDefault and the second StopPropagation.
//
// A value of 1 means to turn it on and 0 means it's off and should be ignored.
func (ed *EventDescriptor) Namespace() string {
	var bits = "-"

	// set flag for prevent-default behaviour.
	if ed.PreventDefault {
		bits += "1"
	} else {
		bits += "0"
	}

	// set flag for stop-propagation.
	if ed.StopPropagation {
		bits += "1"
	} else {
		bits += "0"
	}

	return ed.Name + bits
}

// Contains returns true/false if giving string is the same
// name as giving event.
func (ed *EventDescriptor) Contains(s string) error {
	if ed.Name != s {
		return nerror.New("value is not the same as event name")
	}
	return nil
}

// Match returns true/false if giving attributes matched.
func (ed *EventDescriptor) Match(other Attr) bool {
	if other.Key() != eventName {
		return false
	}
	if other.Text() != ed.Name {
		return false
	}
	return true
}

// Mount implements the RenderMount interface.
func (ed *EventDescriptor) Mount(n *Node) error {
	n.Events.Add(*ed)
	return nil
}

// Unmount implements the RenderMount interface.
func (ed *EventDescriptor) Unmount(n *Node) {
	n.Events.Remove(*ed)
}

// EncodeObject implements encoding using the npkg.EncodableObject interface.
func (ed *EventDescriptor) EncodeObject(enc npkg.ObjectEncoder) error {
	var err error
	if err = enc.String("name", ed.Name); err != nil {
		return err
	}
	if err = enc.Bool("preventDefault", ed.PreventDefault); err != nil {
		return err
	}
	if err = enc.Bool("stopPropagation", ed.StopPropagation); err != nil {
		return err
	}
	return nil
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

// Len returns the underline length of events in map.
func (na *EventHashList) Len() int {
	return len(na.nodes)
}

// EncodeList encodes underline events into provided list encoder.
func (na *EventHashList) EncodeList(enc npkg.ListEncoder) error {
	for _, events := range na.nodes {
		if len(events) == 0 {
			continue
		}
		if err := enc.AddObject(&events[0]); err != nil {
			return nerror.WrapOnly(err)
		}
	}
	return nil
}

// Respond delivers giving event to all descriptors of events within hash.
func (na *EventHashList) Respond(s natomic.Signal) {
	if na.nodes == nil || len(na.nodes) == 0 {
		return
	}

	var name = strings.ToLower(s.Type())
	var descSet = na.nodes[name]

	for _, desc := range descSet {
		// if we are expected to prevent default,
		// then wrap it before sending signal.
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

// EncodeEvents encodes all giving event within provided event hash list.
func (na *EventHashList) EncodeEvents(encoder *strings.Builder) error {
	if na.nodes != nil {
		var count int
		for _, events := range na.nodes {
			if len(events) == 0 {
				continue
			}
			if count > 0 {
				if _, err := encoder.WriteString(spacer); err != nil {
					return err
				}
			}
			if _, err := encoder.WriteString(events[0].Namespace()); err != nil {
				return err
			}
			count++
		}
	}
	return nil
}

// Add adds giving node into giving list if it has
// giving attribute value.
func (na *EventHashList) Add(event EventDescriptor) {
	if na.nodes == nil {
		na.nodes = map[string][]EventDescriptor{}
	}

	var name = strings.ToLower(event.Name)
	na.nodes[event.Name] = append(na.nodes[name], event)
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

	var name = strings.ToLower(event.Name)
	var set = na.nodes[name]
	for index, desc := range set {
		if desc.SignalResponder != nil && desc.SignalResponder == event.SignalResponder {
			set = append(set[:index], set[index+1:]...)
			na.nodes[name] = set
			return
		}
		if desc.EventResponder != nil && desc.EventResponder == event.EventResponder {
			set = append(set[:index], set[index+1:]...)
			na.nodes[name] = set
			return
		}
	}

	if len(set) == 0 {
		delete(na.nodes, name)
	}
}

// RemoveSignalResponder removes giving event descriptor for giving  responder.
func (na *EventHashList) RemoveSignalResponder(event string, r natomic.SignalResponder) {
	if na.nodes == nil {
		na.nodes = map[string][]EventDescriptor{}
	}

	var set = na.nodes[event]
	event = strings.ToLower(event)
	for index, desc := range set {
		if desc.SignalResponder == r {
			set = append(set[:index], set[index+1:]...)
			na.nodes[event] = set
			return
		}
	}

	if len(set) == 0 {
		delete(na.nodes, event)
	}
}

// RemoveEventResponder removes giving event descriptor for giving  responder.
func (na *EventHashList) RemoveEventResponder(event string, r EventDescriptorResponder) {
	if na.nodes == nil {
		na.nodes = map[string][]EventDescriptor{}
	}

	var set = na.nodes[event]
	event = strings.ToLower(event)
	for index, desc := range set {
		if desc.EventResponder == r {
			set = append(set[:index], set[index+1:]...)
			na.nodes[event] = set
			return
		}
	}

	if len(set) == 0 {
		delete(na.nodes, event)
	}
}
