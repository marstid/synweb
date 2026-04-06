package main

import (
	"context"
	"os"

	"github.com/marstid/synweb/internal/config"
	"github.com/marstid/synweb/internal/logger"
	"github.com/marstid/synweb/internal/search"
	"github.com/marstid/synweb/internal/server"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		println("ERROR: Failed to load configuration")
		println(err.Error())
		os.Exit(1)
	}

	if err := config.Validate(cfg); err != nil {
		println("ERROR: " + err.Error())
		os.Exit(1)
	}

	log := logger.New(cfg.LogLevel)

	log.Info("Starting synweb MCP server", "version", "0.0.1")

	client := search.NewClient(cfg.APIBaseURL, cfg.SyntheticAPIKey, log)
	handler := search.NewHandler(client, log)
	srv := server.New("synweb", "0.0.1", log, handler)

	if err := srv.Run(ctx); err != nil {
		log.Error("Server error", "error", err)
		os.Exit(1)
	}
}
