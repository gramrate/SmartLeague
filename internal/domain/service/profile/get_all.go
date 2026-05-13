package profile

import (
	"SmartLeague/internal/domain/dto"
	"context"
	"math"
)

func (s *service) GetAll(ctx context.Context, req *dto.GetAllProfilesRequest) (*dto.GetAllProfilesResponse, error) {
	limit := 10
	offset := 0
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}

	items, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	outItems := make([]*dto.Profile, 0, len(items))
	for _, it := range items {
		outItems = append(outItems, toDTO(it))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	currentPage := (offset / limit) + 1
	if totalPages == 0 {
		totalPages = 1
		currentPage = 1
	}

	return &dto.GetAllProfilesResponse{
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

