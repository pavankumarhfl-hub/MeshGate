package ratelimit

import "testing"

func TestLimiterRejectsAfterBurst(t *testing.T) {
	l := New(2, 2)
	if !l.Allow("client") || !l.Allow("client") { t.Fatal("expected initial requests to pass") }
	if l.Allow("client") { t.Fatal("expected request to be rejected") }
}
