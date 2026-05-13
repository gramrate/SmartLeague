package service_provider

import (
	seriesRepo "SmartLeague/internal/adapters/repository/sql/series"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/service/series"
	"context"

	"github.com/google/uuid"
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

func (s *ServiceProvider) SeriesService() seriesService {
	if s.seriesService == nil {
		repo, err := seriesRepo.NewRepo(s.SQLDB())
		if err != nil {
			s.Logger().Panicf("failed to init series repo: %v", err)
		}
		s.seriesService = series.NewService(repo)
	}
	return s.seriesService
}

