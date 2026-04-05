package middleware

import (
	"bufio"
	"bytes"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	noWritten = -1
)

type resWriter struct {
	Writer  gin.ResponseWriter // the actual ResponseWriter to flush to
	status  int                // the HTTP response code from WriteHeader
	Body    *bytes.Buffer      // the response content body
	Flushed bool
}

func NewResWriter(w gin.ResponseWriter) resWriter {
	return resWriter{
		Writer: w,
		status: w.Status(),
		Body:   &bytes.Buffer{},
	}
}

func (w resWriter) Header() http.Header {
	return w.Writer.Header() // use the actual response header
}

func (w resWriter) Write(buf []byte) (int, error) {
	w.Body.Write(buf)
	return len(buf), nil
}

func (w resWriter) WriteString(s string) (n int, err error) {
	n, err = w.Write([]byte(s))
	return
}

func (w resWriter) Written() bool {
	return w.Body.Len() != noWritten
}

func (w resWriter) WriteHeader(status int) {
	w.status = status
}

func (w resWriter) WriteHeaderNow() {
}

func (w resWriter) Status() int {
	return w.status
}

func (w resWriter) Size() int {
	return w.Body.Len()
}

func (w resWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.Writer.(http.Hijacker).Hijack()
}

func (w resWriter) CloseNotify() <-chan bool {
	return w.Writer.(http.CloseNotifier).CloseNotify()
}

func (w resWriter) Pusher() http.Pusher {
	return w.Pusher()
}

// Fake Flush
// TBD
func (w resWriter) Flush() {
	w.realFlush()
}

func (w *resWriter) realFlush() {
	if w.Flushed {
		return
	}
	w.Writer.WriteHeader(w.status)
	if w.Body.Len() > 0 {
		_, err := w.Writer.Write(w.Body.Bytes())
		if err != nil {
			panic(err)
		}
		w.Body.Reset()
	}
	w.Flushed = true
}
