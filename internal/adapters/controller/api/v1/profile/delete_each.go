package profile

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// DeleteEach Delete profile by id (admin)
//
// @Summary Delete profile by id
// @Tags profile
// @Produce json
// @Param id path string true "Profile ID"
// @Success 204
// @Failure 400 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/profile/{id} [delete]
func (h *handler) DeleteEach(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	if err := h.profileService.Delete(c.Request().Context(), &dto.DeleteProfileRequest{ID: id}); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
