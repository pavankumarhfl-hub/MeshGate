package ratelimit

import (
	"net/http"
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
	limit  int
	burst  int
	window time.Duration
	seen   map[string]*entry
}

type entry struct { count int; reset time.Time }

func New(limit, burst int) *Limiter {
	return &Limiter{limit: limit, burst: burst, window: time.Second, seen: make(map[string]*entry)}
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock(); defer l.mu.Unlock()
	now := time.Now()
	e := l.seen[key]
	if e == nil || now.After(e.reset) { e = &entry{reset: now.Add(l.window)}; l.seen[key] = e }
	cap := l.limit
	if l.burst > cap { cap = l.burst }
	if e.count >= cap { return false }
	e.count++
	return true
}

func Middleware(l *Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Forwarded-For")
		if key == "" { key = r.RemoteAddr }
		if !l.Allow(key) { http.Error(w, "rate limit exceeded", http.StatusTooManyRequests); return }
		next.ServeHTTP(w, r)
	})
}
