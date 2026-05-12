package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mangahub/pkg/utils"
	pb "mangahub/proto"
)

func main() {
	cfg := utils.LoadConfig()
	port := cfg.Server.GRPCPort
	if port == 0 {
		port = 9092
	}
	addr := fmt.Sprintf("localhost:%d", port)

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getID := getCmd.String("id", "", "")

	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	searchQuery := searchCmd.String("query", "", "")
	searchStatus := searchCmd.String("status", "", "")

	if len(os.Args) < 2 {
		fmt.Println("Usage: cli [get|search] [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get":
		getCmd.Parse(os.Args[2:])
		if *getID == "" {
			fmt.Println("Error: -id is required")
			os.Exit(1)
		}
		runGet(addr, *getID)
	case "search":
		searchCmd.Parse(os.Args[2:])
		runSearch(addr, *searchQuery, *searchStatus)
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}

func getClient(addr string) (*grpc.ClientConn, pb.MangaServiceClient) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	return conn, pb.NewMangaServiceClient(conn)
}

func runGet(addr, id string) {
	conn, client := getClient(addr)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.GetManga(ctx, &pb.GetMangaRequest{Id: id})
	if err != nil {
		log.Fatalf("Error calling GetManga: %v", err)
	}

	fmt.Printf("Manga ID: %s\nTitle: %s\nStatus: %s\nChapters: %d\n", res.Id, res.Title, res.Status, res.TotalChapters)
}

func runSearch(addr, query, status string) {
	conn, client := getClient(addr)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.SearchManga(ctx, &pb.SearchRequest{
		Query:  query,
		Status: status,
	})
	if err != nil {
		log.Fatalf("Error calling SearchManga: %v", err)
	}

	fmt.Printf("Found %d results:\n", len(res.Mangas))
	for _, m := range res.Mangas {
		fmt.Printf("- [%s] %s (Status: %s, Chapters: %d)\n", m.Id, m.Title, m.Status, m.TotalChapters)
	}
}
