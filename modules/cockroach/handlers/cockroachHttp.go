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

// @BasePath /api/v1

// DetectCockroach godoc
// @Summary Detect if image contains cockroach
// @Schemes
// @Description Analyzes image to detect presence of cockroach
// @Tags cockroach
// @Accept json
// @Produce json
// @Param request body models.AddCockroachData true "Request body"
// @Success 200 {object} map[string]interface{} "Success response with message"
// @Router /cockroach [post]
func (h *cockroachHttpHandler) DetectCockroach(c *gin.Context) {
	reqBody := new(models.AddCockroachData)

	if err := c.ShouldBindJSON(reqBody); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": err.Error()},
		)
		c.Error(err)
		return
	}

	if err := h.cockroachUsecase.CockroachDataProcessing(reqBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Processing data failed"})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success 🪳🪳🪳"})
	return
}
