package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	deviceusecase "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/domain/usecase"
)

type ActuatorHandler struct {
	uc *deviceusecase.ActuatorControlUsecase
}

func NewActuatorHandler(uc *deviceusecase.ActuatorControlUsecase) *ActuatorHandler {
	return &ActuatorHandler{uc: uc}
}

type controlRequest struct {
	Action string `json:"action" binding:"required"`
}

func (h *ActuatorHandler) Control(c *gin.Context) {
	// (Optional) we could validate user ownership via device relationship later
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid actuator id"})
		return
	}
	var req controlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.Control(id64, req.Action, "manual"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "actuator updated", "actuator_id": id64, "action": req.Action})
}
