package broadcast

import (
	"bytes"
	"io"
	"net/http"
)

type responseWrapper struct {
	buffer *bytes.Buffer
	header http.Header
	status int
}

func newResponseWrapper() *responseWrapper {
	return &responseWrapper{
		buffer: bytes.NewBuffer(nil),
		header: make(http.Header),
	}
}

// copy all content to parent response writer.
func (w *responseWrapper) copy(rw http.ResponseWriter) {
	rw.WriteHeader(w.status)

	for k, v := range w.header {
		for _, vv := range v {
			rw.Header().Add(k, vv)
		}
	}
	io.Copy(rw, w.buffer)
}

// Header returns the header map that will be sent by WriteHeader.
func (w *responseWrapper) Header() http.Header {
	return w.header
}

// Write writes the data to the connection as part of an HTTP reply.
func (w *responseWrapper) Write(p []byte) (int, error) {
	return w.buffer.Write(p)
}

// WriteHeader sends an HTTP response header with the provided status code.
func (w *responseWrapper) WriteHeader(statusCode int) {
	w.status = statusCode
}

// IsOK returns true if response contained 200 status code.
func (w *responseWrapper) IsOK() bool {
	return w.status == 200
}

// Body returns response body to rely to user.
func (w *responseWrapper) Body() []byte {
	return w.buffer.Bytes()
}
