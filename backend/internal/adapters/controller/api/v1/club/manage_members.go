package club

import (
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/types"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// SetMemberRole Set club member role
//
// @Summary Set club member role
// @Tags club
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path string true "Club ID"
// @Param member_id path string true "Member Profile ID"
// @Param request body object{state=int} true "New club state (member/resident/leader)"
// @Success 204
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/club/{id}/member/{member_id}/role [post]
func (h *handler) SetMemberRole(c echo.Context) error {
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
	var body struct {
		State int16 `json:"state"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	if err := h.clubService.SetMemberRole(c.Request().Context(), requesterID, clubID, memberID, types.ClubState(body.State)); err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// KickMember Kick member from club
//
// @Summary Kick member from club
// @Tags club
// @Produce json
// @Security CookieAuth
// @Param id path string true "Club ID"
// @Param member_id path string true "Member Profile ID"
// @Success 204
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/club/{id}/member/{member_id}/kick [post]
func (h *handler) KickMember(c echo.Context) error {
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
	if err := h.clubService.KickMember(c.Request().Context(), requesterID, clubID, memberID); err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// BlockMember Block member in club (also removes from club)
//
// @Summary Block member in club
// @Tags club
// @Produce json
// @Security CookieAuth
// @Param id path string true "Club ID"
// @Param member_id path string true "Member Profile ID"
// @Success 204
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/club/{id}/member/{member_id}/block [post]
func (h *handler) BlockMember(c echo.Context) error {
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
	if err := h.clubService.BlockMember(c.Request().Context(), requesterID, clubID, memberID); err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
