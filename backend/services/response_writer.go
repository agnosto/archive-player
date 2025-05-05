package services

import (
	"bytes"
	"net/http"
	"strconv"
)

// CustomResponseWriter implements http.ResponseWriter for Wails
type CustomResponseWriter struct {
	headers    http.Header
	body       *bytes.Buffer
	statusCode int
}

// NewCustomResponseWriter creates a new custom response writer
func NewCustomResponseWriter() *CustomResponseWriter {
	return &CustomResponseWriter{
		headers:    make(http.Header),
		body:       new(bytes.Buffer),
		statusCode: http.StatusOK,
	}
}

// Header returns the header map
func (w *CustomResponseWriter) Header() http.Header {
	return w.headers
}

// Write writes the data to the buffer
func (w *CustomResponseWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

// WriteHeader sets the status code
func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// GetResponse returns the response as a map
func (w *CustomResponseWriter) GetResponse() map[string]interface{} {
	response := make(map[string]interface{})
	response["statusCode"] = w.statusCode
	response["headers"] = w.headers
	response["body"] = w.body.Bytes()

	// Add Content-Length header if not present
	if _, exists := w.headers["Content-Length"]; !exists {
		w.headers.Set("Content-Length", strconv.Itoa(w.body.Len()))
	}

	return response
}
