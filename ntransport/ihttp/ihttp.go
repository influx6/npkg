package ihttp

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/influx6/npkg/ntransport"
)

// Response is an implementation of http.ResponseWriter that
// records its mutations for later inspection in tests.
type Response struct {

	// Code is the HTTP response code set by WriteHeader.
	//
	// Note that if a Handler never calls WriteHeader or Write,
	// this might end up being 0, rather than the implicit
	// http.StatusOK. To get the implicit value, use the Result
	// method.
	Code int

	// HeaderMap contains the headers explicitly set by the Handler.
	// It is an internal detail.
	//
	// Deprecated: HeaderMap exists for historical compatibility
	// and should not be used. To access the headers returned by a handler,
	// use the Response.Header map as returned by the Result method.
	HeaderMap ntransport.Headers

	// Body is the buffer to which the Handler's Write calls are sent.
	// If nil, the Writes are silently discarded.
	Body *bytes.Buffer

	// Flushed is whether the Handler called Flush.
	Flushed bool

	snapHeader  map[string][]string // snapshot of HeaderMap at first Write
	wroteHeader bool
}

// NewRecorder returns an initialized Response.
func NewRecorder() *Response {
	return &Response{
		HeaderMap: make(map[string][]string),
		Body:      new(bytes.Buffer),
		Code:      200,
	}
}

// DefaultRemoteAddr is the default remote address to return in RemoteAddr if
// an explicit DefaultRemoteAddr isn't set on Response.
const DefaultRemoteAddr = "1.2.3.4"

// Header returns the response headers.
func (rw *Response) Header() map[string][]string {
	m := rw.HeaderMap
	if m == nil {
		m = make(map[string][]string)
		rw.HeaderMap = m
	}
	return m
}

// writeHeader writes a header if it was not written yet and
// detects Content-Type if needed.
//
// bytes or str are the beginning of the response body.
// We pass both to avoid unnecessarily generate garbage
// in rw.WriteString which was created for performance reasons.
// Non-nil bytes win.
func (rw *Response) writeHeader(b []byte, str string) {
	if rw.wroteHeader {
		return
	}

	if len(str) > 512 {
		str = str[:512]
	}

	m := rw.Header()
	_, hasType := m["Content-Type"]
	hasTE := len(m["Transfer-Encoding"]) != 0
	if !hasType && !hasTE {
		if b == nil {
			b = []byte(str)
		}
		m["Content-Type"] = append(m["Content-Type"], DetectContentType(b))
	}
	rw.WriteHeader(200)
}

// Write always succeeds and writes to rw.Body, if not nil.
func (rw *Response) Write(buf []byte) (int, error) {
	rw.writeHeader(buf, "")
	if rw.Body != nil {
		rw.Body.Write(buf)
	}
	return len(buf), nil
}

// WriteString always succeeds and writes to rw.Body, if not nil.
func (rw *Response) WriteString(str string) (int, error) {
	rw.writeHeader(nil, str)
	if rw.Body != nil {
		rw.Body.WriteString(str)
	}
	return len(str), nil
}

// WriteHeader sets rw.Code. After it is called, changing rw.Header
// will not affect rw.HeaderMap.
func (rw *Response) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.Code = code
	rw.wroteHeader = true
	if rw.HeaderMap == nil {
		rw.HeaderMap = make(map[string][]string)
	}
	rw.snapHeader = CloneHeader(rw.HeaderMap)
}

// Flush sets rw.Flushed to true.
func (rw *Response) Flush() {
	if !rw.wroteHeader {
		rw.WriteHeader(200)
	}
	rw.Flushed = true

}

//**************************************************************************
// internal functions from http, and httptest and /internal/net/httpguts
//**************************************************************************

// ParseContentLength trims whitespace from s and returns -1 if no value
// is set, or the value if it's >= 0.
//
// This a modified version of same function found in net/http/transfer.go. This
// one just ignores an invalid header.
func ParseContentLength(cl string) int64 {
	cl = strings.TrimSpace(cl)
	if cl == "" {
		return -1
	}

	n, err := strconv.ParseInt(cl, 10, 64)
	if err != nil {
		return -1
	}
	return n
}

