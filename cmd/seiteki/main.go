package main

import (
	"log"

	"github.com/zekroTJA/seiteki/internal/config"
	"github.com/zekroTJA/seiteki/pkg/seiteki"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("failed getting config: %s", err.Error())
	}

	server, err := seiteki.New(cfg)
	if err != nil {
		log.Fatalf("failed creating server: %s", err.Error())
	}

	log.Printf("serving dir: %s", cfg.StaticDir)
	log.Printf("listening on addr: %s", cfg.Addr)
	err = server.ListenAndServeBlocking()
	log.Fatalf("failed starting server: %s", err.Error())
}
