package middlewares

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	}
	return lrw.ResponseWriter.Write(b)
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w}
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

// PanicRecoveryMiddleware recovers from panics
func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				log.Printf("Panic recovered: %v", err)

				// Return error to client
				responseWithJSON(w, http.StatusInternalServerError, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// TODO: To move to another place to be reused and avoid duplicated code
func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
