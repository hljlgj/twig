package twig

import (
	"net/http"
	"reflect"
	"strings"
	"unsafe"
)

func WriteContentType(w http.ResponseWriter, val string) {
	header := w.Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, val)
	}
}

func IsTLS(r *http.Request) bool {
	return r.TLS != nil
}

func Byte(w http.ResponseWriter, code int, contentType string, bs []byte) (err error) {
	WriteContentType(w, contentType)
	w.WriteHeader(code)
	_, err = w.Write(bs)
	return
}

func UnsafeToBytes(s string) []byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: strHeader.Data,
		Len:  strHeader.Len,
		Cap:  strHeader.Len,
	}))
}

func UnsafeString(w http.ResponseWriter, code int, str string) error {
	return Byte(w, code, MIMETextPlainCharsetUTF8, UnsafeToBytes(str))
}

func String(w http.ResponseWriter, code int, str string) error {
	return Byte(w, code, MIMETextPlainCharsetUTF8, []byte(str))
}

func IsWebSocket(r *http.Request) bool {
	upgrade := r.Header.Get(HeaderUpgrade)
	return upgrade == "websocket" || upgrade == "Websocket"
}

func IsXMLHTTPRequest(r *http.Request) bool {
	return strings.Contains(
		r.Header.Get(HeaderXRequestedWith),
		XMLHttpRequest,
	)
}

func Scheme(r *http.Request) string {
	if IsTLS(r) {
		return "https"
	}
	if scheme := r.Header.Get(HeaderXForwardedProto); scheme != "" {
		return scheme
	}
	if scheme := r.Header.Get(HeaderXForwardedProtocol); scheme != "" {
		return scheme
	}
	if ssl := r.Header.Get(HeaderXForwardedSsl); ssl == "on" {
		return "https"
	}
	if scheme := r.Header.Get(HeaderXUrlScheme); scheme != "" {
		return scheme
	}
	return "http"
}
