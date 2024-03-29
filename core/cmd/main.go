package main

import (
	"github.com/andibalo/ramein/core"
	"github.com/andibalo/ramein/core/internal/config"
	"github.com/andibalo/ramein/core/internal/db"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.InitConfig()

	database := db.InitDB(cfg)

	server := core.NewServer(cfg, database)

	// Listen from a different goroutine
	go func() {
		if err := server.Listen(cfg.AppAddress()); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c
	cfg.Logger().Info("Gracefully shutting down...")
	_ = server.Shutdown()

	cfg.Logger().Info("Running cleanup tasks...")

	// Your cleanup tasks go here

	database.Close()
	cfg.Logger().Info("Fiber was successful shutdown.")
}
