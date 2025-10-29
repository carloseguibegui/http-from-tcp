package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type Writer struct {
	writer io.Writer
}

type Response struct {
}

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func GetDefaultHeaders(contentLength int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-length", strconv.Itoa(contentLength))
	h.Set("Connection", "close")
	h.Set("Content-type", "text/plain")
	return h
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) WriteStatusLine(s StatusCode) error {
	statusLine := []byte("")
	switch s {
	case StatusOK:
		statusLine = []byte("HTTP/1.1 200 OK")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request")
	case StatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error")
	default:
		return fmt.Errorf("Unrecognized error code")
	}
	statusLine = append(statusLine, "\r\n"...)
	_, err := w.writer.Write(statusLine)
	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	b := []byte{}
	h.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Append(b, "\r\n")
	_, err := w.writer.Write(b)
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	return 0, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return 0, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	t := []byte{}
	h.ForEach(func(n, v string) {
		t = fmt.Appendf(t, "%s: %s\r\n", n, v)
	})
	t = fmt.Append(t, "\r\n")
	_, err := w.writer.Write(t)
	return err
}
