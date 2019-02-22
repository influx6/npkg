package httptests

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gokit/npkg/nhttp"
)

// GET returns a new Context using GET method.
func Get(path string, body io.Reader, res *httptest.ResponseRecorder) *nhttp.NContext {
	return NewRequest("GET", path, body, res)
}

// Delete returns a new Context using DELETE method.
func Delete(path string, body io.Reader, res *httptest.ResponseRecorder) *nhttp.NContext {
	return NewRequest("DELETE", path, body, res)
}

// Put returns a new Context using PUT method.
func Put(path string, body io.Reader, res *httptest.ResponseRecorder) *nhttp.NContext {
	return NewRequest("PUT", path, body, res)
}

// Post returns a new Context using PUT method.
func Post(path string, body io.Reader, res *httptest.ResponseRecorder) *nhttp.NContext {
	return NewRequest("POST", path, body, res)
}

// Patch returns a new Context using PUT method.
func Patch(path string, body io.Reader, res *httptest.ResponseRecorder) *nhttp.NContext {
	return NewRequest("PATCH", path, body, res)
}

// NewRequest returns a new instance of a nhttp.Context with provided parameters.
func NewRequest(method string, path string, body io.Reader, res http.ResponseWriter) *nhttp.NContext {
	req := httptest.NewRequest(method, path, body)
	return nhttp.NewContext(
		nhttp.SetRequest(req),
		nhttp.SetResponseWriter(res),
	)
}
