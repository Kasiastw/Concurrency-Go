package limiter

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

type ipLimit struct {
	mu     sync.Mutex
	window time.Duration
	limit  int
	items  map[string]*rate.Limiter
}

func (l ipLimit) Allow(addr string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if lim, ok := l.items[addr]; ok {
		return lim.Allow()
	}
	l.items[addr] = rate.NewLimiter(rate.Every(l.window), l.limit)
	return true
}

var tooManyRequestsMessage = []byte("The limit of requests has been reached")

func newLimiter(window time.Duration, limit int) ipLimit {
	return ipLimit{
		mu:     sync.Mutex{},
		window: window,
		limit:  limit,
		items:  map[string]*rate.Limiter{},
	}
}

func LimitHandler(window time.Duration, limit int, nextHandler http.HandlerFunc) http.HandlerFunc {
	lim := newLimiter(window, limit)
	return func(w http.ResponseWriter, r *http.Request) {
		if lim.Allow(r.RemoteAddr) {
			nextHandler(w, r)
		} else {
			limitationExhaustedHandler(w, r)
		}
	}
}

func limitationExhaustedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write(tooManyRequestsMessage)
}
