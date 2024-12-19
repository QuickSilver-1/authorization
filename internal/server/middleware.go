package server

import (
	"auth/internal/logger"
	"fmt"
	"net/http"
	"time"
)

type writer struct {
	http.ResponseWriter
	statusCode int
}

// Middleware логирует запросы и время их выполнения
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Log.Info(fmt.Sprintf("Request %s %s", r.Method, r.URL.Path))

		wrappedWriter := &writer{w, http.StatusOK}
		next.ServeHTTP(wrappedWriter, r)

		logger.Log.Info(fmt.Sprintf("Completed %s %s with %d in %v", r.Method, r.URL.Path, wrappedWriter.statusCode, time.Since(start)))
	})
}
