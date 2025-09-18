package handlers

import (
	"fmt"
	"net/http"

	"city/internal/models"
	"city/internal/storage"

	"github.com/gin-gonic/gin"
)

type CitizensHandler struct{ store *storage.MemoryStore }

func NewCitizensHandler(s *storage.MemoryStore) *CitizensHandler { return &CitizensHandler{store: s} }

func (h *CitizensHandler) Create(c *gin.Context) {
	var req models.CitizenCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" {
		req.Name = "Citizen"
	}
	cz, err := h.store.AddCitizen(req.Name, req.Job)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cz)
}

func (h *CitizensHandler) List(c *gin.Context) {
	city, err := h.store.GetCity()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, city.Citizens)
}

func (h *CitizensHandler) ChangeJob(c *gin.Context) {
	id := c.Param("id")
	var req models.ChangeJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cz, err := h.store.ChangeCitizenJob(id, req.Job)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cz)
}

func (h *CitizensHandler) ChangeHappiness(c *gin.Context) {
	id := c.Param("id")
	var req models.ChangeHappinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cz, err := h.store.AdjustCitizenHappiness(id, req.Delta)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cz)
}

func (h *CitizensHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.store.RemoveCitizen(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CitizensHandler) JobsStats(c *gin.Context) {
	city, err := h.store.GetCity()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	counters := map[models.Job]int{}
	for _, cz := range city.Citizens {
		counters[cz.Job]++
	}
	c.JSON(http.StatusOK, counters)
}

func (h *CitizensHandler) MassAdd(c *gin.Context) {
	var req models.MassAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Count <= 0 {
		req.Count = 1
	}
	if req.Prefix == "" {
		req.Prefix = "Citizen"
	}
	added := make([]models.Citizen, 0, req.Count)
	for i := 1; i <= req.Count; i++ {
		name := fmt.Sprintf("%s_%d", req.Prefix, i)
		cz, err := h.store.AddCitizen(name, req.Job)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		added = append(added, cz)
	}
	c.JSON(http.StatusCreated, added)
}
