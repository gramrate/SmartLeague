package server

import (
	_ "SmartLeague/docs"
	"SmartLeague/internal/adapters/app"
	"SmartLeague/internal/adapters/controller/api/middleware/auth"
	"SmartLeague/internal/adapters/controller/api/middleware/role"
	"SmartLeague/internal/adapters/controller/api/v1/club"
	"SmartLeague/internal/adapters/controller/api/v1/ping"
	"SmartLeague/internal/adapters/controller/api/v1/profile"
	"SmartLeague/internal/adapters/controller/api/v1/series"
	"SmartLeague/internal/adapters/controller/api/v1/token"
	"SmartLeague/internal/adapters/controller/api/v1/user"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Setup(app *app.App) {
	app.Server.Logger.SetOutput(io.Discard)
	app.Server.HideBanner = true
	app.Server.Debug = false

	//app.Server.Use(middleware.Recover())

	app.Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://t-leech.vercel.app", "http://localhost:4200", "http://localhost:4000"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	app.Server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		HandleError: true,
		LogError:    true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				app.ServiceProvider.Logger().Infow("request completed",
					"ip", v.RemoteIP,
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
				)
			} else {
				app.ServiceProvider.Logger().Errorw("request failed",
					"ip", v.RemoteIP,
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"error", v.Error.Error(),
				)
			}
			return nil
		},
	}))

	addRouters(app)
}

func addRouters(app *app.App) {
	server := app.Server
	serviceProvider := app.ServiceProvider

	authMiddleware := auth.NewAuthMiddleware(serviceProvider.TokenService(), serviceProvider.CookieService())
	roleMiddleware := role.NewRoleMiddleware(serviceProvider.UserService())

	apiV1 := server.Group("/api/v1")

	apiV1.GET("/swagger/*", echoSwagger.WrapHandler)
	apiV1.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/api/v1/swagger/index.html")
	})

	pingHandler := ping.NewHandler()
	pingHandler.Setup(server.Group(""))

	refreshTokenHandler := token.NewHandler(serviceProvider.TokenService(), serviceProvider.CookieService(), serviceProvider.JWTConfig(), serviceProvider.ServerConfig(), serviceProvider.Validator(), serviceProvider.Decoder())
	refreshTokenHandler.Setup(apiV1)

	userHandler := user.NewHandler(serviceProvider.UserService(), serviceProvider.CookieService(), serviceProvider.JWTConfig(), serviceProvider.ServerConfig(), authMiddleware, roleMiddleware, serviceProvider.Validator(), serviceProvider.Decoder())
	userHandler.Setup(apiV1)

	profileHandler := profile.NewHandler(serviceProvider.ProfileService(), authMiddleware, roleMiddleware, serviceProvider.Validator(), serviceProvider.Decoder())
	profileHandler.Setup(apiV1)

	clubHandler := club.NewHandler(serviceProvider.ClubService(), authMiddleware, roleMiddleware, serviceProvider.Validator(), serviceProvider.Decoder())
	clubHandler.Setup(apiV1)

	seriesHandler := series.NewHandler(serviceProvider.SeriesService(), serviceProvider.GameService(), authMiddleware, roleMiddleware, serviceProvider.Validator(), serviceProvider.Decoder())
	seriesHandler.Setup(apiV1)
}
