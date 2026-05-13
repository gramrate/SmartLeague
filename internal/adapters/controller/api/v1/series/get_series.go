package series

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *handler) GetSeries(c echo.Context) error {
	// Get series by id
	//
	// @Summary Get series by id
	// @Tags series
	// @Produce json
	// @Param id path string true "Series ID"
	// @Success 200 {object} dto.GetSeriesResponse
	// @Failure 400 {object} dto.HTTPStatus
	// @Failure 403 {object} dto.HTTPStatus
	// @Failure 500 {object} dto.HTTPStatus
	// @Router /api/v1/series/{id} [get]
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	resp, err := h.seriesService.GetSeries(c.Request().Context(), maybeRequesterID(c.Get("user_id")), &dto.GetSeriesRequest{ID: id})
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}
