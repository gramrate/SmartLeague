package ping

import (
	"SmartLeague/internal/domain/dto"
	"github.com/labstack/echo/v4"
)

// Ping Healthcheck
//
// @Summary      Checking the server performance
// @Description  Checking the server performance.
// @Tags         ping
// @Accept       json
// @Produce      json
// @Success      200      {object}  dto.PingResponse "Successful check"
// @Router       /ping [get]
func (h *handler) Ping(c echo.Context) error {
	return c.JSON(200, dto.PingResponse{
		Message: "ok",
	})
}
