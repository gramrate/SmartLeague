package user

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/dto"
	"context"
	"errors"
)

// UpdateCurrent updates current user by ID.
func (s *userService) UpdateCurrent(ctx context.Context, req *dto.UpdateCurrentUserRequest) (*dto.UpdateCurrentUserResponse, error) {
	userToUpdate, err := s.userRepo.GetById(ctx, req.ID)
	switch {
	case errors.Is(err, errorz.UserNotFound):
		return nil, errorz.UserNotFound
	case err != nil:
		return nil, err
	}
	if req.Name != nil {
		userToUpdate.Name = *req.Name
	}
	if req.Nickname != nil {
		userToUpdate.Nickname = *req.Nickname
	}
	if req.ShowName != nil {
		userToUpdate.ShowName = *req.ShowName
	}
	if req.Description != nil {
		userToUpdate.Description = req.Description
	}
	if req.ClubID != nil {
		userToUpdate.ClubID = req.ClubID
	}

	updatedUser, err := s.userRepo.Update(ctx, *userToUpdate)
	switch {
	case errors.Is(err, errorz.UserNotFound):
		return nil, errorz.UserNotFound
	case err != nil:
		return nil, err
	}

	resp := dto.UpdateCurrentUserResponse(toDTO(updatedUser))
	return &resp, nil
}
