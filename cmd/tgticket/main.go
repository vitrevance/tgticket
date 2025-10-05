package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/vitrevance/tgticket/internal/config"
	"github.com/vitrevance/tgticket/internal/server"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the config file")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := server.NewServer(cfg)
	srv.RegisterRoutes()

	log.Printf("Listening on %s...", cfg.ServerAddr)
	log.Fatal(http.ListenAndServe(cfg.ServerAddr, nil))
}
