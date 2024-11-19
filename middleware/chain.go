package middleware

import (
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

func Chain(middleware ...Middleware) http.Handler {
	var handler http.Handler
	for i := range middleware {
		handler = middleware[len(middleware)-1-i](handler)
	}

	return handler
}
