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

func TestProxyRejectsOversizedBody(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	p, err := New(upstream.URL, time.Second, 0)
	if err != nil { t.Fatal(err) }
	req := httptest.NewRequest(http.MethodPost, "http://gateway/", http.NoBody)
	req.Body = http.MaxBytesReader(httptest.NewRecorder(), req.Body, 1)
	rec := httptest.NewRecorder()
	p.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK { t.Fatalf("unexpected status: %d", rec.Code) }
}
