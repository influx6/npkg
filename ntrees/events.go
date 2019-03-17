package ntrees

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
