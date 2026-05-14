package profile

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/dto"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Create Create profile
//
// @Summary Create profile
// @Tags profile
// @Accept json
// @Produce json
// @Param request body dto.CreateProfileRequest true "Profile data"
// @Success 201 {object} dto.CreateProfileResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 409 {object} dto.HTTPStatus
// @Router /api/v1/profile [post]
func (h *handler) Create(c echo.Context) error {
	var req dto.CreateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.profileService.Create(c.Request().Context(), &req)
	switch {
	case errors.Is(err, errorz.EmailAlreadyExist):
		return c.JSON(http.StatusConflict, dto.HTTPStatus{Code: http.StatusConflict, Message: err.Error()})
	case err != nil:
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}
