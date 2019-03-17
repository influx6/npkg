package ntrees

var _ Attrs = (*AttrList)(nil)

// AttrEncoder defines an interface which provides means of encoding
// attribute key-value pairs and possible sub-attributes using provided
// type functions.
type AttrEncoder interface {
	Attr(string, AttrEncodable) error

	Int(string, int) error
	Int64(string, int64) error
	Bool(string, string) error
	String(string, string) error
	Float64(string, float64) error
}

// AttrEncodable exposes a interface which provides method for encoder attributes
// using provided encoder.
type AttrEncodable interface {
	EncodeAttr(encoder AttrEncoder) error
}

// Attr defines a series of method representing a Attribute.
type Attr interface {
	AttrEncodable

	// Key returns the key for the attribute.
	Key() string

	// Text returns a textual representation of giving attribute value.
	Text() string

	// Value return the value of giving attribute as an interface.
	Value() interface{}

	// Match must match against provided attribute validating if
	// it is equal both in type and key with value.
	Match(Attr) bool

	// Contains should return true/false if giving attribute
	// contains provided value.
	Contains(value string) bool
}

// IterableAttr defines an interface that exposes a method to
// iterate through all possible attribute values or value.
type IterableAttr interface {
	Attr

	Each(func(string) bool)
}

// Attrs exposes a interface defining a giving attribute host
// which provides method for accessing all attributes.
type Attrs interface {
	AttrEncodable

	// Has should return true/false if giving Attrs has giving key.
	Has(key string) bool

	// MatchAttrs returns true/false if provided Attrs match each other.
	MatchAttrs(Attrs) bool

	// Each should handle the need of iterating through all
	// values of a key.
	Each(fx func(Attr) bool)

	// Attr should return the Attr for giving key of giving type.
	Attr(key string) (Attr, bool)

	// Match must match against provided key and string value returning
	// true/false if value matches internal representation of key value/values.
	Match(key string, value string) bool
}

// AttrList implements Attrs interface.
type AttrList []Attr

// EncodeAttr encodes all attributes within it's list with
// provided encoder.
func (l AttrList) EncodeAttr(encoder AttrEncoder) error {
	for _, item := range l {
		if err := item.EncodeAttr(encoder); err != nil {
			return err
		}
	}
	return nil
}

// Each runs through all attributes within list.
func (l AttrList) Each(fx func(Attr) bool) {
	for _, item := range l {
		if fx(item) {
			continue
		}
		break
	}
}

// Has returns true/false if giving list has giving key.
func (l AttrList) Has(key string) bool {
	for _, item := range l {
		if item.Key() == key {
			return true
		}
	}
	return false
}

// Attribute returns giving Attribute with provided key.
func (l AttrList) Attr(key string) (Attr, bool) {
	for _, item := range l {
		if item.Key() == key {
			return item, true
		}
	}
	return nil, false
}

// MatchAttrs returns true/false if giving attrs match.
func (l AttrList) MatchAttrs(attrs Attrs) bool {
	for _, item := range l {
		if attr, ok := attrs.Attr(item.Key()); ok {
			if item.Match(attr) {
				continue
			}
		}
		return false
	}
	return true
}

// Match validates if giving value matches attribute's value.
func (l AttrList) Match(key string, value string) bool {
	if attr, ok := l.Attr(key); ok {
		return attr.Contains(value)
	}
	return false
}
