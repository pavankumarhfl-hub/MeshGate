package breaker

import (
	"testing"
	"time"
)

func TestBreakerOpensAfterThreshold(t *testing.T) {
	b := New(2, time.Second)
	b.Failure(); b.Failure()
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
