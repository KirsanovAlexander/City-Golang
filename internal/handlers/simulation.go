package handlers

import (
	"net/http"

	"city/internal/models"
	"city/internal/storage"

	"github.com/gin-gonic/gin"
)

type SimulationHandler struct{ store *storage.MemoryStore }

func NewSimulationHandler(s *storage.MemoryStore) *SimulationHandler {
	return &SimulationHandler{store: s}
}

func (h *SimulationHandler) Tick(c *gin.Context) {
	city, err := h.store.Tick()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, city)
}

func (h *SimulationHandler) RandomEvent(c *gin.Context) {
	ev, err := h.store.RandomEvent()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ev)
}

func (h *SimulationHandler) CustomEvent(c *gin.Context) {
	var req models.CustomEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ev, err := h.store.CustomEvent(req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ev)
}

func (h *SimulationHandler) EventsHistory(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.EventsHistory())
}

func (h *SimulationHandler) Stats(c *gin.Context) {
	st, err := h.store.Stats()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, st)
}
