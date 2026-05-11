package manga

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mangahub/pkg/models"
)

type Handler struct{
	repo *Repository
}

func NewHandler(repo *Repository)*Handler{
	return &Handler{repo: repo}
}

func (h *Handler) GetMangaByID(c *gin.Context){
	id := c.Param("id")
	
	manga, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error" : "Manga not found",
		})
		return
	}

	c.JSON(http.StatusOK, manga)
}

func (h *Handler) List(c *gin.Context) {
	var filter models.SearchFilters

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Invalid Filter",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {page = 1 }
	if limit < 1 || limit > 50 {limit =20 }

	offset := (page - 1) * limit;

	mangas, err := h.repo.List(filter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Fail to retrive data",
		})
		return 
	}

	c.JSON(http.StatusOK, gin.H{
		"data" : mangas,
		"meta" : gin.H{
			"page" : page,
			"limit" : limit, 
		},
	})
}