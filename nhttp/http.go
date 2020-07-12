package nhttp

import (
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/influx6/npkg"
	"github.com/influx6/npkg/njson"
)

const (
	maxsize = 32 << 40

	// MultipartKey defines the key used to store multipart Form.
	MultipartKey = "MultiPartForm"
)

// HTTPError defines custom error that can be used to specify
// status code and message.
type HTTPError struct {
	Code int
	Err  error
}

// Error returns error string. Implements error interface.
func (h HTTPError) Error() string {
	return h.Err.Error()
}

// Gzip Compression
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// handlerImpl implements http.Handler interface.
type handlerImpl struct {
	ContextHandler
}

// ServeHTTP implements http.Handler.ServeHttp method.
func (h handlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(SetRequest(r), SetResponseWriter(w))
	if err := ctx.InitForms(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer ctx.ClearFlashMessages()

	if err := h.ContextHandler(ctx); err != nil {
		if httperr, ok := err.(HTTPError); ok {
			http.Error(w, httperr.Error(), httperr.Code)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// HTTPFunc returns a http.HandleFunc which wraps the Handler for usage
// with a server.
func HTTPFunc(nx ContextHandler, befores ...func()) http.HandlerFunc {
	return handlerImpl{ContextHandler: func(ctx *NContext) error {
		ctx.response.beforeFuncs = append(ctx.response.beforeFuncs, befores...)
		return nx(ctx)
	}}.ServeHTTP
}

// ServeHandler returns a http.Handler which serves request to the provided Handler.
func ServeHandler(h ContextHandler) http.Handler {
	return handlerImpl{ContextHandler: h}
}

// GetFileMimeType returns associated mime type for giving file extension.
func GetFileMimeType(path string) string {
	ext := filepath.Ext(path)
	extVal := mime.TypeByExtension(ext)
	if extVal == "" {
		extVal = mediaTypes[ext]
	}
	return extVal
}

// JSONError writes the giving error message to the provided writer.
func JSONError(w http.ResponseWriter, statusCode int, errorCode string, message string, err error) error {
	w.WriteHeader(statusCode)

	var encoder = njson.JSONB()
	encoder.ObjectFor("error", func(enc npkg.ObjectEncoder) {
		enc.String("message", message)
		enc.Int("status_code", statusCode)
		enc.String("error_code", errorCode)

		if encodableErr, ok := err.(npkg.EncodableObject); ok {
			enc.Object("incident", encodableErr)
		}
		enc.String("incident", err.Error())
	})

	var _, werr = encoder.WriteTo(w)
	return werr
}

// ParseAuthorization returns the scheme and token of the Authorization string
// if it's valid.
func ParseAuthorization(val string) (authType string, token string, err error) {
	authSplit := strings.SplitN(val, " ", 2)
	if len(authSplit) != 2 {
		err = errors.New("invalid authorization: Expected content: `AuthType Token`")
		return
	}

	authType = strings.TrimSpace(authSplit[0])
	token = strings.TrimSpace(authSplit[1])

	return
}

// ParseTokens parses the base64 encoded token sent as part of the Authorization string,
// It expects all parts of string to be seperated with ':', returning splitted slice.
func ParseTokens(val string) ([]string, error) {
	decoded, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(decoded), ":"), nil
}
