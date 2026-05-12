package grpc

import (
	"context"
	"log"

	"mangahub/internal/manga"
	"mangahub/pkg/models"
	pb "mangahub/proto"
)

type Server struct {
	pb.UnimplementedMangaServiceServer
	repo *manga.Repository
}

func NewServer(repo *manga.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) GetManga(ctx context.Context, req *pb.GetMangaRequest) (*pb.MangaResponse, error) {
	log.Printf("[gRPC] GetManga request for ID: %s", req.Id)
	m, err := s.repo.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.MangaResponse{
		Id:            m.ID,
		Title:         m.Title,
		Status:        m.Status,
		TotalChapters: int32(m.TotalChapters),
	}, nil
}

func (s *Server) SearchManga(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Printf("[gRPC] SearchManga request for Query: %s", req.Query)
	filters := models.SearchFilters{
		Query:  req.Query,
		Status: req.Status,
	}

	mangas, err := s.repo.List(filters, 50, 0)
	if err != nil {
		return nil, err
	}

	var pbMangas []*pb.MangaResponse
	for _, m := range mangas {
		pbMangas = append(pbMangas, &pb.MangaResponse{
			Id:            m.ID,
			Title:         m.Title,
			Status:        m.Status,
			TotalChapters: int32(m.TotalChapters),
		})
	}

	return &pb.SearchResponse{Mangas: pbMangas}, nil
}
