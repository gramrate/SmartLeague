package series

import (
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/types"
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

	var raw struct {
		Name         string  `json:"name"`
		ScoringRules string  `json:"scoring_rules"`
		StartAt      string  `json:"start_at"`
		EndAt        string  `json:"end_at"`
		Description  *string `json:"description"`
		PriceRub     int     `json:"price_rub"`
		IsClosed     bool    `json:"is_closed"`
		Status       int16   `json:"status"`
	}
	if err := c.Bind(&raw); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	startAt, err := parseDateTimeInput(raw.StartAt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid start_at format"})
	}
	endAt, err := parseDateTimeInput(raw.EndAt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid end_at format"})
	}

	req := dto.CreateSeriesRequest{
		Name:         raw.Name,
		ScoringRules: raw.ScoringRules,
		StartAt:      startAt,
		EndAt:        endAt,
		Description:  raw.Description,
		PriceRub:     raw.PriceRub,
		IsClosed:     raw.IsClosed,
		GameType:     types.GameTypeSportMafia,
		Status:       types.SeriesStatus(raw.Status),
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
