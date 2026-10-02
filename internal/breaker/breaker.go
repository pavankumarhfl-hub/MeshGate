package breaker

import (
	"errors"
	"sync"
	"time"
)

var ErrOpen = errors.New("circuit breaker open")

type State int
const ( Closed State = iota; Open; HalfOpen )

type Breaker struct {
	mu sync.Mutex
	state State
	failures int
	threshold int
	openedAt time.Time
	cooldown time.Duration
}

func New(threshold int, cooldown time.Duration) *Breaker { return &Breaker{threshold: threshold, cooldown: cooldown} }

func (b *Breaker) Allow() error {
	b.mu.Lock(); defer b.mu.Unlock()
	if b.state == Open && time.Since(b.openedAt) < b.cooldown { return ErrOpen }
	if b.state == Open { b.state = HalfOpen }
	return nil
}

func (b *Breaker) Success() { b.mu.Lock(); defer b.mu.Unlock(); b.failures = 0; b.state = Closed }
func (b *Breaker) Failure() { b.mu.Lock(); defer b.mu.Unlock(); b.failures++; if b.failures >= b.threshold { b.state = Open; b.openedAt = time.Now() } }
func (b *Breaker) State() State { b.mu.Lock(); defer b.mu.Unlock(); if b.state == Open && time.Since(b.openedAt) >= b.cooldown { return HalfOpen }; return b.state }
