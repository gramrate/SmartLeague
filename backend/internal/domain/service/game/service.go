package game

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"context"
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type repo interface {
	GetProfileClubState(ctx context.Context, profileID uuid.UUID) (clubID *uuid.UUID, state types.ClubState, err error)

	GetSeriesByID(ctx context.Context, id uuid.UUID) (*model.Series, error)
	IsSeriesParticipant(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID) (bool, error)
	IsSeriesJudge(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID) (bool, error)
	FindOrCreateGuest(ctx context.Context, seriesID uuid.UUID, nickname string) (*model.Guest, error)

	CreateGame(ctx context.Context, g model.Game) (*model.Game, error)
	GetGameByID(ctx context.Context, id uuid.UUID) (*model.Game, error)
	ListGamesBySeries(ctx context.Context, seriesID uuid.UUID, limit, offset int, includeDrafts bool) ([]*model.Game, int, error)
	UpdateGame(ctx context.Context, id uuid.UUID, patch model.GameUpdatePatch) (*model.Game, error)
	ReplaceGameParticipants(ctx context.Context, gameID uuid.UUID, participantIDs []uuid.UUID) error
	UpsertGameResults(ctx context.Context, gameID uuid.UUID, rows []model.GameResultRow) error
	ClearGameResults(ctx context.Context, gameID uuid.UUID) error
	ListGameParticipants(ctx context.Context, gameID uuid.UUID) ([]uuid.UUID, error)
	ListGameResults(ctx context.Context, gameID uuid.UUID) ([]model.GameResultRow, error)
	DeleteGame(ctx context.Context, id uuid.UUID) error
	SetGameJudge(ctx context.Context, gameID uuid.UUID, judgeID *uuid.UUID, confirmed bool) error
	UpsertGameDraft(ctx context.Context, d *model.GameDraft) error
	GetGameDraft(ctx context.Context, gameID uuid.UUID) (*model.GameDraft, error)
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

func (s *Service) canManageSeries(ctx context.Context, requesterID *uuid.UUID, clubID uuid.UUID) (bool, error) {
	if requesterID == nil {
		return false, nil
	}
	profileClubID, profileClubState, err := s.repo.GetProfileClubState(ctx, *requesterID)
	if err != nil {
		return false, err
	}
	return profileClubID != nil && *profileClubID == clubID && canManageClub(profileClubState), nil
}

// canEditGame checks if requester is a club manager OR a series judge.
// Judges can draft/publish games but cannot create or delete them.
func (s *Service) canEditGame(ctx context.Context, requesterID *uuid.UUID, ser *model.Series) (bool, error) {
	if requesterID == nil {
		return false, nil
	}
	ok, err := s.canManageSeries(ctx, requesterID, ser.ClubID)
	if err != nil || ok {
		return ok, err
	}
	return s.repo.IsSeriesJudge(ctx, ser.ID, *requesterID)
}

func (s *Service) canAccessSeries(ctx context.Context, requesterID *uuid.UUID, series *model.Series) (bool, error) {
	if !series.IsClubOnly {
		return true, nil
	}
	if series.ShowToAll {
		return true, nil
	}
	if requesterID == nil {
		return false, nil
	}
	isParticipant, err := s.repo.IsSeriesParticipant(ctx, series.ID, *requesterID)
	if err != nil {
		return false, err
	}
	if isParticipant {
		return true, nil
	}
	clubID, _, err := s.repo.GetProfileClubState(ctx, *requesterID)
	if err != nil {
		return false, err
	}
	return clubID != nil && *clubID == series.ClubID, nil
}

func toGameDTO(g *model.Game) *dto.Game {
	d := &dto.Game{
		ID:                 g.ID,
		SeriesID:           g.SeriesID,
		Name:               g.Name,
		Number:             g.Number,
		Description:        g.Description,
		HostID:             g.HostID,
		Status:             g.Status,
		IsTournament:       g.IsTournament,
		GameJudgeID:        g.JudgeID,
		GameJudgeConfirmed: g.JudgeConfirmed,
	}
	if g.JudgeID != nil {
		d.GameJudgeNickname = g.JudgeNickname
	}
	return d
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
	canAccess, err := s.canAccessSeries(ctx, requesterID, ser)
	if err != nil {
		return nil, err
	}
	if !canAccess {
		return nil, errorz.Unauthorized
	}
	if game.Status == types.GameStatusDraft {
		canEdit, err := s.canEditGame(ctx, requesterID, ser)
		if err != nil {
			return nil, err
		}
		if !canEdit {
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
	canAccess, err := s.canAccessSeries(ctx, requesterID, ser)
	if err != nil {
		return nil, err
	}
	if !canAccess {
		return nil, errorz.Unauthorized
	}

	limit := 10
	offset := 0
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}

	includeDrafts, err := s.canEditGame(ctx, requesterID, ser)
	if err != nil {
		return nil, err
	}

	items, total, err := s.repo.ListGamesBySeries(ctx, req.SeriesID, limit, offset, includeDrafts)
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

// resolveSlotPlayer resolves a ManageGameRow to either a profile UUID or a guest UUID.
// Returns (profileID, guestID, error).
func (s *Service) resolveSlotPlayer(ctx context.Context, seriesID uuid.UUID, r dto.ManageGameRow) (*uuid.UUID, *uuid.UUID, error) {
	if r.ProfileID != nil {
		return r.ProfileID, nil, nil
	}
	if r.GuestNickname != nil && strings.TrimSpace(*r.GuestNickname) != "" {
		guest, err := s.repo.FindOrCreateGuest(ctx, seriesID, *r.GuestNickname)
		if err != nil {
			return nil, nil, err
		}
		return nil, &guest.ID, nil
	}
	return nil, nil, nil
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
	canAccess, err := s.canAccessSeries(ctx, requesterID, ser)
	if err != nil {
		return nil, err
	}
	if !canAccess {
		return nil, errorz.Unauthorized
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
			GuestID:       rr.GuestID,
			GuestNickname: rr.GuestNickname,
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

	canEdit, err := s.canEditGame(ctx, requesterID, ser)
	if err != nil {
		return nil, err
	}

	var draftData *dto.GameDraftData
	if canEdit {
		draft, err := s.repo.GetGameDraft(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		if draft != nil {
			var draftRows []dto.ManageGameRow
			_ = json.Unmarshal(draft.Rows, &draftRows)
			draftData = &dto.GameDraftData{
				Rows:           draftRows,
				JudgeID:        draft.JudgeID,
				JudgeConfirmed: draft.JudgeConfirmed,
			}
		}
	}

	resp := dto.GetGameFullResponse(dto.GameFull{
		Game:           *toGameDTO(game),
		ParticipantIDs: participantIDs,
		Results:        dtoResults,
		Draft:          draftData,
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

func (s *Service) SaveDraft(ctx context.Context, requesterID uuid.UUID, req *dto.SaveGameDraftRequest) error {
	game, err := s.repo.GetGameByID(ctx, req.GameID)
	if err != nil {
		return err
	}
	ser, err := s.repo.GetSeriesByID(ctx, game.SeriesID)
	if err != nil {
		return err
	}
	canEdit, err := s.canEditGame(ctx, &requesterID, ser)
	if err != nil {
		return err
	}
	if !canEdit {
		return errorz.Unauthorized
	}

	return s.repo.UpsertGameDraft(ctx, &model.GameDraft{
		GameID:         req.GameID,
		Rows:           req.RawRows,
		JudgeID:        req.JudgeID,
		JudgeConfirmed: req.JudgeConfirmed,
	})
}

func (s *Service) Publish(ctx context.Context, requesterID uuid.UUID, req *dto.PublishGameRequest) error {
	if len(req.Rows) != sportMafiaParticipantsCount {
		return errorz.InvalidRequest
	}

	game, err := s.repo.GetGameByID(ctx, req.GameID)
	if err != nil {
		return err
	}
	ser, err := s.repo.GetSeriesByID(ctx, game.SeriesID)
	if err != nil {
		return err
	}
	canEdit, err := s.canEditGame(ctx, &requesterID, ser)
	if err != nil {
		return err
	}
	if !canEdit {
		return errorz.Unauthorized
	}

	seenSlots := map[int]bool{}
	bestMoveCount := 0
	participantIDs := make([]uuid.UUID, 0, sportMafiaParticipantsCount)
	resultRows := make([]model.GameResultRow, 0, sportMafiaParticipantsCount)

	for _, r := range req.Rows {
		if r.Slot < 1 || r.Slot > 10 || seenSlots[r.Slot] || r.Role == nil {
			return errorz.InvalidRequest
		}
		if !isValidMafiaRole(*r.Role) {
			return errorz.InvalidRequest
		}
		seenSlots[r.Slot] = true

		profileID, guestID, err := s.resolveSlotPlayer(ctx, ser.ID, r)
		if err != nil {
			return err
		}
		// Every slot must have a player (profile or guest)
		if profileID == nil && guestID == nil {
			return errorz.InvalidRequest
		}
		// Registered profile must be a series participant
		if profileID != nil {
			ok, err := s.repo.IsSeriesParticipant(ctx, ser.ID, *profileID)
			if err != nil {
				return err
			}
			if !ok {
				return errorz.InvalidRequest
			}
			participantIDs = append(participantIDs, *profileID)
		}

		if r.BestMove != nil {
			if !isValidBestMove(*r.BestMove) {
				return errorz.InvalidRequest
			}
			bestMoveCount++
		}

		place := r.Slot
		resultRows = append(resultRows, model.GameResultRow{
			GameID:        req.GameID,
			ProfileID:     profileID,
			GuestID:       guestID,
			Place:         &place,
			Role:          r.Role,
			BestMove:      r.BestMove,
			FirstKilled:   false,
			Compensation:  r.Compensation,
			YellowCards:   r.YellowCards,
			Removed:       r.Removed,
			VictoryPoints: 0,
			ExtraPoints:   r.ExtraPoints,
			TotalPoints:   r.TotalPoints,
		})
	}
	if len(seenSlots) != sportMafiaParticipantsCount || bestMoveCount > 1 {
		return errorz.InvalidRequest
	}

	if err := s.repo.ReplaceGameParticipants(ctx, req.GameID, participantIDs); err != nil {
		return err
	}
	if err := s.repo.ClearGameResults(ctx, req.GameID); err != nil {
		return err
	}
	if err := s.repo.UpsertGameResults(ctx, req.GameID, resultRows); err != nil {
		return err
	}
	// Tournament series require judge selection confirmed before publishing
	if ser.IsTournament && !req.JudgeConfirmed {
		return errorz.InvalidRequest
	}

	// Save draft to reflect published state
	if req.RawRows != nil {
		_ = s.repo.UpsertGameDraft(ctx, &model.GameDraft{
			GameID:         req.GameID,
			Rows:           req.RawRows,
			JudgeID:        req.JudgeID,
			JudgeConfirmed: req.JudgeConfirmed,
		})
	}

	finished := types.GameStatusFinished
	if _, err = s.repo.UpdateGame(ctx, req.GameID, model.GameUpdatePatch{Status: &finished}); err != nil {
		return err
	}
	return s.repo.SetGameJudge(ctx, req.GameID, req.JudgeID, req.JudgeConfirmed)
}
