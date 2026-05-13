package profile

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Delete current profile
//
// @Summary Delete current profile
// @Tags profile
// @Produce json
// @Success 204
// @Failure 401 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/profile [delete]
func (h *handler) DeleteMe(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}

	if err := h.profileService.Delete(c.Request().Context(), &dto.DeleteProfileRequest{ID: userID}); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

