package nerror

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/influx6/npkg"
	"github.com/influx6/npkg/nframes"
	"github.com/influx6/npkg/njson"
)

type IsError interface {
	IsErr(err error) bool
}

// IsAny returns true/false any of giving error matches set of error.
func IsAny(err error, set ...error) bool {
	err = UnwrapDeep(err)
	for _, item := range set {
		if isErrItem, ok := item.(IsError); ok {
			if !isErrItem.IsErr(err) {
				continue
			}
			return true
		}
		if IsSameTypeName(err, item) {
			return true
		}
		if item == err {
			return true
		}
	}
	return false
}

// ErrorStack allows you to stack multiple messages together
// which can be combined into a single error message when
// calling ErrorStack.Err().
//
// If no messages are found then no error is returned.
type ErrorStack []string

func (es *ErrorStack) Add(msg string, m ...interface{}) {
	if len(m) > 0 {
		*es = append(*es, "- "+fmt.Sprintf(msg, m...)+"\n")
		return
	}
	*es = append(*es, "- "+msg+"\n")
}

func (es *ErrorStack) Err() error {
	if len(*es) > 0 {
		return New(strings.Join(*es, ""))
	}
	return nil
}

// ErrorOption defines a giving function which receiving
// a giving error returns another modified with respect
// to it's internal intent.
type ErrorOption func(error) error

// Apply applies giving set of ErrorOptions to provided error.
func Apply(err error, ops ...ErrorOption) error {
	for _, op := range ops {
		err = op(err)
	}
	return err
}

// Frames will attempt to add call stack into provided error if
// error is a PointingError type without a stack.
func Frames(frames nframes.Frames) ErrorOption {
	return func(e error) error {
		pe := unwrapAs(e)
		pe.Frames = frames.Details()
		return pe
	}
}

// Meta adds giving header map as meta information to error.
func Meta(err error, header map[string]interface{}) ErrorOption {
	return func(e error) error {
		pe := unwrapAs(e)
		pe.Meta = header
		return pe
	}
}

// Stacked returns an error from provided message and parameter
// list if provided. It adds necessary information related
// to point of return.
func Stacked() ErrorOption {
	return func(e error) error {
		next := unwrapAs(e)
		next.Frames = nframes.GetFrameDetails(3, 32)
		return next
	}
}

// StackedBy returns an error from provided message and parameter
// list if provided. It adds necessary information related
// to point of return.
func StackedBy(n int) ErrorOption {
	return func(e error) error {
		next := unwrapAs(e)
		next.Frames = nframes.GetFrameDetails(3, 32)
		return next
	}
}

// StackWrap returns a new error which wraps existing error value if
// present and also collects current stack trace into returned error.
// It formats message accordingly with arguments from
// variadic list v.
func StackWrap(err error, message string, v ...interface{}) *PointingError {
	if len(v) != 0 {
		message = fmt.Sprintf(message, v...)
	}

	var next PointingError
	next.Parent = err
	next.Message = message
	next.Frames = nframes.GetFrameDetails(3, 32)
	return &next
}

// NewStack returns a new error which wraps existing error value if
// present and also collects current stack trace into returned error.
// It formats message accordingly with arguments from
// variadic list v.
func NewStack(message string, v ...interface{}) *PointingError {
	if len(v) != 0 {
		message = fmt.Sprintf(message, v...)
	}

	var next PointingError
	next.Message = message
	next.Frames = nframes.GetFrameDetails(3, 32)
	return &next
}

// New returns an error from provided message and parameter
// list if provided. It adds necessary information related
// to point of return.
func New(message string, v ...interface{}) *PointingError {
	if len(v) != 0 {
		message = fmt.Sprintf(message, v...)
	}

	var next PointingError
	next.Message = message
	next.Frames = nframes.GetFrameDetails(3, 32)
	return &next
}

// NewBy returns an error from provided message and parameter
// list if provided. It adds necessary information related
// to point of return.
func NewBy(n int, message string, v ...interface{}) *PointingError {
	if len(v) != 0 {
		message = fmt.Sprintf(message, v...)
	}

	var next PointingError
	next.Message = message
	next.Frames = nframes.GetFrameDetails(3, n)
	return &next
}

// WrapBy returns a new error which wraps existing error value if
// present. It formats message accordingly with arguments from
// variadic list v.
func WrapBy(n int, err error, message string, v ...interface{}) *PointingError {
	if len(v) != 0 {
		message = fmt.Sprintf(message, v...)
	}

	var next PointingError
	next.Parent = err
	next.Message = message
	next.Frames = nframes.GetFrameDetails(3, n)
	return &next
}

// WrapOnly returns a new error which wraps existing error value if
// present.
func WrapOnly(err error) *PointingError {
	if tm, ok := err.(*PointingError); ok {
		return tm
	}
	return wrapOnlyBy(err, 4, 32)
}

// Unwrap returns the underline error of giving PointingError.
func Unwrap(e error) error {
	if tm, ok := e.(*PointingError); ok {
		if tm.Parent == nil {
			return tm
		}
		return tm.Parent
	}
	return e
}

// UnwrapDeep returns the root error wrapped by all possible PointingError types.
// It attempts to retrieve the original error.
func UnwrapDeep(e error) error {
	if tm, ok := e.(*PointingError); ok {
		if tm.Parent == nil {
			return tm
		}

		return UnwrapDeep(tm.Parent)
	}

	return e
}

