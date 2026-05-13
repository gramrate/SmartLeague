package profile

import (
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"SmartLeague/internal/domain/utils/password"
	"context"
	"strings"
)

func (s *service) UpdateEach(ctx context.Context, req *dto.UpdateEachProfileRequest) (*dto.UpdateEachProfileResponse, error) {
	var passwordHash *string
	if req.Password != nil {
		h, err := password.PasswordHash(*req.Password)
		if err != nil {
			return nil, err
		}
		passwordHash = &h
	}

	var role *types.Role
	if req.Role != nil {
		role = req.Role
	}

	var email *string
	if req.Email != nil {
		e := normalizeEmail(*req.Email)
		email = &e
	}

	var nickname *string
	if req.Nickname != nil {
		n := strings.TrimSpace(*req.Nickname)
		nickname = &n
	}

	var name *string
	if req.Name != nil {
		n := strings.TrimSpace(*req.Name)
		name = &n
	}

	updated, err := s.repo.Update(ctx, req.ID, model.ProfileUpdatePatch{
		Nickname:     nickname,
		Name:         name,
		ShowName:     req.ShowName,
		Description:  req.Description,
		ClubID:       req.ClubID,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	})
	if err != nil {
		return nil, err
	}
	resp := dto.UpdateEachProfileResponse(*toDTO(updated))
	return &resp, nil
}
