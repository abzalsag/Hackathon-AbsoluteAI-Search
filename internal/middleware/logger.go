package middleware

import (
	"net/http"
	"time"

	"hackathon-absolutecinema/pkg/logger" // <-- Стандартты "log" емес, осы импорт қосылды
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// <-- *log.Logger орнына *logger.Logger қолданылды
func LoggerMiddleware(next http.Handler, logger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		// Сіздің logger пакетіңізде Infof немесе Printf бар деп есептейміз
		logger.Infof("%s %s %d %s",
			r.Method,
			r.URL.Path,
			rw.status,
			time.Since(start),
		)
	})
}
