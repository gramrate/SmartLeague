package game

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type repo interface {
	GetProfileClubState(ctx context.Context, profileID uuid.UUID) (clubID *uuid.UUID, state types.ClubState, err error)

	GetSeriesByID(ctx context.Context, id uuid.UUID) (*model.Series, error)
	IsSeriesParticipant(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID) (bool, error)

	CreateGame(ctx context.Context, g model.Game) (*model.Game, error)
	GetGameByID(ctx context.Context, id uuid.UUID) (*model.Game, error)
	ListGamesBySeries(ctx context.Context, seriesID uuid.UUID, limit, offset int) ([]*model.Game, int, error)
	UpdateGame(ctx context.Context, id uuid.UUID, patch model.GameUpdatePatch) (*model.Game, error)
	ReplaceGameParticipants(ctx context.Context, gameID uuid.UUID, participantIDs []uuid.UUID) error
	UpsertGameResults(ctx context.Context, gameID uuid.UUID, rows []model.GameResultRow) error
	ListGameParticipants(ctx context.Context, gameID uuid.UUID) ([]uuid.UUID, error)
	ListGameResults(ctx context.Context, gameID uuid.UUID) ([]model.GameResultRow, error)
	DeleteGame(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo repo
}

const sportMafiaParticipantsCount = 10

func NewService(repo repo) *Service {
	return &Service{repo: repo}
}

func canManageClub(state types.ClubState) bool {
	return state == types.ClubStateLeader || state == types.ClubStatePresident
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

func (s *Service) Create(ctx context.Context, requesterID uuid.UUID, req *dto.CreateGameRequest) (*dto.CreateGameResponse, error) {
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
		Number:      0,
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

func (s *Service) Get(ctx context.Context, requesterID *uuid.UUID, req *dto.GetGameRequest) (*dto.GetGameResponse, error) {
	game, err := s.repo.GetGameByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	ser, err := s.repo.GetSeriesByID(ctx, game.SeriesID)
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
	resp := dto.GetGameResponse(*toGameDTO(game))
	return &resp, nil
}

func (s *Service) ListBySeries(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesGamesRequest) (*dto.GetSeriesGamesResponse, error) {
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

func (s *Service) Update(ctx context.Context, requesterID uuid.UUID, req *dto.UpdateGameRequest) (*dto.UpdateGameResponse, error) {
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

func (s *Service) SetParticipants(ctx context.Context, requesterID uuid.UUID, req *dto.SetGameParticipantsRequest) error {
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

	if len(req.ParticipantIDs) != sportMafiaParticipantsCount {
		return errorz.InvalidRequest
	}

	// ensure participants are in the series
	for _, pid := range req.ParticipantIDs {
		ok, err := s.repo.IsSeriesParticipant(ctx, ser.ID, pid)
		if err != nil {
			return err
		}
		if !ok {
			return errorz.InvalidRequest
		}
	}

	return s.repo.ReplaceGameParticipants(ctx, req.GameID, req.ParticipantIDs)
}

func (s *Service) UpsertResults(ctx context.Context, requesterID uuid.UUID, req *dto.UpsertGameResultsRequest) error {
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
	if len(req.Rows) != sportMafiaParticipantsCount {
		return errorz.InvalidRequest
	}

	rows := make([]model.GameResultRow, 0, len(req.Rows))
	for _, rrow := range req.Rows {
		if rrow.BestMove != nil && !isValidBestMove(*rrow.BestMove) {
			return errorz.InvalidRequest
		}
		if rrow.Role != nil && !isValidMafiaRole(*rrow.Role) {
			return errorz.InvalidRequest
		}
		rows = append(rows, model.GameResultRow{
			GameID:        req.GameID,
			ProfileID:     rrow.ProfileID,
			Place:         rrow.Place,
			Role:          rrow.Role,
			BestMove:      rrow.BestMove,
			FirstKilled:   rrow.FirstKilled,
			Compensation:  rrow.Compensation,
			YellowCards:   rrow.YellowCards,
			Removed:       rrow.Removed,
			VictoryPoints: rrow.VictoryPoints,
			ExtraPoints:   rrow.ExtraPoints,
			TotalPoints:   rrow.TotalPoints,
		})
	}
	return s.repo.UpsertGameResults(ctx, req.GameID, rows)
}

func isValidBestMove(raw string) bool {
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';' || r == ' '
	})
	if len(parts) == 0 || len(parts) > 3 {
		return false
	}
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return false
		}
	}
	return true
}

func isValidMafiaRole(role types.MafiaRole) bool {
	switch role {
	case types.MafiaRoleCivilian, types.MafiaRoleMafia, types.MafiaRoleDon, types.MafiaRoleSheriff:
		return true
	default:
		return false
	}
}

func (s *Service) GetFull(ctx context.Context, requesterID *uuid.UUID, req *dto.GetGameRequest) (*dto.GetGameFullResponse, error) {
	game, err := s.repo.GetGameByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	ser, err := s.repo.GetSeriesByID(ctx, game.SeriesID)
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

	participantIDs, err := s.repo.ListGameParticipants(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	results, err := s.repo.ListGameResults(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	dtoResults := make([]dto.GameResultRow, 0, len(results))
	for _, rr := range results {
		dtoResults = append(dtoResults, dto.GameResultRow{
			ProfileID:     rr.ProfileID,
			Place:         rr.Place,
			Role:          rr.Role,
			BestMove:      rr.BestMove,
			FirstKilled:   rr.FirstKilled,
			Compensation:  rr.Compensation,
			YellowCards:   rr.YellowCards,
			Removed:       rr.Removed,
			VictoryPoints: rr.VictoryPoints,
			ExtraPoints:   rr.ExtraPoints,
			TotalPoints:   rr.TotalPoints,
		})
	}

	resp := dto.GetGameFullResponse(dto.GameFull{
		Game:           *toGameDTO(game),
		ParticipantIDs: participantIDs,
		Results:        dtoResults,
	})
	return &resp, nil
}

func (s *Service) Delete(ctx context.Context, requesterID uuid.UUID, req *dto.DeleteGameRequest) error {
	game, err := s.repo.GetGameByID(ctx, req.ID)
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
	return s.repo.DeleteGame(ctx, req.ID)
}
