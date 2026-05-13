package profile

import (
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/model"
	"context"
)

func (s *service) UpdateCurrent(ctx context.Context, req *dto.UpdateCurrentProfileRequest) (*dto.UpdateCurrentProfileResponse, error) {
	updated, err := s.repo.Update(ctx, req.ID, model.ProfileUpdatePatch{
		Nickname:    req.Nickname,
		Name:        req.Name,
		ShowName:    req.ShowName,
		Description: req.Description,
		Club:        req.Club,
	})
	if err != nil {
		return nil, err
	}
	resp := dto.UpdateCurrentProfileResponse(*toDTO(updated))
	return &resp, nil
}
