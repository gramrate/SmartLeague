package series

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GetGame Get game by id
//
// @Summary Get game by id
// @Tags game
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} dto.GetGameResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/game/{id} [get]
func (h *handler) GetGame(c echo.Context) error {
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}
	resp, err := h.gameService.Get(c.Request().Context(), maybeRequesterID(c.Get("user_id")), &dto.GetGameRequest{ID: gameID})
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

// GetGameFull Get full game object (participants + results)
//
// @Summary Get full game object
// @Tags game
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} dto.GetGameFullResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/game/{id}/full [get]
func (h *handler) GetGameFull(c echo.Context) error {
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}
	resp, err := h.gameService.GetFull(c.Request().Context(), maybeRequesterID(c.Get("user_id")), &dto.GetGameRequest{ID: gameID})
	if err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

// DeleteGame Delete game by id
//
// @Summary Delete game by id
// @Tags game
// @Produce json
// @Security CookieAuth
// @Param id path string true "Game ID"
// @Success 204
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/game/{id} [delete]
func (h *handler) DeleteGame(c echo.Context) error {
	requesterID, ok := c.Get("user_id").(uuid.UUID)
	if !ok || requesterID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, dto.HTTPStatus{Code: http.StatusUnauthorized, Message: "unauthorized"})
	}
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}
	if err := h.gameService.Delete(c.Request().Context(), requesterID, &dto.DeleteGameRequest{ID: gameID}); err != nil {
		return c.JSON(http.StatusForbidden, dto.HTTPStatus{Code: http.StatusForbidden, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
