package club

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Delete Delete club by id
//
// @Summary Delete club by id
// @Tags club
// @Produce json
// @Param id path string true "Club ID"
// @Success 204
// @Failure 400 {object} dto.HTTPStatus
// @Failure 401 {object} dto.HTTPStatus
// @Failure 403 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/club/{id} [delete]
func (h *handler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	if err := h.clubService.Delete(c.Request().Context(), &dto.DeleteClubRequest{ID: id}); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
