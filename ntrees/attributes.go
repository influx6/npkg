package ntrees

import (
	"io"
	"strconv"
	"strings"
	"sync"
)

var _ Attrs = (*AttrList)(nil)
var stringPool = sync.Pool{
	New: func() interface{} {
		var content strings.Builder
		return &content
	},
}

// AttrEncoder defines an interface which provides means of encoding
// attribute key-value pairs and possible sub-attributes using provided
// type functions.
type AttrEncoder interface {
	Attr(string, AttrEncodable) error

	Int(string, int) error
	Float(string, float64) error
	List(string, ...string) error
	QuotedString(string, string) error
	UnquotedString(string, string) error
}

// AttrEncodable exposes a interface which provides method for encoder attributes
// using provided encoder.
type AttrEncodable interface {
	EncodeAttr(encoder AttrEncoder) error
}

// DOMAttrEncoder implements a not to optimized AttrEncoder interface.
type DOMAttrEncoder struct {
	Key     string
	Content *strings.Builder
}

// NewDOMAttrEncoder returns a new DOMAttrEncoder.
func NewDOMAttrEncoder(key string) *DOMAttrEncoder {
	var content strings.Builder
	return &DOMAttrEncoder{
		Key:     key,
		Content: &content,
	}
}

// DOMAttrEncoderWith returns a new DOMAttrEncoder.
func DOMAttrEncoderWith(key string, content *strings.Builder) *DOMAttrEncoder {
	return &DOMAttrEncoder{
		Key:     key,
		Content: content,
	}
}

// String returns the encoded attribute list of elements.
func (dm *DOMAttrEncoder) WriteTo(w io.Writer) (int64, error) {
	var n, err = w.Write([]byte(dm.Content.String()))
	return int64(n), err
}

// String returns the encoded attribute list of elements.
func (dm *DOMAttrEncoder) String() string {
	return dm.Content.String()
}

// Attr implements encoding of multi-attribute based values.
func (dm *DOMAttrEncoder) WithAttr(key string, fn func(encoder AttrEncoder) error) error {
	var err error
	var content = stringPool.Get().(*strings.Builder)
	defer stringPool.Put(content)

	content.Reset()

	if dm.Key != "" {
		key = dm.Key + "." + key
	}

	var dmer = &DOMAttrEncoder{
		Key:     key,
		Content: content,
	}
	if err = fn(dmer); err != nil {
		return err
	}

	if dm.Content.Len() > 0 {
		if _, err = dm.Content.WriteString(" "); err != nil {
			return err
		}
	}
	if _, err = dm.Content.WriteString(content.String()); err != nil {
		return err
	}
	return nil
}

// Attr implements encoding of multi-attribute based values.
func (dm *DOMAttrEncoder) Attr(key string, attrs AttrEncodable) error {
	var err error
	var content = stringPool.Get().(*strings.Builder)
	defer stringPool.Put(content)

	content.Reset()

	if dm.Key != "" {
		key = dm.Key + "." + key
	}

	var dmer = &DOMAttrEncoder{
		Key:     key,
		Content: content,
	}

	if err = attrs.EncodeAttr(dmer); err != nil {
		return err
	}

	if dm.Content.Len() > 0 {
		if _, err = dm.Content.WriteString(" "); err != nil {
			return err
		}
	}
	if _, err = dm.Content.WriteString(content.String()); err != nil {
		return err
	}
	return nil
}

// Attr encodes giving list of string values for string key.
func (dm *DOMAttrEncoder) List(key string, set ...string) error {
	var err error
	var content = stringPool.Get().(*strings.Builder)
	defer stringPool.Put(content)

	content.Reset()

	for i, s := 0, len(set); i < s; i++ {
		if _, err = content.WriteString(set[i]); err != nil {
			return err
		}
		if i < (s - 1) {
			if _, err = content.WriteString(","); err != nil {
				return err
			}
		}
	}

	if dm.Content.Len() > 0 {
		if _, err = dm.Content.WriteString(" "); err != nil {
			return err
		}
	}
	if dm.Key != "" {
		if _, err = dm.Content.WriteString(dm.Key + "."); err != nil {
			return err
		}
	}
	if _, err = dm.Content.WriteString(key); err != nil {
		return err
	}
	if _, err = dm.Content.WriteString("=\""); err != nil {
		return err
	}
	if _, err = dm.Content.WriteString(content.String()); err != nil {
		return err
	}
	if _, err = dm.Content.WriteString("\""); err != nil {
		return err
	}
	return nil
}

// QuotedString encodes giving string value for string key.
func (dm *DOMAttrEncoder) QuotedString(key string, val string) error {
	var err error
	if dm.Content.Len() > 0 {
		if _, err = dm.Content.WriteString(" "); err != nil {
			return err
		}
	}
	if dm.Key != "" {
		if _, err = dm.Content.WriteString(dm.Key + "."); err != nil {
			return err
		}
	}
	if _, err = dm.Content.WriteString(key); err != nil {
		return err
	}
	if _, err = dm.Content.WriteString("=\""); err != nil {
		return err
	}
	if _, err = dm.Content.WriteString(val); err != nil {
		return err
	}
	if _, err = dm.Content.WriteString("\""); err != nil {
		return err
	}
	return nil
}

