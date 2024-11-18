package middleware

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(next http.Handler) http.Handler

func Chain(middleware ...Middleware) http.Handler {
	var handler http.Handler
	for i := range middleware {
		handler = middleware[len(middleware)-1-i](handler)
	}

	return handler
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request
		log.Printf(
			"%s %s %s %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}
