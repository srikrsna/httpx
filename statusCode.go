package httpx

import (
	"net/http"
	"net"
	"bufio"
	"errors"
)

type notFoundWriter struct {
	http.ResponseWriter
	status int
}

func New404Handler(next http.Handler, notFound http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		nfw := &notFoundWriter{ResponseWriter: w}

		next.ServeHTTP(nfw, r)

		if nfw.status == http.StatusNotFound {
			notFound.ServeHTTP(w, r)
		}
	})
}

func (h *notFoundWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok :=  h.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return  nil, nil, errors.New("not a hijacker")
}

func (h *notFoundWriter) Write(b []byte) (int, error) {
	if h.status == http.StatusNotFound {
		return len(b), nil
	}

	return h.ResponseWriter.Write(b)
}

func (h *notFoundWriter) WriteHeader(statusCode int) {
	h.status = statusCode

	if statusCode != http.StatusNotFound {
		h.ResponseWriter.WriteHeader(statusCode)
	}
}

