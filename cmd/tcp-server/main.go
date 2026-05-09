package main

import (
	"fmt"
	"log"

	"mangahub/pkg/database"
	"mangahub/pkg/utils"
)

func main() {
	cfg := utils.LoadConfig()
	logger := utils.NewLogger(cfg.Logging.Level, cfg.Logging.Path)

	db, err := database.New(cfg.Database.Path)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	logger.Info("TCP sync server starting", "port", cfg.Server.TCPPort)
	fmt.Printf("TCP sync server ready on :%d\n", cfg.Server.TCPPort)

	// TODO (Dev B): call internal/tcp server.Listen()
	select {}
}
