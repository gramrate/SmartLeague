package club

import (
	"SmartLeague/internal/domain/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetAll(c echo.Context) error {
	var req dto.GetAllClubsRequest
	if err := h.formDecoder.Decode(&req, c.QueryParams()); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}
	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{Code: http.StatusBadRequest, Message: err.Error()})
	}

	resp, err := h.clubService.GetAll(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

