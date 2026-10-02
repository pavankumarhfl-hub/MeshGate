package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProxyForwardsRequest(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/hello" { t.Fatalf("unexpected path: %s", r.URL.Path) }
		w.Header().Set("X-Upstream", "ok")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello"))
	}))
	defer upstream.Close()

	p, err := New(upstream.URL, time.Second, 1)
	if err != nil { t.Fatal(err) }
	r := httptest.NewRequest(http.MethodGet, "http://gateway/hello", nil)
	rec := httptest.NewRecorder()
	p.ServeHTTP(rec, r)
	if rec.Code != http.StatusOK || rec.Body.String() != "hello" { t.Fatalf("unexpected response: %d %q", rec.Code, rec.Body.String()) }
}

func TestProxyRetriesServerErrors(t *testing.T) {
	calls := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		if calls < 3 { http.Error(w, "temporary", http.StatusBadGateway); return }
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	p, err := New(upstream.URL, time.Second, 2)
	if err != nil { t.Fatal(err) }
	rec := httptest.NewRecorder()
	p.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "http://gateway/", nil))
	if rec.Code != http.StatusOK || calls != 3 { t.Fatalf("expected successful third attempt, status=%d calls=%d", rec.Code, calls) }
}
