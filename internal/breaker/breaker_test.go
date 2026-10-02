package breaker

import (
	"sync"
	"testing"
	"time"
)

func TestBreakerOpensAfterThreshold(t *testing.T) {
	b := New(2, time.Second)
	b.Failure()
	b.Failure()
	if b.State() != Open { t.Fatalf("expected open state, got %v", b.State()) }
	if b.Allow() != ErrOpen { t.Fatal("expected requests to be rejected while open") }
}

func TestBreakerRecovers(t *testing.T) {
	b := New(1, time.Millisecond)
	b.Failure()
	time.Sleep(3 * time.Millisecond)
	if err := b.Allow(); err != nil { t.Fatal(err) }
	b.Success()
	if b.State() != Closed { t.Fatal("expected closed state after success") }
}

func TestBreakerAllowsSingleHalfOpenProbe(t *testing.T) {
	b := New(1, time.Millisecond)
	b.Failure()
	time.Sleep(3 * time.Millisecond)

	var wg sync.WaitGroup
	allowed := make(chan struct{}, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if b.Allow() == nil { allowed <- struct{}{} }
		}()
	}
	wg.Wait()
	close(allowed)

	count := 0
	for range allowed { count++ }
	if count != 1 { t.Fatalf("expected one half-open probe, got %d", count) }
}
