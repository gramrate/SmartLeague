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

const maxParticipantsSportMafia = 10

type repo interface {
	GetProfileClubState(ctx context.Context, profileID uuid.UUID) (clubID *uuid.UUID, state types.ClubState, err error)
	IsProfileBannedInClub(ctx context.Context, profileID uuid.UUID, clubID uuid.UUID) (bool, error)

	CreateSeries(ctx context.Context, s model.Series) (*model.Series, error)
	GetSeriesByID(ctx context.Context, id uuid.UUID) (*model.Series, error)
	ListSeriesByClub(ctx context.Context, clubID uuid.UUID, includeClosed, includeClubOnly bool, limit, offset int) ([]*model.Series, int, error)
	ListAllSeries(ctx context.Context, limit, offset int, query, clubQuery, from, to *string, isRating *bool, requesterClubID *uuid.UUID, showPast, showClosed bool) ([]*model.SeriesListItem, int, error)
	UpdateSeries(ctx context.Context, id uuid.UUID, patch model.SeriesUpdatePatch) (*model.Series, error)
	DeleteSeries(ctx context.Context, id uuid.UUID) error

	AddSeriesParticipant(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID) error
	RemoveSeriesParticipant(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID) error
	CountSeriesParticipants(ctx context.Context, seriesID uuid.UUID) (int, error)
	ListSeriesParticipants(ctx context.Context, seriesID uuid.UUID, limit, offset int, query *string) ([]*model.User, int, error)
	IsSeriesParticipant(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID) (bool, error)
	ListPaidSeriesParticipants(ctx context.Context, seriesID uuid.UUID) ([]uuid.UUID, error)
	SetSeriesParticipantPaid(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID, paid bool) error

	ListSeriesLeaderboard(ctx context.Context, seriesID uuid.UUID, limit, offset int) ([]*model.LeaderboardRow, int, error)

	ListSeriesJudges(ctx context.Context, seriesID uuid.UUID) ([]*model.SeriesJudge, error)
	SetSeriesJudge(ctx context.Context, seriesID, profileID uuid.UUID, role model.JudgeRole) error
	RemoveSeriesJudge(ctx context.Context, seriesID, profileID uuid.UUID) error
	IsSeriesJudge(ctx context.Context, seriesID, profileID uuid.UUID) (bool, error)
	ClearAllSeriesJudges(ctx context.Context, seriesID uuid.UUID) error
	ClearAllGameJudgesForSeries(ctx context.Context, seriesID uuid.UUID) error
	ClearAllGameDraftJudgesForSeries(ctx context.Context, seriesID uuid.UUID) error
	ClearAllSeriesPayments(ctx context.Context, seriesID uuid.UUID) error
}

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{repo: repo}
}

func canManageClub(state types.ClubState) bool {
	return state == types.ClubStateLeader || state == types.ClubStatePresident
}

func maxParticipantsForGameType(gameType types.GameType) int {
	return maxParticipantsSportMafia
}

func (s *Service) canAccessSeries(ctx context.Context, requesterID *uuid.UUID, series *model.Series) (bool, error) {
	if !series.IsClubOnly {
		return true, nil
	}
	// show_to_all=true means everyone can view (just can't join unless in club)
	if series.ShowToAll {
		return true, nil
	}
	// show_to_all=false: only participants, judges, and club members can view
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

func seriesToDTO(s *model.Series, creatorID *uuid.UUID) *dto.Series {
	return &dto.Series{
		ID:           s.ID,
		ClubID:       s.ClubID,
		CreatorID:    creatorID,
		Name:         s.Name,
		Description:  s.Description,
		StartAt:      s.StartAt,
		EndAt:        s.EndAt,
		PriceRub:     s.PriceRub,
		IsRating:     s.IsRating,
		IsClubOnly:   s.IsClubOnly,
		ShowToAll:    s.ShowToAll,
		IsClosed:     s.IsClosed,
		IsTournament: s.IsTournament,
		GameType:     s.GameType,
	}
}

func judgeToDTO(j *model.SeriesJudge) *dto.SeriesJudge {
	out := &dto.SeriesJudge{
		ProfileID: j.ProfileID,
		Role:      int16(j.Role),
	}
	if j.User != nil {
		out.Nickname = j.User.Nickname
		out.Name = j.User.Name
		out.ShowName = j.User.ShowName
	}
	return out
}

func profileToDTO(p *model.User) *dto.User {
	return &dto.User{
		ID:          p.ID,
		Nickname:    p.Nickname,
		Name:        p.Name,
		ShowName:    p.ShowName,
		Description: p.Description,
		Email:       p.Email,
		ClubID:      p.ClubID,
		ClubState:   p.ClubState,
		Role:        p.Role,
	}
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
		Description:  req.Description,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
		PriceRub:     req.PriceRub,
		IsRating:     req.IsRating != nil && *req.IsRating,
		IsClubOnly:   req.IsClubOnly != nil && *req.IsClubOnly,
		ShowToAll:    req.ShowToAll == nil || *req.ShowToAll,
		IsClosed:     req.IsClosed,
		IsTournament: req.IsTournament,
		GameType:     types.GameTypeSportMafia,
	})
	if err != nil {
		return nil, err
	}

	resp := dto.CreateSeriesResponse(*seriesToDTO(created, &requesterID))
	return &resp, nil
}

