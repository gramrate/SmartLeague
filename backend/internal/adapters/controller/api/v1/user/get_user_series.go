package user

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GetUserSeries Get user series
//
// @Summary Get user series
// @Tags user
// @Produce json
// @Param id path string true "User ID"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Param q query string false "search by series name"
// @Param from query string false "filter start boundary, YYYY-MM-DD"
// @Param to query string false "filter end boundary, YYYY-MM-DD"
// @Param is_rating query boolean false "filter by rating flag (true/false)"
// @Param show_past query boolean false "show past series (end_at < now)"
// @Param show_closed query boolean false "show closed registration series (is_closed = true)"
// @Success 200 {object} dto.GetUserSeriesResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/user/{id}/series [get]
func (h *handler) GetUserSeries(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	var req dto.GetUserSeriesRequest
	if err := h.formDecoder.Decode(&req, c.QueryParams()); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	req.UserID = userID

	if viewerID, ok := c.Get("user_id").(uuid.UUID); ok && viewerID != uuid.Nil {
		req.ViewerID = &viewerID
	}

	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.userService.GetUserSeries(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}
