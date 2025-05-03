package server

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
)

type etagResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	headersSent bool
	body        *bytes.Buffer
}

func (rw *etagResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

func (rw *etagResponseWriter) Write(b []byte) (int, error) {
	if rw.body == nil {
		rw.body = &bytes.Buffer{}
	}
	return rw.body.Write(b)
}

func etagMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &etagResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           &bytes.Buffer{},
		}

		next.ServeHTTP(rw, r)

		sum := sha1.Sum(rw.body.Bytes())
		etag := `"` + hex.EncodeToString(sum[:]) + `"`

		if match := r.Header.Get("If-None-Match"); match == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.Header().Set("ETag", etag)
		for k, v := range rw.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rw.statusCode)
		_, _ = io.Copy(w, rw.body)
	})
}
