package handlers

import (
	"net/http"

	"city/internal/models"
	"city/internal/storage"

	"github.com/gin-gonic/gin"
)

type ResourcesHandler struct{ store *storage.MemoryStore }

func NewResourcesHandler(s *storage.MemoryStore) *ResourcesHandler {
	return &ResourcesHandler{store: s}
}

func (h *ResourcesHandler) Trade(c *gin.Context) {
	var req models.TradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.store.Trade(req.Resource, req.Amount, req.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ResourcesHandler) History(c *gin.Context) {
	hist := h.store.ResourceHistory()
	c.JSON(http.StatusOK, hist)
}

func (h *ResourcesHandler) Adjust(c *gin.Context) {
	var req models.AdjustResourcesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.store.AdjustResources(req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
