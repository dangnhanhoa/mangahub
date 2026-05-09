package main

import (
	"fmt"

	"mangahub/pkg/utils"
)

func main() {
	cfg := utils.LoadConfig()
	logger := utils.NewLogger(cfg.Logging.Level, cfg.Logging.Path)

	logger.Info("UDP notification server starting", "port", cfg.Server.UDPPort)
	fmt.Printf("UDP notification server ready on :%d\n", cfg.Server.UDPPort)

	// TODO (Dev B): call internal/udp server.Listen()
	select {}
}
