package handlers

import (
	"net/http"
	"template-golang/modules/cockroach/models"
	"template-golang/modules/cockroach/usecases"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type cockroachHttpHandler struct {
	cockroachUsecase usecases.CockroachUsecase
}

func NewCockroachHttpHandler(cockroachUsecase usecases.CockroachUsecase) *cockroachHttpHandler {
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

	validate := validator.New(validator.WithRequiredStructEnabled())

	// Validate the request body
	if err := validate.Struct(reqBody); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": err.Error()},
		)
		c.Error(err)
		return
	}

	if err := h.cockroachUsecase.ProcessData(reqBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Processing data failed"})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success 🪳🪳🪳"})
	return
}
