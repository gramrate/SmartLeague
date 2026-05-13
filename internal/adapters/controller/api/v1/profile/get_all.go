package profile

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Get all profiles
//
// @Summary Get all profiles
// @Tags profile
// @Produce json
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {object} dto.GetAllProfilesResponse
// @Failure 400 {object} dto.HTTPStatus
// @Failure 500 {object} dto.HTTPStatus
// @Router /api/v1/profile/all [get]
func (h *handler) GetAll(c echo.Context) error {
	var req dto.GetAllProfilesRequest
	if err := h.formDecoder.Decode(&req, c.QueryParams()); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.profileService.GetAll(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

