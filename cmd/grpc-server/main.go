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

	logger.Info("gRPC server starting", "port", cfg.Server.GRPCPort)
	fmt.Printf("gRPC server ready on :%d\n", cfg.Server.GRPCPort)

	// TODO (Dev A): register gRPC services and call grpcServer.Serve()
	select {}
}
