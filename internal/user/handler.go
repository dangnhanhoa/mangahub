package user
import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"mangahub/pkg/models"
)

type Broadcaster interface {
	Broadcast(data interface{})
}
type Handler struct {
	repo        *Repository
	broadcaster Broadcaster
}
func NewHandler(repo *Repository, b Broadcaster) *Handler {
	return &Handler{
		repo:        repo,
		broadcaster: b,
	}
}
func (h *Handler) AddToLibrary(c *gin.Context) {
	userID := c.MustGet("user_id").(string) 
	var req models.UserProgress
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	req.UserID = userID
	if req.Status == "" {
		req.Status = "reading"
	}
	if err := h.repo.AddOrUpdate(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "add to library success", "data": req})
}

func (h *Handler) ListLibrary(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	status := c.Query("status") 
	list, err := h.repo.List(userID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi truy xuất dữ liệu"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *Handler) UpdateProgress(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	var req struct {
		MangaID string `json:"manga_id" binding:"required"`
		Chapter int    `json:"chapter" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "thiếu manga_id hoặc chapter"})
		return
	}
	progress, err := h.repo.Get(userID, req.MangaID)
	if err != nil {
		progress = &models.UserProgress{
			UserID:  userID,
			MangaID: req.MangaID,
			Status:  "reading",
		}
	}
	progress.CurrentChapter = req.Chapter
	if err := h.repo.AddOrUpdate(progress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi lưu tiến độ"})
		return
	}
	
	if h.broadcaster != nil {
		updateMsg := models.ProgressUpdate{
			UserID:    userID,
			Username:  "User_" + userID[:5], 
			MangaID:   req.MangaID,
			Chapter:   req.Chapter,
			Timestamp: time.Now().Unix(),
		}
		
		go h.broadcaster.Broadcast(updateMsg)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Update progress success",
		"data":    progress,
	})
}