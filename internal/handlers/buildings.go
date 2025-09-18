package handlers

import (
	"net/http"

	"city/internal/models"
	"city/internal/storage"

	"github.com/gin-gonic/gin"
)

type BuildingsHandler struct{ store *storage.MemoryStore }

func NewBuildingsHandler(s *storage.MemoryStore) *BuildingsHandler {
	return &BuildingsHandler{store: s}
}

func (h *BuildingsHandler) Create(c *gin.Context) {
	var req models.BuildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b, err := h.store.AddBuilding(req.Type)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

func (h *BuildingsHandler) Upgrade(c *gin.Context) {
	id := c.Param("id")
	b, err := h.store.UpgradeBuilding(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BuildingsHandler) Repair(c *gin.Context) {
	id := c.Param("id")
	var req models.RepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Amount == 0 {
		req.Amount = 10
	}
	b, err := h.store.RepairBuilding(id, req.Amount)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BuildingsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.store.RemoveBuilding(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *BuildingsHandler) List(c *gin.Context) {
	city, err := h.store.GetCity()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, city.Buildings)
}

func (h *BuildingsHandler) Effects(c *gin.Context) {
	// basic static description
	effects := []gin.H{
		{"type": models.BuildingFarm, "effect": "+5 food per level per day"},
		{"type": models.BuildingFactory, "effect": "+7 money, -3 energy per level per day"},
		{"type": models.BuildingPowerPlant, "effect": "+10 energy per level per day"},
		{"type": models.BuildingHouse, "effect": "housing/morale, no direct resource"},
	}
	c.JSON(http.StatusOK, effects)
}
