package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/pavankumarhfl-hub/MeshGate/internal/config"
	"github.com/pavankumarhfl-hub/MeshGate/internal/proxy"
	"github.com/pavankumarhfl-hub/MeshGate/internal/ratelimit"
)

type Gateway struct { cfg config.Config; limiter *ratelimit.Limiter; mu sync.RWMutex; routes map[string]http.Handler }

func New(cfg config.Config) *Gateway { return &Gateway{cfg: cfg, limiter: ratelimit.New(cfg.RateLimit,cfg.Burst), routes: make(map[string]http.Handler)} }

func (g *Gateway) AddRoute(path, target string) error { p, err := proxy.New(target,g.cfg.UpstreamTimeout,g.cfg.MaxRetries); if err != nil{return err}; g.mu.Lock(); g.routes[path]=p; g.mu.Unlock(); return nil }

func (g *Gateway) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(http.StatusOK); _,_=w.Write([]byte("ok\n"))})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(http.StatusOK); _,_=w.Write([]byte("ready\n"))})
	mux.HandleFunc("/routes", g.routeHandler)
	mux.HandleFunc("/", g.dispatch)
	return g.logging(g.limit(mux))
}

func (g *Gateway) dispatch(w http.ResponseWriter,r *http.Request) { g.mu.RLock(); h:=g.routes[r.URL.Path]; g.mu.RUnlock(); if h==nil {http.Error(w,"route not found",http.StatusNotFound);return}; h.ServeHTTP(w,r) }
func (g *Gateway) routeHandler(w http.ResponseWriter,_ *http.Request){ g.mu.RLock(); defer g.mu.RUnlock(); out:=make([]string,0,len(g.routes)); for p:=range g.routes{out=append(out,p)}; w.Header().Set("Content-Type","application/json"); _=json.NewEncoder(w).Encode(out) }
func (g *Gateway) limit(next http.Handler) http.Handler{return ratelimit.Middleware(g.limiter,next)}
func (g *Gateway) logging(next http.Handler) http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){start:=time.Now(); next.ServeHTTP(w,r); _=start})}
func (g *Gateway) Run(s *http.Server) error { if err:=g.AddRoute("/example","http://127.0.0.1:9000"); err!=nil{return err}; errCh:=make(chan error,1); go func(){errCh<-s.ListenAndServe()}(); return <-errCh }
func shutdown(ctx context.Context,s *http.Server) error{return s.Shutdown(ctx)}
