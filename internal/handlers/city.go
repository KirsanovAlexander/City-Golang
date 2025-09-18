package handlers

import (
	"net/http"

	"city/internal/models"
	"city/internal/storage"

	"github.com/gin-gonic/gin"
)

type CityHandler struct {
	store *storage.MemoryStore
}

func NewCityHandler(s *storage.MemoryStore) *CityHandler { return &CityHandler{store: s} }

func (h *CityHandler) Create(c *gin.Context) {
	var req models.CreateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" {
		req.Name = "My City"
	}
	if req.Difficulty == "" {
		req.Difficulty = "normal"
	}
	city := h.store.CreateCity(req.Name, req.Difficulty)
	c.JSON(http.StatusCreated, city)
}

func (h *CityHandler) Get(c *gin.Context) {
	city, err := h.store.GetCity()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, city)
}

func (h *CityHandler) Reset(c *gin.Context) {
	h.store.Reset()
	c.Status(http.StatusNoContent)
}

func (h *CityHandler) UpdateSettings(c *gin.Context) {
	var req models.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.UpdateSettings(req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	city, _ := h.store.GetCity()
	c.JSON(http.StatusOK, city.Settings)
}
