package club

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// SetLeader Set club leader for member
//
// @Summary Set club leader
// @Tags club
// @Produce json
// @Param id path string true "Club ID"
// @Param member_id path string true "Member Profile ID"
// @Success 204
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/club/{id}/leader/{member_id} [post]
func (h *handler) SetLeader(c echo.Context) error {
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || requesterID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}

	clubID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid club id"})
	}

	memberID, err := uuid.Parse(c.Param("member_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid member id"})
	}

	if err := h.clubService.SetLeader(c.Request().Context(), requesterID, clubID, memberID); err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
