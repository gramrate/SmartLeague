package series

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"context"
	"math"

	"github.com/google/uuid"
)

type repo interface {
	GetProfileClubState(ctx context.Context, profileID uuid.UUID) (clubID *uuid.UUID, state types.ClubState, err error)

	CreateSeries(ctx context.Context, s model.Series) (*model.Series, error)
	GetSeriesByID(ctx context.Context, id uuid.UUID) (*model.Series, error)
	ListSeriesByClub(ctx context.Context, clubID uuid.UUID, includeClosed bool, limit, offset int) ([]*model.Series, int, error)
	UpdateSeries(ctx context.Context, id uuid.UUID, patch model.SeriesUpdatePatch) (*model.Series, error)
	DeleteSeries(ctx context.Context, id uuid.UUID) error

	CreateGame(ctx context.Context, g model.Game) (*model.Game, error)
	GetGameByID(ctx context.Context, id uuid.UUID) (*model.Game, error)
	ListGamesBySeries(ctx context.Context, seriesID uuid.UUID, limit, offset int) ([]*model.Game, int, error)
	UpdateGame(ctx context.Context, id uuid.UUID, patch model.GameUpdatePatch) (*model.Game, error)
	ReplaceGameParticipants(ctx context.Context, gameID uuid.UUID, participantIDs []uuid.UUID) error
	UpsertGameResults(ctx context.Context, gameID uuid.UUID, rows []model.GameResultRow) error
	ListSeriesLeaderboard(ctx context.Context, seriesID uuid.UUID, limit, offset int) ([]*model.LeaderboardRow, int, error)
}

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{repo: repo}
}

func toSeriesDTO(s *model.Series) *dto.Series {
	return &dto.Series{
		ID:           s.ID,
		ClubID:       s.ClubID,
		Name:         s.Name,
		ScoringRules: s.ScoringRules,
		StartAt:      s.StartAt,
		EndAt:        s.EndAt,
		Description:  s.Description,
		PriceRub:     s.PriceRub,
		IsClosed:     s.IsClosed,
		GameType:     s.GameType,
		Status:       s.Status,
	}
}

func toGameDTO(g *model.Game) *dto.Game {
	return &dto.Game{
		ID:          g.ID,
		SeriesID:    g.SeriesID,
		Name:        g.Name,
		Number:      g.Number,
		Description: g.Description,
		HostID:      g.HostID,
		Status:      g.Status,
	}
}

func canManageClub(state types.ClubState) bool {
	return state == types.ClubStateLeader || state == types.ClubStatePresident
}

func (s *Service) CreateSeries(ctx context.Context, requesterID uuid.UUID, req *dto.CreateSeriesRequest) (*dto.CreateSeriesResponse, error) {
	clubID, clubState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if clubID == nil || !canManageClub(clubState) {
		return nil, errorz.Unauthorized
	}

	created, err := s.repo.CreateSeries(ctx, model.Series{
		ID:           uuid.New(),
		ClubID:       *clubID,
		CreatorID:    requesterID,
		Name:         req.Name,
		ScoringRules: req.ScoringRules,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
		Description:  req.Description,
		PriceRub:     req.PriceRub,
		IsClosed:     req.IsClosed,
		GameType:     req.GameType,
		Status:       req.Status,
	})
	if err != nil {
		return nil, err
	}

	resp := dto.CreateSeriesResponse(*toSeriesDTO(created))
	return &resp, nil
}

func (s *Service) GetSeries(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesRequest) (*dto.GetSeriesResponse, error) {
	ser, err := s.repo.GetSeriesByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if ser.IsClosed {
		if requesterID == nil {
			return nil, errorz.Unauthorized
		}
		clubID, _, err := s.repo.GetProfileClubState(ctx, *requesterID)
		if err != nil {
			return nil, err
		}
		if clubID == nil || *clubID != ser.ClubID {
			return nil, errorz.Unauthorized
		}
	}

	resp := dto.GetSeriesResponse(*toSeriesDTO(ser))
	return &resp, nil
}

