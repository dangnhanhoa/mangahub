package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"mangahub/internal/auth"
	"mangahub/internal/manga"

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
	
	//Auth service init
	router := gin.Default()
	authService := auth.NewService(db)
	authHandler := auth.NewHandler(authService)
	
	// Auth Routes
	authRoutes := router.Group("/auth")
	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)

	authRoutes.POST("/logout", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message": "Logout success",
		})

	})

	// Protected Route
	protectedRoute := router.Group("/protected")
	protectedRoute.Use(auth.AuthMiddleware())
	{
		protectedRoute.GET("/test",func(c *gin.Context){
			userID := c.GetString("user_id")

			c.JSON(http.StatusOK, gin.H{
				"message": "success",
				"user_id": userID,
			})
		})
	}

	//Manga service init
	mangaRepo := manga.NewRepository(db)
	mangaHandler := manga.NewHandler(mangaRepo)

	// Manga Routes
	mangaRoutes := router.Group("/manga")
	mangaRoutes.GET("", mangaHandler.List)
	mangaRoutes.GET("/:id", mangaHandler.GetMangaByID)

	port := cfg.Server.HTTPPort
	addr := fmt.Sprintf(":%d",port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
	select {}
}
