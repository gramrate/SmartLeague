package series

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAllSeries Get all series list (paginated)
//
// @Summary Get all series list
// @Tags series
// @Produce json
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Param show_past query boolean false "show past series (end_at < now)"
// @Param show_closed query boolean false "show closed registration series (is_closed = true)"
// @Success 200 {object} dto.GetAllSeriesResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/series/all [get]
func (h *handler) GetAllSeries(c echo.Context) error {
	var req dto.GetAllSeriesRequest
	if err := h.formDecoder.Decode(&req, c.QueryParams()); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.seriesService.GetAllSeries(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}