func (s *Service) GetClubSeries(ctx context.Context, requesterID *uuid.UUID, req *dto.GetClubSeriesRequest) (*dto.GetClubSeriesResponse, error) {
	limit := 10
	offset := 0
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}

	includeClosed := false
	if requesterID != nil {
		clubID, _, err := s.repo.GetProfileClubState(ctx, *requesterID)
		if err != nil {
			return nil, err
		}
		includeClosed = clubID != nil && *clubID == req.ClubID
	}

	items, total, err := s.repo.ListSeriesByClub(ctx, req.ClubID, includeClosed, limit, offset)
	if err != nil {
		return nil, err
	}

	outItems := make([]*dto.Series, 0, len(items))
	for _, it := range items {
		outItems = append(outItems, toSeriesDTO(it))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	currentPage := (offset / limit) + 1
	if totalPages == 0 {
		totalPages = 1
		currentPage = 1
	}

	return &dto.GetClubSeriesResponse{
		Items: outItems,
		Pagination: dto.PaginationInfo{
			TotalItems:  total,
			TotalPages:  totalPages,
			CurrentPage: currentPage,
			HasNext:     offset+limit < total,
			HasPrevious: offset > 0,
		},
	}, nil
}

func (s *Service) UpdateSeries(ctx context.Context, requesterID uuid.UUID, req *dto.UpdateSeriesRequest) (*dto.UpdateSeriesResponse, error) {
	ser, err := s.repo.GetSeriesByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	clubID, clubState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if clubID == nil || *clubID != ser.ClubID || !canManageClub(clubState) {
		return nil, errorz.Unauthorized
	}

	patch := model.SeriesUpdatePatch{
		Name:         req.Name,
		ScoringRules: req.ScoringRules,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
		Description:  req.Description,
		PriceRub:     req.PriceRub,
		IsClosed:     req.IsClosed,
		GameType:     req.GameType,
		Status:       req.Status,
	}

	updated, err := s.repo.UpdateSeries(ctx, req.ID, patch)
	if err != nil {
		return nil, err
	}
	resp := dto.UpdateSeriesResponse(*toSeriesDTO(updated))
	return &resp, nil
}

func (s *Service) DeleteSeries(ctx context.Context, requesterID uuid.UUID, req *dto.DeleteSeriesRequest) error {
	ser, err := s.repo.GetSeriesByID(ctx, req.ID)
	if err != nil {
		return err
	}
	clubID, clubState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return err
	}
	if clubID == nil || *clubID != ser.ClubID || !canManageClub(clubState) {
		return errorz.Unauthorized
	}
	return s.repo.DeleteSeries(ctx, req.ID)
}

