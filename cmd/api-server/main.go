package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"mangahub/internal/auth"
	"mangahub/internal/manga"
	"mangahub/internal/user"
	"mangahub/internal/websocket"
	"mangahub/internal/tcp"
	"mangahub/internal/udp"

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

	// User service init
	tcpAddr := fmt.Sprintf("localhost:%d", cfg.Server.TCPPort)
	if cfg.Server.TCPPort == 0 {
		tcpAddr = "localhost:9090"
	}
	tcpClient := tcp.NewClient(tcpAddr)
	
	userRepo := user.NewRepository(db)
	userHandler := user.NewHandler(userRepo, tcpClient)

	//User Routes
	userRoutes := router.Group("/user")
	userRoutes.Use(auth.AuthMiddleware())
	{
		userRoutes.POST("/library", userHandler.AddToLibrary)
		userRoutes.GET("/library", userHandler.ListLibrary)
		userRoutes.PUT("/progress", userHandler.UpdateProgress)
	}

	//websocket
	wsHub := websocket.NewHub()
	wsHandler := websocket.NewHandler(wsHub)

	chatRoutes := router.Group("/chat/:mangaId")
	chatRoutes.Use(auth.AuthMiddleware())
	{
		chatRoutes.GET("/ws",wsHandler.ServeWS)
	}

	// UDP Trigger API
	udpAddr := fmt.Sprintf("localhost:%d", cfg.Server.UDPPort)
	if cfg.Server.UDPPort == 0 {
		udpAddr = "localhost:9091"
	}
	udpClient := udp.NewClient(udpAddr)

	adminRoutes := router.Group("/admin")
	adminRoutes.Use(auth.AuthMiddleware())
	{
		adminRoutes.POST("/notify", func(c *gin.Context) {
			var req struct {
				MangaID string `json:"manga_id"`
				Message string `json:"message"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "bad request"})
				return
			}
			udpClient.TriggerBroadcast(req.MangaID, req.Message)
			c.JSON(200, gin.H{"message": "notification triggered"})
		})
	}

	port := cfg.Server.HTTPPort
	addr := fmt.Sprintf(":%d",port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
	select {}
}
