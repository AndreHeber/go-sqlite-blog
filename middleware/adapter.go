package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

// First, create context keys for request and response writer
type ctxKey string
const (
    requestKey  ctxKey = "request"
    responseKey ctxKey = "response"
)

type Adapter struct {
	errorInResponse bool
}

func Init(errorInResponse bool) *Adapter {
	return &Adapter{errorInResponse: errorInResponse}
}

// Create an adapter function
func (a *Adapter) HttpToContextHandler(h func(context.Context) error) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), requestKey, r)
        ctx = context.WithValue(ctx, responseKey, w)
        
        if err := h(ctx); err != nil {
			slog.Error("middleware: HttpToContextHandler", "error", err)

            // Handle error appropriately
			if a.errorInResponse {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				// return common error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
        }
    }
}

// Helper functions to extract from context
func GetRequest(ctx context.Context) *http.Request {
    return ctx.Value(requestKey).(*http.Request)
}

func GetResponseWriter(ctx context.Context) http.ResponseWriter {
    return ctx.Value(responseKey).(http.ResponseWriter)
}
