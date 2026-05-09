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

	logger.Info("HTTP API server starting", "port", cfg.Server.HTTPPort)
	fmt.Printf("HTTP API server ready on :%d\n", cfg.Server.HTTPPort)

	// TODO (Dev A): register Gin routes and call router.Run()
	select {}
}
