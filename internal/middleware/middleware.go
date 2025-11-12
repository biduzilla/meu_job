package middleware

import (
	"expvar"
	"fmt"
	"meu_job/internal/config"
	"meu_job/internal/services"
	"meu_job/utils/errors"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	totalRequestsReceived           = expvar.NewInt("total_requests_received")
	totalResponsesSent              = expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_μs")
	totalResponsesSentByStatus      = expvar.NewMap("total_responses_sent_by_status")
)

type Middleware struct {
	errRsp      errors.ErrorResponseInterface
	userService services.UserService
	config      config.Config
}

func New(
	errRsp errors.ErrorResponseInterface,
	userService services.UserService,
	config config.Config,
) *Middleware {
	return &Middleware{
		errRsp:      errRsp,
		userService: userService,
		config:      config,
	}
}

func (m *Middleware) RateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.config.Limiter.Enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				m.errRsp.ServerErrorResponse(w, r, err)
				return
			}

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(
						rate.Limit(m.config.Limiter.RPS),
						m.config.Limiter.Burst,
					),
				}

				clients[ip].lastSeen = time.Now()
				if !clients[ip].limiter.Allow() {
					mu.Unlock()
					m.errRsp.RateLimitExceededResponse(w, r)
					return
				}
				mu.Unlock()

			}
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				m.errRsp.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
