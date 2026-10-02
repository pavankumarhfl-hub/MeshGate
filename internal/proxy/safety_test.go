package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProxyDoesNotRetryPOST(t *testing.T) {
	calls := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		http.Error(w, "temporary", http.StatusBadGateway)
	}))
	defer upstream.Close()

	p, err := New(upstream.URL, time.Second, 3)
	if err != nil { t.Fatal(err) }
	rec := httptest.NewRecorder()
	p.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "http://gateway/", nil))
	if calls != 1 { t.Fatalf("POST was retried: calls=%d", calls) }
}
