package user

import (
	"SmartLeague/internal/domain/dto"
	"context"

	"github.com/google/uuid"
)

func (s *userService) GetUserGames(ctx context.Context, req *dto.GetUserGamesRequest) (*dto.GetUserGamesResponse, error) {
	limit := 10
	if req.Limit != nil {
		limit = *req.Limit
	}
	offset := 0
	if req.Offset != nil {
		offset = *req.Offset
	}

	var viewerClubID *uuid.UUID
	if req.ViewerID != nil {
		if viewer, err := s.userRepo.GetById(ctx, *req.ViewerID); err == nil && viewer.ClubID != nil {
			viewerClubID = viewer.ClubID
		}
	}

	games, seriesNames, totalItems, err := s.userRepo.GetGamesByProfileID(ctx, req.UserID, limit, offset, viewerClubID)
	if err != nil {
		return nil, err
	}

	items := make([]*dto.PlayerGame, 0, len(games))
	for i, g := range games {
		items = append(items, &dto.PlayerGame{
			ID:         g.ID,
			SeriesID:   g.SeriesID,
			SeriesName: seriesNames[i],
			Name:       g.Name,
			Number:     g.Number,
			Status:     g.Status,
			CreatedAt:  g.CreatedAt,
		})
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}
	currentPage := (offset / limit) + 1

	return &dto.GetUserGamesResponse{
		Items: items,
		Pagination: dto.PaginationInfo{
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: currentPage,
			HasNext:     currentPage < totalPages,
			HasPrevious: currentPage > 1,
		},
	}, nil
}

func (s *userService) GetUserSeries(ctx context.Context, req *dto.GetUserSeriesRequest) (*dto.GetUserSeriesResponse, error) {
	limit := 10
	showPast := false
	showClosed := false
	if req.Limit != nil {
		limit = *req.Limit
	}
	offset := 0
	if req.Offset != nil {
		offset = *req.Offset
	}
	if req.ShowPast != nil {
		showPast = *req.ShowPast
	}
	if req.ShowClosed != nil {
		showClosed = *req.ShowClosed
	}

	var viewerClubID *uuid.UUID
	if req.ViewerID != nil {
		if viewer, err := s.userRepo.GetById(ctx, *req.ViewerID); err == nil && viewer.ClubID != nil {
			viewerClubID = viewer.ClubID
		}
	}

	seriesItems, totalItems, err := s.userRepo.GetSeriesByProfileID(
		ctx, req.UserID, limit, offset, req.Query, req.From, req.To, req.IsRating, showPast, showClosed, viewerClubID,
	)
	if err != nil {
		return nil, err
	}

	items := make([]*dto.PlayerSeries, 0, len(seriesItems))
	for _, sItem := range seriesItems {
		ps := &dto.PlayerSeries{
			ID:           sItem.ID,
			Name:         sItem.Name,
			StartAt:      sItem.StartAt,
			EndAt:        sItem.EndAt,
			PriceRub:     sItem.PriceRub,
			IsRating:     sItem.IsRating,
			IsClubOnly:   sItem.IsClubOnly,
			ShowToAll:    sItem.ShowToAll,
			IsClosed:     sItem.IsClosed,
			IsTournament: sItem.IsTournament,
			IsJudge:      sItem.IsJudge,
		}
		if sItem.JudgeRole != nil {
			r := int16(*sItem.JudgeRole)
			ps.JudgeRole = &r
		}
		items = append(items, ps)
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}
	currentPage := (offset / limit) + 1

	return &dto.GetUserSeriesResponse{
		Items: items,
		Pagination: dto.PaginationInfo{
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: currentPage,
			HasNext:     currentPage < totalPages,
			HasPrevious: currentPage > 1,
		},
	}, nil
}
