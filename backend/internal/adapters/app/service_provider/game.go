package service_provider

import (
	seriesRepo "SmartLeague/internal/adapters/repository/sql/series"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/service/game"
	"context"

	"github.com/google/uuid"
)

type gameService interface {
	Create(ctx context.Context, requesterID uuid.UUID, req *dto.CreateGameRequest) (*dto.CreateGameResponse, error)
	ListBySeries(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesGamesRequest) (*dto.GetSeriesGamesResponse, error)
	Get(ctx context.Context, requesterID *uuid.UUID, req *dto.GetGameRequest) (*dto.GetGameResponse, error)
	Update(ctx context.Context, requesterID uuid.UUID, req *dto.UpdateGameRequest) (*dto.UpdateGameResponse, error)
	SaveDraft(ctx context.Context, requesterID uuid.UUID, req *dto.SaveGameDraftRequest) error
	Publish(ctx context.Context, requesterID uuid.UUID, req *dto.PublishGameRequest) error
	SetParticipants(ctx context.Context, requesterID uuid.UUID, req *dto.SetGameParticipantsRequest) error
	UpsertResults(ctx context.Context, requesterID uuid.UUID, req *dto.UpsertGameResultsRequest) error
	GetFull(ctx context.Context, requesterID *uuid.UUID, req *dto.GetGameRequest) (*dto.GetGameFullResponse, error)
	Delete(ctx context.Context, requesterID uuid.UUID, req *dto.DeleteGameRequest) error
}

func (s *ServiceProvider) GameService() gameService {
	if s.gameService == nil {
		repo, err := seriesRepo.NewRepo(s.SQLDB())
		if err != nil {
			s.Logger().Panicf("failed to init series repo for game service: %v", err)
		}
		s.gameService = game.NewService(repo)
	}
	return s.gameService
}
