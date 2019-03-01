package ntransport

import (
	"context"
	"io"
	"net/url"
	"time"
)

// set of possible values for Cookie.SameSite.
const (
	SameSiteDefaultMode = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
)

// A Cookie represents an HTTP cookie as sent in the Set-Cookie header of an
// HTTP response or the Cookie header of an HTTP request.
//
// See https://tools.ietf.org/html/rfc6265 for details.
type Cookie struct {
	Secure     bool
	HTTPOnly   bool
	Name       string
	Value      string
	RawExpires string // for reading cookies only
	Raw        string
	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	Unparsed   []string  // Raw text of unparsed attribute-value pairs

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge int

	// SameSite allows a server to define a cookie attribute making it impossible for
	// the browser to send this cookie along with cross-site requests. The main
	// goal is to mitigate the risk of cross-origin information leakage, and provide
	// some protection against cross-site request forgery attacks.
	//
	// See https://tools.ietf.org/html/draft-ietf-httpbis-cookie-same-site-00 for details.
	SameSite int
}

// Headers defines a map type providing similar header composition
// as used with http.Request and http.Response.
type Headers map[string][]string

// Add adds giving value into giving key if existing
// else creating new version.
func (h Headers) Add(key, value string) {
	h[key] = append(h[key], value)
}

// Set sets giving key with giving value as a new
// slice.
func (h Headers) Set(key, value string) {
	h[key] = []string{value}
}

// Get returns first value of giving key if it exists
// else returns an empty string.
func (h Headers) Get(key string) string {
	if vals, ok := h[key]; ok {
		if len(vals) > 0 {
			return vals[0]
		}
	}
	return ""
}

// Del removes giving key and associated value.
func (h Headers) Del(key string) {
	delete(h, key)
}

// Delete removes giving key and associated value.
func (h Headers) Delete(key string) {
	delete(h, key)
}

// Clone deep clones giving headers.
func (h Headers) Clone() Headers {
	return CloneHeader(h)
}

// CloneHeader clones giving header map returning a deep copy of it.
func CloneHeader(h Headers) Headers {
	h2 := make(map[string][]string, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

// Request implements a underline carrier of a request object which will be used
// by a transport to request giving resource.
type Request struct {
	Method  string
	Headers Headers
	Cookies []Cookie
	URL     *url.URL
	Body    io.ReadCloser
	Params  map[string]string
}

// Response defines the representation of a response from a tansport. This
// ensures we are able to represent responses as the transport see's fit and users
// are able to infer a giving response underline type through it's Provider() value
type Response interface {
	// Provider returns the unique representation of the underline response or provider
	// for which the response comes from. This allows user take provided giving string to
	// infer existing underline response type for futher access level.
	Provider() string

	// Status returns giving status of response after request was made.
	Status() int

	// Headers returns the headers received for giving response.
	Headers() Headers

	// Cookies returns the underline cookies attached to response if any.
	Cookies() []Cookie

	// Body returns underline body for response if any received from the server.
	Body() io.ReadCloser

	// Meta contains additional attached data from server or request.
	Meta() map[string]string
}

// Transport defines what we expect from a handler of requests.
// It will be responsible for the serialization of request to server and
// delivery of response or error from server.
type Transport interface {
	Send(ctx context.Context, request *Request) (Response, error)
}

//*****************************************************************************
// Services
//*****************************************************************************

// ResponseWriter defines what we expect from an implementation which handles the
// delivery of responses for a giving request.
type ResponseWriter interface {
	// WriteHead writes the response status and headers to the underline
	// connection. The response writer caches the status and headers and will
	// immediately write headers into response once ResponseWriter.Write is called
	// to write underline body. So the last call to WriteHead when ResponseWriter.Write
	// will prevail over previous.
	WriteHead(status int, headers Headers)

	// Write provides a method to stream respond body to underline
	// connection writer. A response writer underline writer will be
	// closed after the end of the Serve call.
	Write(b []byte) (int, error)

	// SentHeaders returns true/false if headers has been already
	// written by the call to ResponseWriter.WriteHead().
	SentHeaders() bool

	// Finished returns true/false if giving response header and body has
	// both be written and flushed to connection.
	Finished() bool
}

// Service defines what we expect from a service server which is responsible for the
// handling of incoming requests for a giving service type and the response for that giving
// request.
type Service interface {
	Serve(w ResponseWriter, r *Request)
}