// CloneHeader clones giving header map returning a deep copy of it.
func CloneHeader(h map[string][]string) map[string][]string {
	h2 := make(map[string][]string, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

// The algorithm uses at most sniffLen bytes to make its decision.
const sniffLen = 512

// DetectContentType implements the algorithm described
// at https://mimesniff.spec.whatwg.org/ to determine the
// Content-Type of the given data. It considers at most the
// first 512 bytes of data. DetectContentType always returns
// a valid MIME type: if it cannot determine a more specific one, it
// returns "application/octet-stream".
func DetectContentType(data []byte) string {
	if len(data) > sniffLen {
		data = data[:sniffLen]
	}

	// Index of the first non-whitespace byte in data.
	firstNonWS := 0
	for ; firstNonWS < len(data) && isWS(data[firstNonWS]); firstNonWS++ {
	}

	for _, sig := range sniffSignatures {
		if ct := sig.match(data, firstNonWS); ct != "" {
			return ct
		}
	}
	return "application/octet-stream" // fallback
}

func isWS(b byte) bool {
	switch b {
	case '\t', '\n', '\x0c', '\r', ' ':
		return true
	}
	return false
}

type sniffSig interface {
	// match returns the MIME type of the data, or "" if unknown.
	match(data []byte, firstNonWS int) string
}

// Data matching the table in section 6.
var sniffSignatures = []sniffSig{

	htmlSig("<!DOCTYPE HTML"),

	htmlSig("<HTML"),

	htmlSig("<HEAD"),

	htmlSig("<SCRIPT"),

	htmlSig("<IFRAME"),

	htmlSig("<H1"),

	htmlSig("<DIV"),

	htmlSig("<FONT"),

	htmlSig("<TABLE"),

	htmlSig("<A"),

	htmlSig("<STYLE"),

	htmlSig("<TITLE"),

	htmlSig("<B"),

	htmlSig("<BODY"),

	htmlSig("<BR"),

	htmlSig("<P"),

	htmlSig("<!--"),

	&maskedSig{mask: []byte("\xFF\xFF\xFF\xFF\xFF"), pat: []byte("<?xml"), skipWS: true, ct: "text/xml; charset=utf-8"},

	&exactSig{[]byte("%PDF-"), "application/pdf"},

	&exactSig{[]byte("%!PS-Adobe-"), "application/postscript"},

	// UTF BOMs.

	&maskedSig{mask: []byte("\xFF\xFF\x00\x00"), pat: []byte("\xFE\xFF\x00\x00"), ct: "text/plain; charset=utf-16be"},

	&maskedSig{mask: []byte("\xFF\xFF\x00\x00"), pat: []byte("\xFF\xFE\x00\x00"), ct: "text/plain; charset=utf-16le"},

	&maskedSig{mask: []byte("\xFF\xFF\xFF\x00"), pat: []byte("\xEF\xBB\xBF\x00"), ct: "text/plain; charset=utf-8"},

	&exactSig{[]byte("GIF87a"), "image/gif"},

	&exactSig{[]byte("GIF89a"), "image/gif"},

	&exactSig{[]byte("\x89\x50\x4E\x47\x0D\x0A\x1A\x0A"), "image/png"},

	&exactSig{[]byte("\xFF\xD8\xFF"), "image/jpeg"},

	&exactSig{[]byte("BM"), "image/bmp"},

	&maskedSig{

		mask: []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF\xFF\xFF"),

		pat: []byte("RIFF\x00\x00\x00\x00WEBPVP"),

		ct: "image/webp",
	},

	&exactSig{[]byte("\x00\x00\x01\x00"), "image/vnd.microsoft.icon"},

	&maskedSig{

		mask: []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF"),

		pat: []byte("RIFF\x00\x00\x00\x00WAVE"),

		ct: "audio/wave",
	},

	&maskedSig{

		mask: []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF"),

		pat: []byte("FORM\x00\x00\x00\x00AIFF"),

		ct: "audio/aiff",
	},

	&maskedSig{

		mask: []byte("\xFF\xFF\xFF\xFF"),

		pat: []byte(".snd"),

		ct: "audio/basic",
	},

	&maskedSig{

		mask: []byte("\xFF\xFF\xFF\xFF\xFF"),

		pat: []byte("OggS\x00"),

		ct: "application/ogg",
	},

	&maskedSig{

		mask: []byte("\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"),

		pat: []byte("MThd\x00\x00\x00\x06"),

		ct: "audio/midi",
	},

	&maskedSig{

		mask: []byte("\xFF\xFF\xFF"),

		pat: []byte("ID3"),

		ct: "audio/mpeg",
	},

	&maskedSig{

		mask: []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF"),

		pat: []byte("RIFF\x00\x00\x00\x00AVI "),

		ct: "video/avi",
	},

	// Fonts

	&maskedSig{

		// 34 NULL bytes followed by the string "LP"

		pat: []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x4C\x50"),

		// 34 NULL bytes followed by \xF\xF

		mask: []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xFF\xFF"),

		ct: "application/vnd.ms-fontobject",
	},

	&exactSig{[]byte("\x00\x01\x00\x00"), "font/ttf"},

	&exactSig{[]byte("OTTO"), "font/otf"},

	&exactSig{[]byte("ttcf"), "font/collection"},

	&exactSig{[]byte("wOFF"), "font/woff"},

	&exactSig{[]byte("wOF2"), "font/woff2"},

	&exactSig{[]byte("\x1A\x45\xDF\xA3"), "video/webm"},

	&exactSig{[]byte("\x52\x61\x72\x20\x1A\x07\x00"), "application/x-rar-compressed"},
	&exactSig{[]byte("\x50\x4B\x03\x04"), "application/zip"},
	&exactSig{[]byte("\x1F\x8B\x08"), "application/x-gzip"},
	&exactSig{[]byte("\x00\x61\x73\x6D"), "application/wasm"},
	mp4Sig{},
	textSig{}, // should be last
}

type exactSig struct {
	sig []byte
	ct  string
}

func (e *exactSig) match(data []byte, firstNonWS int) string {
	if bytes.HasPrefix(data, e.sig) {
		return e.ct
	}
	return ""
}

type maskedSig struct {
	mask, pat []byte
	skipWS    bool
	ct        string
}

func (m *maskedSig) match(data []byte, firstNonWS int) string {

	// pattern matching algorithm section 6
	// https://mimesniff.spec.whatwg.org/#pattern-matching-algorithm
	if m.skipWS {
		data = data[firstNonWS:]
	}

	if len(m.pat) != len(m.mask) {
		return ""
	}

	if len(data) < len(m.mask) {
		return ""
	}

	for i, mask := range m.mask {
		db := data[i] & mask
		if db != m.pat[i] {
			return ""
		}
	}
	return m.ct
}

type htmlSig []byte

func (h htmlSig) match(data []byte, firstNonWS int) string {

	data = data[firstNonWS:]
	if len(data) < len(h)+1 {
		return ""
	}

	for i, b := range h {
		db := data[i]
		if 'A' <= b && b <= 'Z' {
			db &= 0xDF
		}

		if b != db {
			return ""
		}
	}

	// Next byte must be space or right angle bracket.
	if db := data[len(h)]; db != ' ' && db != '>' {
		return ""
	}

	return "text/html; charset=utf-8"
}

var mp4ftype = []byte("ftyp")

var mp4 = []byte("mp4")

type mp4Sig struct{}

func (mp4Sig) match(data []byte, firstNonWS int) string {

	// https://mimesniff.spec.whatwg.org/#signature-for-mp4
	// c.f. section 6.2.1
	if len(data) < 12 {
		return ""
	}

	boxSize := int(binary.BigEndian.Uint32(data[:4]))
	if boxSize%4 != 0 || len(data) < boxSize {
		return ""
	}

	if !bytes.Equal(data[4:8], mp4ftype) {
		return ""
	}

	for st := 8; st < boxSize; st += 4 {
		if st == 12 {
			// minor version number
			continue
		}

		if bytes.Equal(data[st:st+3], mp4) {
			return "video/mp4"
		}
	}

	return ""
}

type textSig struct{}

func (textSig) match(data []byte, firstNonWS int) string {
	// c.f. section 5, step 4.
	for _, b := range data[firstNonWS:] {
		switch {
		case b <= 0x08,
			b == 0x0B,
			0x0E <= b && b <= 0x1A,
			0x1C <= b && b <= 0x1F:
			return ""
		}
	}

	return "text/plain; charset=utf-8"
}
