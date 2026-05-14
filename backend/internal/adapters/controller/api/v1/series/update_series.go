package series

import (
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/types"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// UpdateSeries Update series by id
//
// @Summary Update series by id
// @Tags series
// @Accept json
// @Produce json
// @Param id path string true "Series ID"
// @Param request body dto.UpdateSeriesRequest true "Update data"
// @Success 200 {object} dto.UpdateSeriesResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/series/{id} [patch]
func (h *handler) UpdateSeries(c echo.Context) error {
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || requesterID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	var raw struct {
		Name         *string `json:"name"`
		ScoringRules *string `json:"scoring_rules"`
		StartAt      *string `json:"start_at"`
		EndAt        *string `json:"end_at"`
		Description  *string `json:"description"`
		PriceRub     *int    `json:"price_rub"`
		IsClosed     *bool   `json:"is_closed"`
		Status       *int16  `json:"status"`
	}
	if err := c.Bind(&raw); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	req := dto.UpdateSeriesRequest{
		ID:           id,
		Name:         raw.Name,
		ScoringRules: raw.ScoringRules,
		Description:  raw.Description,
		PriceRub:     raw.PriceRub,
		IsClosed:     raw.IsClosed,
	}
	if raw.Status != nil {
		status := types.SeriesStatus(*raw.Status)
		req.Status = &status
	}
	if raw.StartAt != nil && *raw.StartAt != "" {
		startAt, parseErr := parseDateTimeInput(*raw.StartAt)
		if parseErr != nil {
			return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid start_at format"})
		}
		req.StartAt = &startAt
	}
	if raw.EndAt != nil && *raw.EndAt != "" {
		endAt, parseErr := parseDateTimeInput(*raw.EndAt)
		if parseErr != nil {
			return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid end_at format"})
		}
		req.EndAt = &endAt
	}
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.seriesService.UpdateSeries(c.Request().Context(), requesterID, &req)
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}
