package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Adapter struct {
	Request         *http.Request
	ResponseWriter  http.ResponseWriter
	Logger          *slog.Logger
	DB              *sql.DB
	Ctx             context.Context
	Cancel          context.CancelFunc
	ErrorInResponse bool
	LogDBQueries    bool
	ipRateLimiter   *IPRateLimiter
}

func Init(logger *slog.Logger, db *sql.DB, errorInResponse bool, logDBQueries bool, ipRateLimit rate.Limit, burst int) *Adapter {
	return &Adapter{Logger: logger, DB: db, ErrorInResponse: errorInResponse, LogDBQueries: logDBQueries, ipRateLimiter: NewIPRateLimiter(ipRateLimit, burst)}
}

// 1. Simple rate limiter per IP
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.RWMutex
	r   rate.Limit
	b   int
}

func (l *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists := l.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(l.r, l.b)
		l.ips[ip] = limiter
	}

	return limiter
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}
}

// Create an adapter function
func (a *Adapter) HTTPToContextHandler(h func(*Adapter) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limiter := a.ipRateLimiter.getLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			a.Logger.Info("middleware: HttpToContextHandler", "error", "Too many requests", "ip", r.RemoteAddr)
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		if err := r.ParseForm(); err != nil {
			a.Logger.Error("middleware: HttpToContextHandler", "error", err)
			return
		}

		formValues := make([]string, 0, len(r.Form))
		for key, value := range r.Form {
			formValues = append(formValues, fmt.Sprintf("%s=%s", key, value))
		}
		parameters := strings.Join(formValues, "&")

		a.Request = r
		a.ResponseWriter = w

		ctx := context.Background()
		a.Ctx, a.Cancel = context.WithTimeout(ctx, 10*time.Second)
		defer a.Cancel()

		start := time.Now()
		if err := h(a); err != nil {
			a.Logger.Error("middleware: HttpToContextHandler", "error", err)

			// Handle error appropriately
			if a.ErrorInResponse {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				// return common error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}

		// log response info, parameters and duration
		var requestInfo string
		if parameters != "" {
			requestInfo = fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, parameters)
		} else {
			requestInfo = fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		}
		duration := time.Since(start)
		a.Logger.Info(requestInfo, "duration", duration.String())
	}
}
