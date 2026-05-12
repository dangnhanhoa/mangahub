package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	grpcserver "mangahub/internal/grpc"
	"mangahub/internal/manga"
	"mangahub/pkg/database"
	"mangahub/pkg/utils"
	pb "mangahub/proto"
)

func main() {
	cfg := utils.LoadConfig()

	db, err := database.New(cfg.Database.Path)
	if err != nil {
		log.Fatalf("[gRPC] Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := manga.NewRepository(db)
	server := grpcserver.NewServer(repo)

	port := cfg.Server.GRPCPort
	if port == 0 {
		port = 9092
	}
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[gRPC] Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMangaServiceServer(s, server)

	log.Printf("[gRPC] Server is listening on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[gRPC] Failed to serve: %v", err)
	}
}
