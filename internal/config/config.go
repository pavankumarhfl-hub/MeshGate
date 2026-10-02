package config

import "time"

type Config struct {
	ListenAddr       string
	ReadHeaderTimeout time.Duration
	RequestTimeout   time.Duration
	UpstreamTimeout  time.Duration
	MaxRetries       int
	RateLimit        int
	Burst            int
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
	}
}
