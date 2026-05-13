package profile

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Update current profile
//
// @Summary Update current profile
// @Tags profile
// @Accept json
// @Produce json
// @Param request body dto.UpdateCurrentProfileRequest true "Update data"
// @Success 200 {object} dto.UpdateCurrentProfileResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/profile [patch]
func (h *handler) UpdateMe(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}

	var req dto.UpdateCurrentProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	req.ID = userID

	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.profileService.UpdateCurrent(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

