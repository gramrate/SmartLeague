package series

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GetPayments Get paid participants in series (leaders/president only)
//
// @Summary Get series payments
// @Tags series
// @Produce json
// @Security CookieAuth
// @Param id path string true "Series ID"
// @Success 200 {object} dto.GetSeriesPaymentsResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/series/{id}/payments [get]
func (h *handler) GetPayments(c echo.Context) error {
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || requesterID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}
	seriesID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}
	resp, err := h.seriesService.GetPayments(c.Request().Context(), requesterID, &dto.GetSeriesPaymentsRequest{SeriesID: seriesID})
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

// SetPayment Mark/unmark participant payment in paid series (leaders/president only)
//
// @Summary Set participant payment status
// @Tags series
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path string true "Series ID"
// @Param profile_id path string true "Profile ID"
// @Param request body object{paid=bool} true "Payment status"
// @Success 204
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/series/{id}/payment/{profile_id} [post]
func (h *handler) SetPayment(c echo.Context) error {
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || requesterID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}
	seriesID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}
	profileID, err := uuid.Parse(c.Param("profile_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid profile id"})
	}
	var body struct {
		Paid bool `json:"paid"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	if err := h.seriesService.SetPayment(c.Request().Context(), requesterID, &dto.SetSeriesPaymentRequest{
		SeriesID:  seriesID,
		ProfileID: profileID,
		Paid:      body.Paid,
	}); err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