func (s *Service) GetSeries(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesRequest) (*dto.GetSeriesResponse, error) {
	ser, err := s.repo.GetSeriesByID(ctx, req.ID)
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

	var creatorID *uuid.UUID
	if requesterID != nil {
		clubID, clubState, err := s.repo.GetProfileClubState(ctx, *requesterID)
		if err != nil {
			return nil, err
		}
		if clubID != nil && *clubID == ser.ClubID && canManageClub(clubState) {
			creatorID = &ser.CreatorID
		}
	}

	resp := dto.GetSeriesResponse(*seriesToDTO(ser, creatorID))
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
	isLeader := false
	if requesterID != nil {
		clubID, clubState, err := s.repo.GetProfileClubState(ctx, *requesterID)
		if err != nil {
			return nil, err
		}
		includeClosed = clubID != nil && *clubID == req.ClubID
		isLeader = includeClosed && canManageClub(clubState)
	}

	items, total, err := s.repo.ListSeriesByClub(ctx, req.ClubID, includeClosed, includeClosed, limit, offset)
	if err != nil {
		return nil, err
	}

	outItems := make([]*dto.Series, 0, len(items))
	for _, it := range items {
		var creatorID *uuid.UUID
		if isLeader {
			creatorID = &it.CreatorID
		}
		outItems = append(outItems, seriesToDTO(it, creatorID))
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

func (s *Service) GetAllSeries(ctx context.Context, requesterID *uuid.UUID, req *dto.GetAllSeriesRequest) (*dto.GetAllSeriesResponse, error) {
	limit := 10
	offset := 0
	showPast := false
	showClosed := false
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}
	if req.ShowPast != nil {
		showPast = *req.ShowPast
	}
	if req.ShowClosed != nil {
		showClosed = *req.ShowClosed
	}

	var requesterClubID *uuid.UUID
	if requesterID != nil {
		clubID, _, err := s.repo.GetProfileClubState(ctx, *requesterID)
		if err != nil {
			return nil, err
		}
		requesterClubID = clubID
	}

	items, total, err := s.repo.ListAllSeries(ctx, limit, offset, req.Query, req.ClubQuery, req.From, req.To, req.IsRating, requesterClubID, showPast, showClosed)
	if err != nil {
		return nil, err
	}

	outItems := make([]*dto.AllSeriesItem, 0, len(items))
	for _, it := range items {
		outItems = append(outItems, &dto.AllSeriesItem{
			ID:           it.ID,
			ClubID:       it.ClubID,
			ClubName:     it.ClubName,
			Name:         it.Name,
			Description:  it.Description,
			StartAt:      it.StartAt,
			EndAt:        it.EndAt,
			PriceRub:     it.PriceRub,
			IsRating:     it.IsRating,
			IsClubOnly:   it.IsClubOnly,
			ShowToAll:    it.ShowToAll,
			IsClosed:     it.IsClosed,
			IsTournament: it.IsTournament,
			GamesCount:   it.GamesCount,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	currentPage := (offset / limit) + 1
	if totalPages == 0 {
		totalPages = 1
		currentPage = 1
	}

	return &dto.GetAllSeriesResponse{
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
		Description:  req.Description,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
		PriceRub:     req.PriceRub,
		IsRating:     req.IsRating,
		IsClubOnly:   req.IsClubOnly,
		ShowToAll:    req.ShowToAll,
		IsClosed:     req.IsClosed,
		IsTournament: req.IsTournament,
	}

	wasTournament := ser.IsTournament
	wasPaid := ser.PriceRub > 0

	updated, err := s.repo.UpdateSeries(ctx, req.ID, patch)
	if err != nil {
		return nil, err
	}

	if wasTournament && !updated.IsTournament {
		if err := s.repo.ClearAllSeriesJudges(ctx, req.ID); err != nil {
			return nil, err
		}
		if err := s.repo.ClearAllGameJudgesForSeries(ctx, req.ID); err != nil {
			return nil, err
		}
		if err := s.repo.ClearAllGameDraftJudgesForSeries(ctx, req.ID); err != nil {
			return nil, err
		}
	}

	if wasPaid && updated.PriceRub <= 0 {
		if err := s.repo.ClearAllSeriesPayments(ctx, req.ID); err != nil {
			return nil, err
		}
	}

	resp := dto.UpdateSeriesResponse(*seriesToDTO(updated, &updated.CreatorID))
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

func (s *Service) GetParticipants(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesParticipantsRequest) (*dto.GetSeriesParticipantsResponse, error) {
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

	items, total, err := s.repo.ListSeriesParticipants(ctx, req.SeriesID, limit, offset, req.Query)
	if err != nil {
		return nil, err
	}

	outItems := make([]*dto.User, 0, len(items))
	for _, it := range items {
		outItems = append(outItems, profileToDTO(it))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	currentPage := (offset / limit) + 1
	if totalPages == 0 {
		totalPages = 1
		currentPage = 1
	}

	return &dto.GetSeriesParticipantsResponse{
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

func (s *Service) GetPayments(ctx context.Context, requesterID uuid.UUID, req *dto.GetSeriesPaymentsRequest) (*dto.GetSeriesPaymentsResponse, error) {
	ser, err := s.repo.GetSeriesByID(ctx, req.SeriesID)
	if err != nil {
		return nil, err
	}
	clubID, clubState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if clubID == nil || *clubID != ser.ClubID || !canManageClub(clubState) {
		return nil, errorz.PermissionDenied
	}
	if ser.PriceRub <= 0 {
		return &dto.GetSeriesPaymentsResponse{PaidProfileIDs: []uuid.UUID{}}, nil
	}
	ids, err := s.repo.ListPaidSeriesParticipants(ctx, req.SeriesID)
	if err != nil {
		return nil, err
	}
	return &dto.GetSeriesPaymentsResponse{PaidProfileIDs: ids}, nil
}

func (s *Service) SetPayment(ctx context.Context, requesterID uuid.UUID, req *dto.SetSeriesPaymentRequest) error {
	ser, err := s.repo.GetSeriesByID(ctx, req.SeriesID)
	if err != nil {
		return err
	}
	clubID, clubState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return err
	}
	if clubID == nil || *clubID != ser.ClubID || !canManageClub(clubState) {
		return errorz.PermissionDenied
	}
	if ser.PriceRub <= 0 {
		return errorz.InvalidRequest
	}
	isParticipant, err := s.repo.IsSeriesParticipant(ctx, req.SeriesID, req.ProfileID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errorz.InvalidRequest
	}
	return s.repo.SetSeriesParticipantPaid(ctx, req.SeriesID, req.ProfileID, req.Paid)
}

func (s *Service) Join(ctx context.Context, req *dto.JoinSeriesRequest) error {
	ser, err := s.repo.GetSeriesByID(ctx, req.SeriesID)
	if err != nil {
		return err
	}
	banned, err := s.repo.IsProfileBannedInClub(ctx, req.ProfileID, ser.ClubID)
	if err != nil {
		return err
	}
	if banned {
		return errorz.Unauthorized
	}

	if ser.IsClubOnly {
		clubID, _, err := s.repo.GetProfileClubState(ctx, req.ProfileID)
		if err != nil {
			return err
		}
		if clubID == nil || *clubID != ser.ClubID {
			return errorz.Unauthorized
		}
	}

	if ser.IsClosed {
		return errorz.SeriesJoinClosed
	}

	maxParticipants := maxParticipantsForGameType(ser.GameType)
	count, err := s.repo.CountSeriesParticipants(ctx, req.SeriesID)
	if err != nil {
		return err
	}
	if count >= maxParticipants {
		return errorz.InvalidRequest
	}

	return s.repo.AddSeriesParticipant(ctx, req.SeriesID, req.ProfileID)
}

func (s *Service) ListJudges(ctx context.Context, seriesID uuid.UUID) (*dto.GetSeriesJudgesResponse, error) {
	judges, err := s.repo.ListSeriesJudges(ctx, seriesID)
	if err != nil {
		return nil, err
	}
	out := make([]*dto.SeriesJudge, 0, len(judges))
	for _, j := range judges {
		out = append(out, judgeToDTO(j))
	}
	return &dto.GetSeriesJudgesResponse{Items: out}, nil
}

func (s *Service) SetJudge(ctx context.Context, requesterID uuid.UUID, req *dto.SetSeriesJudgeRequest) error {
	ser, err := s.repo.GetSeriesByID(ctx, req.SeriesID)
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
	isParticipant, err := s.repo.IsSeriesParticipant(ctx, req.SeriesID, req.ProfileID)
	if err != nil {
		return err
	}
	isAlreadyJudge, err := s.repo.IsSeriesJudge(ctx, req.SeriesID, req.ProfileID)
	if err != nil {
		return err
	}
	if !isParticipant && !isAlreadyJudge {
		return errorz.InvalidRequest
	}
	if err := s.repo.SetSeriesJudge(ctx, req.SeriesID, req.ProfileID, model.JudgeRole(req.Role)); err != nil {
		return err
	}
	// Remove from participants — judges manage games and are not counted as participants
	return s.repo.RemoveSeriesParticipant(ctx, req.SeriesID, req.ProfileID)
}

func (s *Service) RemoveJudge(ctx context.Context, requesterID uuid.UUID, req *dto.RemoveSeriesJudgeRequest) error {
	ser, err := s.repo.GetSeriesByID(ctx, req.SeriesID)
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
	if err := s.repo.RemoveSeriesJudge(ctx, req.SeriesID, req.ProfileID); err != nil {
		return err
	}
	// Add back to participants when judge role is revoked
	return s.repo.AddSeriesParticipant(ctx, req.SeriesID, req.ProfileID)
}

func (s *Service) Leave(ctx context.Context, req *dto.LeaveSeriesRequest) error {
	ser, err := s.repo.GetSeriesByID(ctx, req.SeriesID)
	if err != nil {
		return err
	}
	isJudge, err := s.repo.IsSeriesJudge(ctx, req.SeriesID, req.ProfileID)
	if err != nil {
		return err
	}
	if isJudge {
		// Judges leave by resigning their role; no is_closed restriction applies
		return s.repo.RemoveSeriesJudge(ctx, req.SeriesID, req.ProfileID)
	}
	if ser.IsClosed {
		return errorz.SeriesJoinClosed
	}
	return s.repo.RemoveSeriesParticipant(ctx, req.SeriesID, req.ProfileID)
}

func (s *Service) GetLeaderboard(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesLeaderboardRequest) (*dto.GetSeriesLeaderboardResponse, error) {
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

	items, total, err := s.repo.ListSeriesLeaderboard(ctx, req.SeriesID, limit, offset)
	if err != nil {
		return nil, err
	}

	outItems := make([]*dto.LeaderboardRow, 0, len(items))
	for _, it := range items {
		outItems = append(outItems, &dto.LeaderboardRow{
			ProfileID:     it.ProfileID,
			GuestID:       it.GuestID,
			GuestNickname: it.GuestNickname,
			Points:        it.Points,
		})
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

func (s *Service) IsProfileBannedInClub(ctx context.Context, profileID uuid.UUID, clubID uuid.UUID) (bool, error) {
	return s.repo.IsProfileBannedInClub(ctx, profileID, clubID)
}
