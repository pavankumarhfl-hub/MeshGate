package proxy

import (
	"context"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/pavankumarhfl-hub/MeshGate/internal/breaker"
)

type Proxy struct {
	client *http.Client
	upstream *url.URL
	breaker *breaker.Breaker
	maxRetries int
}

func New(target string, timeout time.Duration, retries int) (*Proxy, error) {
	u, err := url.Parse(target); if err != nil { return nil, err }
	return &Proxy{client: &http.Client{Timeout: timeout}, upstream: u, breaker: breaker.New(5, 5*time.Second), maxRetries: retries}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := p.breaker.Allow(); err != nil { http.Error(w, err.Error(), http.StatusServiceUnavailable); return }
	var last error
	for attempt := 0; attempt <= p.maxRetries; attempt++ {
		if attempt > 0 { time.Sleep(time.Duration(25*(1<<min(attempt, 6)))*time.Millisecond + time.Duration(rand.Intn(25))*time.Millisecond) }
		err := p.forward(r.Context(), w, r)
		if err == nil { p.breaker.Success(); return }
		last = err
		p.breaker.Failure()
		if !retryable(err) { break }
	}
	http.Error(w, "upstream unavailable: "+last.Error(), http.StatusBadGateway)
}

func (p *Proxy) forward(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	u := *p.upstream
	u.Path = joinPath(u.Path, r.URL.Path)
	u.RawQuery = r.URL.RawQuery
	req, err := http.NewRequestWithContext(ctx, r.Method, u.String(), r.Body); if err != nil { return err }
	req.Header = r.Header.Clone()
	resp, err := p.client.Do(req); if err != nil { return err }
	defer resp.Body.Close()
	for k, values := range resp.Header { for _, v := range values { w.Header().Add(k, v) } }
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil { return err }
	if resp.StatusCode >= 500 { return &upstreamStatus{code: resp.StatusCode} }
	return nil
}

type upstreamStatus struct { code int }
func (e *upstreamStatus) Error() string { return http.StatusText(e.code) }
func retryable(err error) bool { if e, ok := err.(*upstreamStatus); ok { return e.code >= 500 }; return true }
func joinPath(a,b string) string { if a=="/" { return b }; if b=="/" { return a }; return a+"/"+b }
func min(a,b int) int { if a<b{return a}; return b }
