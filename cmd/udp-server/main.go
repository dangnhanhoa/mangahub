package main

import (
	"fmt"
	"log"

	"mangahub/internal/udp"
	"mangahub/pkg/utils"
)

func main() {
	cfg := utils.LoadConfig()

	server := udp.NewServer(cfg.Server.UDPPort)
	if err := server.Start(); err != nil {
		log.Fatalf("[UDP] Failed to start server: %v", err)
	} 
	fmt.Printf("UDP notification server ready on :%d\n", cfg.Server.UDPPort)

}
