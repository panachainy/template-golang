package handlers

import (
	"net/http"
	"template-golang/modules/cockroach/models"
	"template-golang/modules/cockroach/usecases"

	"github.com/gin-gonic/gin"
)

type cockroachHttpHandler struct {
	cockroachUsecase usecases.CockroachUsecase
}

func NewCockroachHttpHandler(cockroachUsecase usecases.CockroachUsecase) CockroachHandler {
	return &cockroachHttpHandler{
		cockroachUsecase: cockroachUsecase,
	}
}

func (h *cockroachHttpHandler) DetectCockroach(c gin.Context) error {
	reqBody := new(models.AddCockroachData)

	if err := c.ShouldBindJSON(reqBody); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": err.Error()},
		)
		return err
	}

	if err := h.cockroachUsecase.CockroachDataProcessing(reqBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Processing data failed"})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success 🪳🪳🪳"})
	return nil
}
