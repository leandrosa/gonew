package middlewares

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func WrapHandlerWithLogging(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		log.Printf("--> %s %s", req.Method, req.URL.Path)

		for key, value := range req.Header {
			// log.Printf("--> key:%s", key)
			// log.Printf("--> value:%s", value)
			log.Printf("----> %s: %s", key, value)
		}

		lrw := NewLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, req)

		statusCode := lrw.statusCode
		log.Printf("<-- %d %s completed in %s", statusCode, http.StatusText(statusCode), time.Since(start))
	})
}
