package series

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// DeleteSeries Delete series by id
//
// @Summary Delete series by id
// @Tags series
// @Produce json
// @Param id path string true "Series ID"
// @Success 204
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/series/{id} [delete]
func (h *handler) DeleteSeries(c echo.Context) error {
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || requesterID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	if err := h.seriesService.DeleteSeries(c.Request().Context(), requesterID, &dto.DeleteSeriesRequest{ID: id}); err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
