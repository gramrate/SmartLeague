package profile

import (
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/types"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// UpdateEach Update profile by id (admin)
//
// @Summary Update profile by id
// @Tags profile
// @Accept json
// @Produce json
// @Param id path string true "Profile ID"
// @Param request body dto.UpdateEachProfileRequest true "Update data"
// @Success 200 {object} dto.UpdateEachProfileResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/profile/{id} [patch]
func (h *handler) UpdateEach(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: "invalid id"})
	}

	var req dto.UpdateEachProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	req.ID = id
	req.RequesterRole = types.RoleAdmin

	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.profileService.UpdateEach(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}
