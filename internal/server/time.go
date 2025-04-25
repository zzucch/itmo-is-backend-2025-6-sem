package server

import (
	"net/http"
	"time"
)

type timedResponseWriter struct {
	http.ResponseWriter
	start       time.Time
	wroteHeader bool
}

func (rw *timedResponseWriter) WriteHeader(statusCode int) {
	if !rw.wroteHeader {
		elapsed := time.Since(rw.start)
		rw.Header().Set("X-Elapsed-Time", elapsed.String())
		rw.wroteHeader = true
	}

	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *timedResponseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}

	return rw.ResponseWriter.Write(b)
}

func timeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &timedResponseWriter{
			ResponseWriter: w,
			start:          time.Now(),
			wroteHeader:    false,
		}

		next.ServeHTTP(rw, r)
	})
}
