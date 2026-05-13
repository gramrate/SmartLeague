package series

import (
	"SmartLeague/internal/adapters/controller/api/middleware/auth"
	"SmartLeague/internal/adapters/controller/api/middleware/role"
	"SmartLeague/internal/adapters/controller/api/validator"
	"SmartLeague/internal/domain/dto"
	"context"

	"github.com/go-playground/form"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type seriesService interface {
	CreateSeries(ctx context.Context, requesterID uuid.UUID, req *dto.CreateSeriesRequest) (*dto.CreateSeriesResponse, error)
	GetSeries(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesRequest) (*dto.GetSeriesResponse, error)
	GetClubSeries(ctx context.Context, requesterID *uuid.UUID, req *dto.GetClubSeriesRequest) (*dto.GetClubSeriesResponse, error)
	UpdateSeries(ctx context.Context, requesterID uuid.UUID, req *dto.UpdateSeriesRequest) (*dto.UpdateSeriesResponse, error)
	DeleteSeries(ctx context.Context, requesterID uuid.UUID, req *dto.DeleteSeriesRequest) error

	CreateGame(ctx context.Context, requesterID uuid.UUID, req *dto.CreateGameRequest) (*dto.CreateGameResponse, error)
	GetSeriesGames(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesGamesRequest) (*dto.GetSeriesGamesResponse, error)
	UpdateGame(ctx context.Context, requesterID uuid.UUID, req *dto.UpdateGameRequest) (*dto.UpdateGameResponse, error)
	SetGameParticipants(ctx context.Context, requesterID uuid.UUID, req *dto.SetGameParticipantsRequest) error
	UpsertGameResults(ctx context.Context, requesterID uuid.UUID, req *dto.UpsertGameResultsRequest) error
	GetSeriesLeaderboard(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesLeaderboardRequest) (*dto.GetSeriesLeaderboardResponse, error)
}

type handler struct {
	seriesService  seriesService
	authMiddleware *auth.Middleware
	roleMiddleware *role.Middleware
	validator      *validator.Validator
	formDecoder    *form.Decoder
}

func NewHandler(
	seriesService seriesService,
	authMiddleware *auth.Middleware,
	roleMiddleware *role.Middleware,
	validator *validator.Validator,
	formDecoder *form.Decoder,
) *handler {
	return &handler{
		seriesService:  seriesService,
		authMiddleware: authMiddleware,
		roleMiddleware: roleMiddleware,
		validator:      validator,
		formDecoder:    formDecoder,
	}
}

func (h *handler) Setup(router *echo.Group) {
	router.POST("/series", h.CreateSeries, h.authMiddleware.RequireAuth)
	router.GET("/series/:id", h.GetSeries)
	router.PATCH("/series/:id", h.UpdateSeries, h.authMiddleware.RequireAuth)
	router.DELETE("/series/:id", h.DeleteSeries, h.authMiddleware.RequireAuth)

	router.GET("/club/:id/series", h.GetClubSeries)

	router.POST("/series/:id/games", h.CreateGame, h.authMiddleware.RequireAuth)
	router.GET("/series/:id/games", h.GetSeriesGames)
	router.GET("/series/:id/leaderboard", h.GetLeaderboard)

	router.PATCH("/game/:id", h.UpdateGame, h.authMiddleware.RequireAuth)
	router.POST("/game/:id/participants", h.SetParticipants, h.authMiddleware.RequireAuth)
	router.POST("/game/:id/results", h.UpsertResults, h.authMiddleware.RequireAuth)
}

