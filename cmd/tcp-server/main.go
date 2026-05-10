package main

import (
	"fmt"
	"log"

	"mangahub/pkg/database"
	"mangahub/pkg/utils"
	"mangahub/internal/tcp"
)

func main() {
	cfg := utils.LoadConfig()
	logger := utils.NewLogger(cfg.Logging.Level, cfg.Logging.Path)

	db, err := database.New(cfg.Database.Path)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	server := tcp.NewServer(cfg.Server.TCPPort)
	if err := server.Start(); err != nil {
		log.Fatalf("TCP Server: %v", err)
	}

	logger.Info("TCP sync server starting", "port", cfg.Server.TCPPort)
	fmt.Printf("TCP sync server ready on :%d\n", cfg.Server.TCPPort)

	// TODO (Dev B): call internal/tcp server.Listen()
	select {}
}
