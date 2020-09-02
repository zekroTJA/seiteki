package main

import (
	"log"

	"github.com/op/go-logging"
	"github.com/zekroTJA/seiteki/internal/config"
	"github.com/zekroTJA/seiteki/pkg/seiteki"
	// logging "github.com/op/go-logging"
)

// Version of the seiteki application wrapper.
// This version is independently specified of
// the seiteki package version.
const Version = "1.1.0"

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("failed getting config: %s", err.Error())
	}

	logger := logging.MustGetLogger("main")

	log.SetPrefix("seiteki | ")

	server, err := seiteki.New(cfg)
	if err != nil {
		log.Fatalf("failed creating server: %s", err.Error())
	}
	server.SetLogger(logger)

	if cfg.RouteMode == "" {
		cfg.RouteMode = seiteki.RouteModeRegex
	}

	log.Printf("serving dir: %s", cfg.StaticDir)
	log.Printf("route mode: %s", cfg.RouteMode)
	log.Printf("listening on addr: %s", cfg.Addr)
	err = server.ListenAndServeBlocking()
	log.Fatalf("failed starting server: %s", err.Error())
}