// Forward wraps giving error, recording where it was
// created and attaches the frames of call.
func Forward(err error) *PointingError {
	next := wrapOnly(err)
	next.Parent = err
	next.Message = err.Error()
	next.Frames = nframes.GetFrameDetails(3, 32)
	return next
}

// Wrap returns a new error which wraps existing error value if
// present. It formats message accordingly with arguments from
// variadic list v.
func Wrap(err error, message string, v ...interface{}) *PointingError {
	if len(v) != 0 {
		message = fmt.Sprintf(message, v...)
	}

	next := wrapOnly(err)
	next.Parent = err
	next.Message = message
	next.Frames = nframes.GetFrameDetails(3, 32)
	return next
}

// wrapOnly returns a new error which wraps existing error value if
// present.
func wrapOnly(err error) *PointingError {
	return wrapOnlyBy(err, 4, 32)
}

// WrapOnlyBy returns a new error which wraps existing error value if
// present.
func wrapOnlyBy(err error, depth int, stack int) *PointingError {
	var next PointingError
	next.Parent = err
	next.Frames = nframes.GetFrameDetails(depth, stack)
	return &next
}

// unwrapAs unwraps giving error to PointingError if it is else wraps
// and returns a new PointingError.
func unwrapAs(e error) *PointingError {
	if tm, ok := e.(*PointingError); ok {
		return tm
	}
	return wrapOnlyBy(e, 4, 32)
}

var _ ErrorMessage = (*PointingError)(nil)
var _ HasMessage = (*PointingError)(nil)

// PointingError defines a custom error type which points to
// both an originating point of return and a parent error if
// wrapped.
type PointingError struct {
	Message string
	Params  map[string]string
	Frames  []nframes.FrameDetail
	Meta    map[string]interface{}
	Parent  error
}

func (pe *PointingError) Add(key, value string) *PointingError {
	if pe.Params == nil {
		pe.Params = map[string]string{}
	}
	pe.Params[key] = value
	return pe
}

func (pe *PointingError) IsError(err error) bool {
	if pe.Parent == nil {
		return false
	}
	if isErrorParent, ok := pe.Parent.(IsError); ok {
		return isErrorParent.IsErr(err)
	}
	if IsSameTypeName(pe.Parent, err) {
		return true
	}
	return pe.Parent == err
}

// Error implements the error interface.
func (pe *PointingError) Error() string {
	return pe.String()
}

type HasMessage interface {
	HasMessage() bool
}

func (pe PointingError) HasMessage() bool {
	return len(pe.Message) > 0
}

type ErrorMessage interface {
	GetMessage() string
}

func (pe *PointingError) GetMessage() string {
	if len(pe.Message) == 0 && pe.Parent != nil {
		if pep, ok := pe.Parent.(ErrorMessage); ok {
			return pep.GetMessage()
		}
	}
	return pe.Message
}

// String returns formatted string.
func (pe *PointingError) String() string {
	var buf = bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	pe.Format(buf)
	return buf.String()
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 128))
	},
}

// FormatMessage formats giving message of an error.
func (pe *PointingError) FormatMessage(buf *bytes.Buffer) {
	if pe.Message != "" {
		buf.WriteString(pe.Message)
	}

	if pe.Parent != nil {
		if peHas, ok := pe.Parent.(HasMessage); ok && peHas.HasMessage() {
			buf.WriteString(": ")
		}
		if pem, ok := pe.Parent.(*PointingError); ok {
			pem.FormatMessage(buf)
		} else {
			buf.WriteString(pe.Parent.Error())
		}
	}
}

// FormatStack formats giving stack information for giving error.
func (pe *PointingError) FormatStack(buf *bytes.Buffer) {
	buf.WriteString("-------------------------------------------")
	var nb = njson.JSONB()
	if len(pe.Params) > 0 {
		buf.WriteString("\n")
		npkg.EncodableStringMap(pe.Params).EncodeObject(nb)
		_, _ = nb.WriteTo(buf)
		buf.WriteString("\n")
		buf.WriteString("-------------------------------------------")
	}
	buf.WriteString("\n")
	for _, frame := range pe.Frames {
		_, _ = fmt.Fprintf(buf, "- [%s] %s:%d", frame.Package, frame.File, frame.Line)
		buf.WriteString("\n")
	}
	if po, ok := pe.Parent.(*PointingError); ok {
		po.FormatStack(buf)
	}
}

// Format writes details of error into provided buffer.
func (pe *PointingError) Format(buf *bytes.Buffer) {
	pe.FormatMessage(buf)
	buf.WriteString("\n")
	pe.FormatStack(buf)
}

func IsSameTypeName(me interface{}, other interface{}) bool {
	var mc = reflect.TypeOf(me)
	var oc = reflect.TypeOf(other)
	return nameOfType(mc) == nameOfType(oc)
}
func nameOfType(t reflect.Type) string {
	stars := ""
	if t == nil {
		return "nil"
	}

	for t.Kind() == reflect.Ptr {
		stars += "*"
		t = t.Elem()
	}

	if t.Kind() == reflect.Interface {
		stars = ""
	}
	return t.PkgPath() + "/" + stars + t.Name()
}
