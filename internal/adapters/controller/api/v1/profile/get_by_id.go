package profile

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GetByID Get profile by id
//
// @Summary Get profile by id
// @Tags profile
// @Produce json
// @Param id path string true "Profile ID"
// @Success 200 {object} dto.GetProfileResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/profile/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	resp, err := h.profileService.GetByID(c.Request().Context(), &dto.GetProfileRequest{ID: id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}
