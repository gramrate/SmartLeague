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

	CreateSeries(ctx context.Context, s model.Series) (*model.Series, error)
	GetSeriesByID(ctx context.Context, id uuid.UUID) (*model.Series, error)
	ListSeriesByClub(ctx context.Context, clubID uuid.UUID, includeClosed bool, limit, offset int) ([]*model.Series, int, error)
	ListAllSeries(ctx context.Context, limit, offset int, showPast, showClosed bool) ([]*model.SeriesListItem, int, error)
	UpdateSeries(ctx context.Context, id uuid.UUID, patch model.SeriesUpdatePatch) (*model.Series, error)
	DeleteSeries(ctx context.Context, id uuid.UUID) error

	AddSeriesParticipant(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID) error
	RemoveSeriesParticipant(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID) error
	CountSeriesParticipants(ctx context.Context, seriesID uuid.UUID) (int, error)
	ListSeriesParticipants(ctx context.Context, seriesID uuid.UUID, limit, offset int, query *string) ([]*model.User, int, error)
	IsSeriesParticipant(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID) (bool, error)

	ListSeriesLeaderboard(ctx context.Context, seriesID uuid.UUID, limit, offset int) ([]*model.LeaderboardRow, int, error)
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

func seriesToDTO(s *model.Series, creatorID *uuid.UUID) *dto.Series {
	return &dto.Series{
		ID:          s.ID,
		ClubID:      s.ClubID,
		CreatorID:   creatorID,
		Name:        s.Name,
		Description: s.Description,
		StartAt:     s.StartAt,
		EndAt:       s.EndAt,
		PriceRub:    s.PriceRub,
		IsClosed:    s.IsClosed,
		GameType:    s.GameType,
		Status:      s.Status,
	}
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
		ID:          uuid.New(),
		ClubID:      *clubID,
		CreatorID:   requesterID,
		Name:        req.Name,
		Description: req.Description,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		PriceRub:    req.PriceRub,
		IsClosed:    req.IsClosed,
		GameType:    types.GameTypeSportMafia,
		Status:      req.Status,
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

	items, total, err := s.repo.ListSeriesByClub(ctx, req.ClubID, includeClosed, limit, offset)
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

func (s *Service) GetAllSeries(ctx context.Context, req *dto.GetAllSeriesRequest) (*dto.GetAllSeriesResponse, error) {
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

	items, total, err := s.repo.ListAllSeries(ctx, limit, offset, showPast, showClosed)
	if err != nil {
		return nil, err
	}

	outItems := make([]*dto.AllSeriesItem, 0, len(items))
	for _, it := range items {
		outItems = append(outItems, &dto.AllSeriesItem{
			ID:          it.ID,
			ClubID:      it.ClubID,
			ClubName:    it.ClubName,
			Name:        it.Name,
			Description: it.Description,
			StartAt:     it.StartAt,
			EndAt:       it.EndAt,
			IsClosed:    it.IsClosed,
			GamesCount:  it.GamesCount,
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
		Name:        req.Name,
		Description: req.Description,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		PriceRub:    req.PriceRub,
		IsClosed:    req.IsClosed,
		Status:      req.Status,
	}

	updated, err := s.repo.UpdateSeries(ctx, req.ID, patch)
	if err != nil {
		return nil, err
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
	_ = requesterID
	if _, err := s.repo.GetSeriesByID(ctx, req.SeriesID); err != nil {
		return nil, err
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

func (s *Service) Join(ctx context.Context, req *dto.JoinSeriesRequest) error {
	ser, err := s.repo.GetSeriesByID(ctx, req.SeriesID)
	if err != nil {
		return err
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

func (s *Service) Leave(ctx context.Context, req *dto.LeaveSeriesRequest) error {
	ser, err := s.repo.GetSeriesByID(ctx, req.SeriesID)
	if err != nil {
		return err
	}
	if ser.IsClosed {
		return errorz.SeriesJoinClosed
	}
	return s.repo.RemoveSeriesParticipant(ctx, req.SeriesID, req.ProfileID)
}

func (s *Service) GetLeaderboard(ctx context.Context, requesterID *uuid.UUID, req *dto.GetSeriesLeaderboardRequest) (*dto.GetSeriesLeaderboardResponse, error) {
	_ = requesterID
	if _, err := s.repo.GetSeriesByID(ctx, req.SeriesID); err != nil {
		return nil, err
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
