package user

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/dto"
	"context"
	"errors"
)

// Delete delete the user by ID.
func (s *userService) Delete(ctx context.Context, req *dto.DeleteUserRequest) error {
	err := s.userRepo.Delete(ctx, req.ID)
	switch {
	case errors.Is(err, errorz.UserNotFound):
		return errorz.UserNotFound
	case err != nil:
		return err
	}

	return nil
}
