package handlers

import (
	"net/http"
	"template-golang/modules/cockroach/models"
	"template-golang/modules/cockroach/usecases"
	"template-golang/modules/request"

	"github.com/labstack/echo/v4"
)

type cockroachHttpHandler struct {
	cockroachUsecase usecases.CockroachUsecase
}

func NewCockroachHttpHandler(cockroachUsecase usecases.CockroachUsecase) CockroachHandler {
	return &cockroachHttpHandler{
		cockroachUsecase: cockroachUsecase,
	}
}

func (h *cockroachHttpHandler) DetectCockroach(c echo.Context) error {
	reqBody := new(models.AddCockroachData)

	wrapper := request.ContextWrapper(c)

	if err := wrapper.Bind(reqBody); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": err.Error()},
		)
	}

	if err := h.cockroachUsecase.CockroachDataProcessing(reqBody); err != nil {
		return response(c, http.StatusInternalServerError, "Processing data failed")
	}

	return response(c, http.StatusOK, "Success 🪳🪳🪳")
}
