package series

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// CreateSeries Create series
//
// @Summary Create series
// @Tags series
// @Accept json
// @Produce json
// @Param request body dto.CreateSeriesRequest true "Series data"
// @Success 201 {object} dto.CreateSeriesResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/series [post]
func (h *handler) CreateSeries(c echo.Context) error {
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || requesterID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}

	var req dto.CreateSeriesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.seriesService.CreateSeries(c.Request().Context(), requesterID, &req)
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, resp)
}
