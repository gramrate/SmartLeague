package series

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GetClubSeries Get club series list (paginated)
//
// @Summary Get club series list
// @Tags series
// @Produce json
// @Param id path string true "Club ID"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {object} dto.GetClubSeriesResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/club/{id}/series [get]
func (h *handler) GetClubSeries(c echo.Context) error {
	clubID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	var req dto.GetClubSeriesRequest
	if err := h.formDecoder.Decode(&req, c.QueryParams()); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	req.ClubID = clubID
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.seriesService.GetClubSeries(c.Request().Context(), maybeRequesterID(c.Get("user_id")), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}
