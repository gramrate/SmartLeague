package user

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/dto"
	"context"
	"errors"
)

// GetByID returns the user by ID.
func (s *userService) GetByID(ctx context.Context, req *dto.GetUserRequest) (*dto.GetUserResponse, error) {
	u, err := s.userRepo.GetById(ctx, req.ID)
	switch {
	case errors.Is(err, errorz.UserNotFound):
		return nil, errorz.UserNotFound
	case err != nil:
		return nil, err
	}

	resp := dto.GetUserResponse(toDTO(u))
	return &resp, nil
}
