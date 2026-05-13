package profile

import (
	"SmartLeague/internal/domain/dto"
	"context"
)

func (s *service) GetByID(ctx context.Context, req *dto.GetProfileRequest) (*dto.GetProfileResponse, error) {
	p, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	resp := dto.GetProfileResponse(*toDTO(p))
	return &resp, nil
}

