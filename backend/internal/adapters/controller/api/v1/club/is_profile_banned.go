package club

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// IsProfileBanned Check if a profile is banned in a club
//
// @Summary Check if profile is banned in club
// @Tags club
// @Produce json
// @Security CookieAuth
// @Param id path string true "Club ID"
// @Param profile_id path string true "Profile ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/club/{id}/profile/{profile_id}/is-banned [get]
func (h *handler) IsProfileBanned(c echo.Context) error {
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || requesterID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}
	clubID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid club id"})
	}
	profileID, err := uuid.Parse(c.Param("profile_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid profile id"})
	}
	banned, err := h.clubService.IsProfileBanned(c.Request().Context(), requesterID, clubID, profileID)
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]bool{"is_banned": banned})
}
