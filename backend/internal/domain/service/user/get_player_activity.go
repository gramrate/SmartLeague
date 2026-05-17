package user

import (
	"SmartLeague/internal/domain/dto"
	"context"
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

	games, seriesNames, totalItems, err := s.userRepo.GetGamesByProfileID(ctx, req.UserID, limit, offset)
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
	if req.Limit != nil {
		limit = *req.Limit
	}
	offset := 0
	if req.Offset != nil {
		offset = *req.Offset
	}

	seriesItems, totalItems, err := s.userRepo.GetSeriesByProfileID(ctx, req.UserID, limit, offset)
	if err != nil {
		return nil, err
	}

	items := make([]*dto.PlayerSeries, 0, len(seriesItems))
	for _, sItem := range seriesItems {
		items = append(items, &dto.PlayerSeries{
			ID:      sItem.ID,
			Name:    sItem.Name,
			StartAt: sItem.StartAt,
			EndAt:   sItem.EndAt,
		})
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
