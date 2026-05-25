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
	GetAllSeries(ctx context.Context, requesterID *uuid.UUID, req *dto.GetAllSeriesRequest) (*dto.GetAllSeriesResponse, error)
	UpdateSeries(ctx context.Context, requesterID uuid.UUID, req *dto.UpdateSeriesRequest) (*dto.UpdateSeriesResponse, error)
	DeleteSeries(ctx context.Context, requesterID uuid.UUID, req *dto.DeleteSeriesRequest) error
	GetParticipants(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesParticipantsRequest) (*dto.GetSeriesParticipantsResponse, error)
	Join(ctx context.Context, req *dto.JoinSeriesRequest) error
	Leave(ctx context.Context, req *dto.LeaveSeriesRequest) error
	GetLeaderboard(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesLeaderboardRequest) (*dto.GetSeriesLeaderboardResponse, error)
}

type gameService interface {
	Create(ctx context.Context, requesterID uuid.UUID, req *dto.CreateGameRequest) (*dto.CreateGameResponse, error)
	ListBySeries(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesGamesRequest) (*dto.GetSeriesGamesResponse, error)
	Get(ctx context.Context, requesterID *uuid.UUID, req *dto.GetGameRequest) (*dto.GetGameResponse, error)
	GetFull(ctx context.Context, requesterID *uuid.UUID, req *dto.GetGameRequest) (*dto.GetGameFullResponse, error)
	Update(ctx context.Context, requesterID uuid.UUID, req *dto.UpdateGameRequest) (*dto.UpdateGameResponse, error)
	Delete(ctx context.Context, requesterID uuid.UUID, req *dto.DeleteGameRequest) error
	SaveDraft(ctx context.Context, requesterID uuid.UUID, req *dto.SaveGameDraftRequest) error
	Publish(ctx context.Context, requesterID uuid.UUID, req *dto.PublishGameRequest) error
	SetParticipants(ctx context.Context, requesterID uuid.UUID, req *dto.SetGameParticipantsRequest) error
	UpsertResults(ctx context.Context, requesterID uuid.UUID, req *dto.UpsertGameResultsRequest) error
}

type handler struct {
	seriesService  seriesService
	gameService    gameService
	authMiddleware *auth.Middleware
	roleMiddleware *role.Middleware
	validator      *validator.Validator
	formDecoder    *form.Decoder
}

func NewHandler(
	seriesService seriesService,
	gameService gameService,
	authMiddleware *auth.Middleware,
	roleMiddleware *role.Middleware,
	validator *validator.Validator,
	formDecoder *form.Decoder,
) *handler {
	return &handler{
		seriesService:  seriesService,
		gameService:    gameService,
		authMiddleware: authMiddleware,
		roleMiddleware: roleMiddleware,
		validator:      validator,
		formDecoder:    formDecoder,
	}
}

func (h *handler) Setup(router *echo.Group) {
	router.POST("/series", h.CreateSeries, h.authMiddleware.RequireAuth)
	router.GET("/series/:id", h.GetSeries, h.authMiddleware.OptionalAuth)
	router.GET("/series/:id/full", h.GetSeriesFull, h.authMiddleware.OptionalAuth)
	router.PATCH("/series/:id", h.UpdateSeries, h.authMiddleware.RequireAuth)
	router.DELETE("/series/:id", h.DeleteSeries, h.authMiddleware.RequireAuth)

	router.GET("/club/:id/series", h.GetClubSeries, h.authMiddleware.OptionalAuth)
	router.GET("/series/all", h.GetAllSeries, h.authMiddleware.OptionalAuth)

	router.GET("/series/:id/participants", h.GetParticipants, h.authMiddleware.OptionalAuth)
	router.POST("/series/:id/join", h.JoinSeries, h.authMiddleware.RequireAuth)
	router.POST("/series/:id/leave", h.LeaveSeries, h.authMiddleware.RequireAuth)

	router.POST("/series/:id/games", h.CreateGame, h.authMiddleware.RequireAuth)
	router.POST("/series/:id/games/draft", h.CreateGameDraft, h.authMiddleware.RequireAuth)
	router.GET("/series/:id/games", h.GetSeriesGames, h.authMiddleware.OptionalAuth)
	router.GET("/series/:id/leaderboard", h.GetLeaderboard, h.authMiddleware.OptionalAuth)

	router.GET("/game/:id", h.GetGame, h.authMiddleware.OptionalAuth)
	router.GET("/game/:id/full", h.GetGameFull, h.authMiddleware.OptionalAuth)
	router.DELETE("/game/:id", h.DeleteGame, h.authMiddleware.RequireAuth)
	router.PATCH("/game/:id", h.UpdateGame, h.authMiddleware.RequireAuth)
	router.POST("/game/:id/draft", h.SaveGameDraft, h.authMiddleware.RequireAuth)
	router.POST("/game/:id/publish", h.PublishGame, h.authMiddleware.RequireAuth)
	router.POST("/game/:id/participants", h.SetParticipants, h.authMiddleware.RequireAuth)
	router.POST("/game/:id/results", h.UpsertResults, h.authMiddleware.RequireAuth)
}
