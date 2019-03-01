package nconns

import (
	"io"
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

// Request implements a underline carrier of a request object.
type Request struct {
	URL     string
	Method  string
	Cookies []Cookie
	Body    io.ReadCloser
	Headers Headers
}