// UnquotedString encodes giving string value for string key.
func (dm *DOMAttrEncoder) UnquotedString(key string, val string) error {
	var err error
	if dm.Content.Len() > 0 {
		if _, err = dm.Content.WriteString(" "); err != nil {
			return err
		}
	}
	if dm.Key != "" {
		if _, err = dm.Content.WriteString(dm.Key + "."); err != nil {
			return err
		}
	}
	if _, err = dm.Content.WriteString(key); err != nil {
		return err
	}
	if _, err = dm.Content.WriteString("="); err != nil {
		return err
	}
	if _, err = dm.Content.WriteString(val); err != nil {
		return err
	}
	return nil
}

// Float encodes giving int value for string key.
func (dm *DOMAttrEncoder) Float(key string, val float64) error {
	var content [8]byte
	var appended = strconv.AppendFloat(content[:0], val, 'f', -1, 64)
	return dm.UnquotedString(key, string(appended))
}

// Int encodes giving int value for string key.
func (dm *DOMAttrEncoder) Int(key string, val int) error {
	var content [8]byte
	var appended = strconv.AppendInt(content[:0], int64(val), 10)
	return dm.UnquotedString(key, string(appended))
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

// IntAttr implements the Attr interface for a string key-value pair.
type IntAttr struct {
	Name string
	Val  int
}

// NewIntAttr returns a new instance of a IntAttr.
func NewIntAttr(n string, v int) IntAttr {
	return IntAttr{Name: n, Val: v}
}

// Key returns giving key or name of attribute.
func (s IntAttr) Key() string {
	return s.Name
}

// Value returns giving value of attribute.
func (s IntAttr) Value() interface{} {
	return s.Val
}

// Text returns giving value of attribute as text.
func (s IntAttr) Text() string {
	return strconv.Itoa(s.Val)
}

// Mount implements the Mounter interface.
func (s IntAttr) Mount(parent *Node) error {
	parent.Attrs.Add(s)
	return nil
}

// Contains returns true/false if provided value is contained in attr.
//
// Since we are dealing with a number, we attempt to convert the provided
// value into a number and match else return false.
func (s IntAttr) Contains(other string) bool {
	var vm, err = strconv.Atoi(other)
	if err != nil {
		return false
	}
	return vm == s.Val
}

// EncodeAttr implements the AttrEncodable interface.
func (s IntAttr) EncodeAttr(encoder AttrEncoder) error {
	return encoder.Int(s.Name, s.Val)
}

// Match returns true/false if giving attributes matches.
func (s IntAttr) Match(other Attr) bool {
	if other.Key() != s.Name {
		return false
	}
	if other.Text() != s.Text() {
		return false
	}
	return true
}

// StringListAttr implements the Attr interface for a string key-value pair.
type StringListAttr struct {
	Join string
	Name string
	Val  []string
}

// NewStringListAttr returns a new instance of a StringListAttr
func NewStringListAttr(n string, v ...string) StringListAttr {
	return StringListAttr{Name: n, Val: v}
}

// Key returns giving key or name of attribute.
func (s StringListAttr) Key() string {
	return s.Name
}

// Value returns giving value of attribute.
func (s StringListAttr) Value() interface{} {
	return s.Val
}

// Text returns giving value of attribute as text.
func (s StringListAttr) Text() string {
	return strings.Join(s.Val, s.Join)
}

// Mount implements the Mounter interface.
func (s StringListAttr) Mount(parent *Node) error {
	parent.Attrs.Add(s)
	return nil
}

// Contains returns true/false if provided value is contained in attr.
func (s StringListAttr) Contains(other string) bool {
	for _, cn := range s.Val {
		if cn == other {
			return true
		}
	}
	return false
}

// EncodeAttr implements the AttrEncodable interface.
func (s StringListAttr) EncodeAttr(encoder AttrEncoder) error {
	return encoder.List(s.Name, s.Val...)
}

// Match returns true/false if giving attributes matches.
func (s StringListAttr) Match(other Attr) bool {
	if other.Key() != s.Name {
		return false
	}
	if other.Text() != s.Text() {
		return false
	}
	return true
}

// StringAttr implements the Attr interface for a string key-value pair.
type StringAttr struct {
	Name string
	Val  string
}

// NewStringAttr returns a new instance of a StringAttr
func NewStringAttr(n string, v string) StringAttr {
	return StringAttr{Name: n, Val: v}
}

// NamespaceAttr returns a new StringAttr with Name set to "namespace".
func NamespaceAttr(v string) StringAttr {
	return StringAttr{Name: "namespace", Val: v}
}

// Key returns giving key or name of attribute.
func (s StringAttr) Key() string {
	return s.Name
}

// Value returns giving value of attribute.
func (s StringAttr) Value() interface{} {
	return s.Val
}

// Text returns giving value of attribute as text.
func (s StringAttr) Text() string {
	return s.Val
}

// Mount implements the Mounter interface.
func (s StringAttr) Mount(parent *Node) error {
	parent.Attrs.Add(s)
	return nil
}

// Contains returns true/false if provided value is contained in attr.
func (s StringAttr) Contains(other string) bool {
	return strings.Contains(s.Val, other)
}

// EncodeAttr implements the AttrEncodable interface.
func (s StringAttr) EncodeAttr(encoder AttrEncoder) error {
	return encoder.QuotedString(s.Name, s.Val)
}

// Match returns true/false if giving attributes matches.
func (s StringAttr) Match(other Attr) bool {
	if other.Key() != s.Name {
		return false
	}
	if other.Text() != s.Val {
		return false
	}
	return true
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

// Add adds giving attribute into list.
func (l *AttrList) Add(v Attr) {
	*l = append(*l, v)
}

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

// Attr returns giving Attribute with provided key.
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
