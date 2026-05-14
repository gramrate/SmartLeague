package club

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Join Join club
//
// @Summary Join club
// @Tags club
// @Produce json
// @Param id path string true "Club ID"
// @Success 204
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/club/{id}/join [post]
func (h *handler) Join(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}

	clubID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	if err := h.clubService.Join(c.Request().Context(), &dto.JoinClubRequest{ProfileID: userID, ClubID: clubID}); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// Leave Leave current club
//
// @Summary Leave club
// @Tags club
// @Produce json
// @Success 204
// @Failure 401 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/club/leave [post]
func (h *handler) Leave(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}

	if err := h.clubService.Leave(c.Request().Context(), &dto.LeaveClubRequest{ProfileID: userID}); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
