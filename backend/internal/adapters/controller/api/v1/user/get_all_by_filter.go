package user

import (
	"SmartLeague/internal/domain/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetAllByFilter Get users by filters
//
// @Summary     Get users by filters
// @Description Retrieves a list of users filtered by role, club_state, club name and nickname query.
// @Tags        user
// @Accept      json
// @Produce     json
// @Param       limit          query     int     false  "Max number of users to return"   minimum(1) maximum(100)  example(10)
// @Param       offset         query     int     false  "Pagination offset"               minimum(0)               example(0)
// @Param       role           query     int     false  "Role enum (0–2)"                 Enums(0,1,2)             example(0)
// @Param       club_state     query     int     false  "Club state (1=member,2=resident,3=leader+president)"     Enums(1,2,3) example(2)
// @Param       club           query     string  false  "Filter by club name"                                     example("smart")
// @Param       q              query     string  false  "Search by nickname"                                       example("mishmish")
// @Success     200  {object}  dto.GetAllByFilterUsersResponse
// @Failure     400  {object}  dto.HTTPStatus "Invalid query parameters"
// @Failure     500  {object}  dto.HTTPStatus "Internal server error"
// @Router      /api/v1/user/all [get]
func (h *handler) GetAllByFilter(c echo.Context) error {
	var req dto.GetAllByFilterUsersRequest

	if err := h.formDecoder.Decode(&req, c.QueryParams()); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	resp, err := h.userService.GetAllByFilter(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HTTPStatus{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}
