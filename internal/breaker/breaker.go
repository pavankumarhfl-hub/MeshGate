package breaker

import (
	"errors"
	"sync"
	"time"
)

var ErrOpen = errors.New("circuit breaker open")

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type Breaker struct {
	mu        sync.Mutex
	state     State
	failures  int
	threshold int
	openedAt  time.Time
	cooldown  time.Duration
	probeInFlight bool
}

func New(threshold int, cooldown time.Duration) *Breaker {
	return &Breaker{threshold: threshold, cooldown: cooldown}
}

func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case Closed:
		return nil
	case Open:
		if time.Since(b.openedAt) < b.cooldown {
			return ErrOpen
		}
		b.state = HalfOpen
		b.probeInFlight = true
		return nil
	case HalfOpen:
		if b.probeInFlight {
			return ErrOpen
		}
		b.probeInFlight = true
		return nil
	default:
		return ErrOpen
	}
}

func (b *Breaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.probeInFlight = false
	b.state = Closed
}

func (b *Breaker) Failure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.probeInFlight = false
	b.failures++
	if b.state == HalfOpen || b.failures >= b.threshold {
		b.state = Open
		b.openedAt = time.Now()
	}
}

func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == Open && time.Since(b.openedAt) >= b.cooldown {
		return HalfOpen
	}
	return b.state
}
