package main

import (
	"fmt"
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
		if err := server.Listen(":8000"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c
	fmt.Println("Gracefully shutting down...")
	_ = server.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()
	fmt.Println("Fiber was successful shutdown.")
}
