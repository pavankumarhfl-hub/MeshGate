package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ListenAddr        string
	ReadHeaderTimeout time.Duration
	RequestTimeout    time.Duration
	UpstreamTimeout   time.Duration
	MaxRetries        int
	RateLimit         int
	Burst             int
	MaxBodyBytes      int64
	ShutdownTimeout   time.Duration
	Routes            map[string]string
}

func Default() Config {
	return Config{
		ListenAddr:        ":8080",
		ReadHeaderTimeout: 5 * time.Second,
		RequestTimeout:    10 * time.Second,
		UpstreamTimeout:   3 * time.Second,
		MaxRetries:        2,
		RateLimit:         100,
		Burst:             100,
		MaxBodyBytes:      1 << 20,
		ShutdownTimeout:   10 * time.Second,
		Routes:            map[string]string{},
	}
}

func FromEnv() Config {
	cfg := Default()
	if v := os.Getenv("MESHGATE_ADDR"); v != "" { cfg.ListenAddr = v }
	if v := os.Getenv("MESHGATE_UPSTREAM_TIMEOUT"); v != "" { if d, err := time.ParseDuration(v); err == nil { cfg.UpstreamTimeout = d } }
	if v := os.Getenv("MESHGATE_REQUEST_TIMEOUT"); v != "" { if d, err := time.ParseDuration(v); err == nil { cfg.RequestTimeout = d } }
	if v := os.Getenv("MESHGATE_MAX_RETRIES"); v != "" { if n, err := strconv.Atoi(v); err == nil && n >= 0 { cfg.MaxRetries = n } }
	if v := os.Getenv("MESHGATE_RATE_LIMIT"); v != "" { if n, err := strconv.Atoi(v); err == nil && n > 0 { cfg.RateLimit = n } }
	if v := os.Getenv("MESHGATE_BURST"); v != "" { if n, err := strconv.Atoi(v); err == nil && n > 0 { cfg.Burst = n } }
	if v := os.Getenv("MESHGATE_MAX_BODY_BYTES"); v != "" { if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 { cfg.MaxBodyBytes = n } }
	if path, target := os.Getenv("MESHGATE_ROUTE_PATH"), os.Getenv("MESHGATE_ROUTE_TARGET"); path != "" && target != "" {
		cfg.Routes[path] = target
	}
	return cfg
}