func (s *Service) CreateGame(ctx context.Context, requesterID uuid.UUID, req *dto.CreateGameRequest) (*dto.CreateGameResponse, error) {
	ser, err := s.repo.GetSeriesByID(ctx, req.SeriesID)
	if err != nil {
		return nil, err
	}
	clubID, clubState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if clubID == nil || *clubID != ser.ClubID || !canManageClub(clubState) {
		return nil, errorz.Unauthorized
	}

	name := ""
	if req.Name != nil {
		name = *req.Name
	}

	created, err := s.repo.CreateGame(ctx, model.Game{
		ID:          uuid.New(),
		SeriesID:    req.SeriesID,
		Name:        name,
		Number:      req.Number,
		Description: req.Description,
		HostID:      req.HostID,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}

	resp := dto.CreateGameResponse(*toGameDTO(created))
	return &resp, nil
}

func (s *Service) GetSeriesGames(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesGamesRequest) (*dto.GetSeriesGamesResponse, error) {
	ser, err := s.repo.GetSeriesByID(ctx, req.SeriesID)
	if err != nil {
		return nil, err
	}
	if ser.IsClosed {
		if requesterID == nil {
			return nil, errorz.Unauthorized
		}
		clubID, _, err := s.repo.GetProfileClubState(ctx, *requesterID)
		if err != nil {
			return nil, err
		}
		if clubID == nil || *clubID != ser.ClubID {
			return nil, errorz.Unauthorized
		}
	}

	limit := 10
	offset := 0
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}

	items, total, err := s.repo.ListGamesBySeries(ctx, req.SeriesID, limit, offset)
	if err != nil {
		return nil, err
	}

	outItems := make([]*dto.Game, 0, len(items))
	for _, it := range items {
		outItems = append(outItems, toGameDTO(it))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	currentPage := (offset / limit) + 1
	if totalPages == 0 {
		totalPages = 1
		currentPage = 1
	}

	return &dto.GetSeriesGamesResponse{
		Items: outItems,
		Pagination: dto.PaginationInfo{
			TotalItems:  total,
			TotalPages:  totalPages,
			CurrentPage: currentPage,
			HasNext:     offset+limit < total,
			HasPrevious: offset > 0,
		},
	}, nil
}

func (s *Service) UpdateGame(ctx context.Context, requesterID uuid.UUID, req *dto.UpdateGameRequest) (*dto.UpdateGameResponse, error) {
	game, err := s.repo.GetGameByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	ser, err := s.repo.GetSeriesByID(ctx, game.SeriesID)
	if err != nil {
		return nil, err
	}
	clubID, clubState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if clubID == nil || *clubID != ser.ClubID || !canManageClub(clubState) {
		return nil, errorz.Unauthorized
	}

	updated, err := s.repo.UpdateGame(ctx, req.ID, model.GameUpdatePatch{
		Name:        req.Name,
		Description: req.Description,
		HostID:      req.HostID,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}
	resp := dto.UpdateGameResponse(*toGameDTO(updated))
	return &resp, nil
}

func (s *Service) SetGameParticipants(ctx context.Context, requesterID uuid.UUID, req *dto.SetGameParticipantsRequest) error {
	game, err := s.repo.GetGameByID(ctx, req.GameID)
	if err != nil {
		return err
	}
	ser, err := s.repo.GetSeriesByID(ctx, game.SeriesID)
	if err != nil {
		return err
	}
	clubID, clubState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return err
	}
	if clubID == nil || *clubID != ser.ClubID || !canManageClub(clubState) {
		return errorz.Unauthorized
	}
	if ser.GameType == types.GameTypeSportMafia && len(req.ParticipantIDs) != 10 {
		return errorz.InvalidRequest
	}
	return s.repo.ReplaceGameParticipants(ctx, req.GameID, req.ParticipantIDs)
}

func (s *Service) UpsertGameResults(ctx context.Context, requesterID uuid.UUID, req *dto.UpsertGameResultsRequest) error {
	game, err := s.repo.GetGameByID(ctx, req.GameID)
	if err != nil {
		return err
	}
	ser, err := s.repo.GetSeriesByID(ctx, game.SeriesID)
	if err != nil {
		return err
	}
	clubID, clubState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return err
	}
	if clubID == nil || *clubID != ser.ClubID || !canManageClub(clubState) {
		return errorz.Unauthorized
	}

	rows := make([]model.GameResultRow, 0, len(req.Rows))
	for _, rrow := range req.Rows {
		rows = append(rows, model.GameResultRow{
			GameID:       req.GameID,
			ProfileID:    rrow.ProfileID,
			Place:        rrow.Place,
			Role:         rrow.Role,
			BestMove:     rrow.BestMove,
			FirstKilled:  rrow.FirstKilled,
			Compensation: rrow.Compensation,
			YellowCards:  rrow.YellowCards,
			Removed:      rrow.Removed,
			ExtraPoints:  rrow.ExtraPoints,
			TotalPoints:  rrow.TotalPoints,
		})
	}
	return s.repo.UpsertGameResults(ctx, req.GameID, rows)
}

func (s *Service) GetSeriesLeaderboard(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesLeaderboardRequest) (*dto.GetSeriesLeaderboardResponse, error) {
	ser, err := s.repo.GetSeriesByID(ctx, req.SeriesID)
	if err != nil {
		return nil, err
	}
	if ser.IsClosed {
		if requesterID == nil {
			return nil, errorz.Unauthorized
		}
		clubID, _, err := s.repo.GetProfileClubState(ctx, *requesterID)
		if err != nil {
			return nil, err
		}
		if clubID == nil || *clubID != ser.ClubID {
			return nil, errorz.Unauthorized
		}
	}

	limit := 10
	offset := 0
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}

	items, total, err := s.repo.ListSeriesLeaderboard(ctx, req.SeriesID, limit, offset)
	if err != nil {
		return nil, err
	}

	outItems := make([]*dto.LeaderboardRow, 0, len(items))
	for _, it := range items {
		outItems = append(outItems, &dto.LeaderboardRow{ProfileID: it.ProfileID, Points: it.Points})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	currentPage := (offset / limit) + 1
	if totalPages == 0 {
		totalPages = 1
		currentPage = 1
	}

	return &dto.GetSeriesLeaderboardResponse{
		Items: outItems,
		Pagination: dto.PaginationInfo{
			TotalItems:  total,
			TotalPages:  totalPages,
			CurrentPage: currentPage,
			HasNext:     offset+limit < total,
			HasPrevious: offset > 0,
		},
	}, nil
}
