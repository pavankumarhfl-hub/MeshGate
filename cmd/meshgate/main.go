package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pavankumarhfl-hub/MeshGate/internal/config"
	"github.com/pavankumarhfl-hub/MeshGate/internal/gateway"
)

func main() {
	cfg := config.Default()
	if addr := os.Getenv("MESHGATE_ADDR"); addr != "" {
		cfg.ListenAddr = addr
	}

	gw := gateway.New(cfg)
	server := &http.Server{Addr: cfg.ListenAddr, Handler: gw.Handler(), ReadHeaderTimeout: cfg.ReadHeaderTimeout}

	log.Printf("meshgate listening on %s", cfg.ListenAddr)
	if err := gw.Run(server); err != nil {
		log.Fatal(err)
	}
}
