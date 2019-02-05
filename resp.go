package twig

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

// RespCallBack ResponseWriter回调函数
type RespCallBack func(http.ResponseWriter)

// ResponseWrap 包装http.ResponseWrite
type ResponseWrap struct {
	before    []RespCallBack
	after     []RespCallBack
	Writer    http.ResponseWriter
	Status    int
	Len       int64
	Committed bool
}

func newResponseWrap(w http.ResponseWriter) *ResponseWrap {
	return &ResponseWrap{Writer: w}
}

func (r *ResponseWrap) Header() http.Header {
	return r.Writer.Header()
}

func (r *ResponseWrap) Flush() {
	r.Writer.(http.Flusher).Flush()
}

func (r *ResponseWrap) Before(fn RespCallBack) {
	r.before = append(r.before, fn)
}

func (r *ResponseWrap) After(fn RespCallBack) {
	r.after = append(r.after, fn)
}

// 设置header时查是否已经输出内容
func (r *ResponseWrap) WriteHeader(code int) {
	if r.Committed {
		return
	}
	for _, fn := range r.before {
		fn(r)
	}
	r.Status = code
	r.Writer.WriteHeader(code)
	r.Committed = true
}

// 输出时候检查是否设置Header
func (r *ResponseWrap) Write(b []byte) (n int, err error) {
	if !r.Committed {
		r.WriteHeader(http.StatusOK)
	}
	n, err = r.Writer.Write(b)
	r.Len += int64(n)
	for _, fn := range r.after {
		fn(r)
	}
	return
}

func (r *ResponseWrap) ReadFrom(src io.Reader) (n int64, e error) {
	if !r.Committed {
		r.WriteHeader(http.StatusOK)
	}

	n, e = io.Copy(r.Writer, src)
	r.Len += n
	return

}

// Hijack Hijack 支持
func (r *ResponseWrap) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.Writer.(http.Hijacker).Hijack()
}

func (r *ResponseWrap) reset(w http.ResponseWriter) {
	r.before = nil
	r.after = nil
	r.Writer = w
	r.Len = 0
	r.Status = http.StatusOK
	r.Committed = false
}
