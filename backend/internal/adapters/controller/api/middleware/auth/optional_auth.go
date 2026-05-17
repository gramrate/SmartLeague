package auth

import (
	"SmartLeague/internal/domain/common/errorz"
	"errors"

	"github.com/labstack/echo/v4"
)

// OptionalAuth parses access token if present and sets user_id in context.
// It never blocks request flow for public endpoints.
func (m *Middleware) OptionalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := m.cookieService.ReadAccessTokenCookie(c.Request())
		if err != nil {
			if errors.Is(err, errorz.NoCookie) {
				return next(c)
			}
			return next(c)
		}

		userID, err := m.tokenService.ParseAccessToken(c.Request().Context(), token)
		if err == nil {
			c.Set("user_id", userID)
		}

		return next(c)
	}
}

