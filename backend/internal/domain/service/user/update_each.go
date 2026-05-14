package user

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/types"
	"SmartLeague/internal/domain/utils/password"
	"context"
	"errors"
)

// UpdateEach updates current user by ID.
//
// A user cannot modify another user who has the same or higher role.
//
// A user cannot assign a role that is equal to or higher than their own.
//
// SuperAdmin has no restrictions.
func (s *userService) UpdateEach(ctx context.Context, req *dto.UpdateEachUserRequest) (*dto.UpdateEachUserResponse, error) {
	userToUpdate, err := s.userRepo.GetById(ctx, req.ID)
	switch {
	case errors.Is(err, errorz.UserNotFound):
		return nil, errorz.UserNotFound
	case err != nil:
		return nil, err
	}

	if req.RequesterRole <= userToUpdate.Role && req.RequesterRole != types.RoleSuperAdmin {
		return nil, errorz.PermissionDenied
	}
	if req.Role != nil && req.RequesterRole != types.RoleSuperAdmin {
		if req.RequesterRole <= *req.Role {
			return nil, errorz.PermissionDenied
		}
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
	if req.ClubState != nil {
		userToUpdate.ClubState = *req.ClubState
	}
	if req.Email != nil {
		userToUpdate.Email = *req.Email
	}
	if req.Password != nil {
		passwordHash, err := password.PasswordHash(*req.Password)
		if err != nil {
			return nil, err
		}
		userToUpdate.PasswordHash = passwordHash
	}
	if req.Role != nil {
		userToUpdate.Role = *req.Role
	}

	updatedUser, err := s.userRepo.Update(ctx, *userToUpdate)
	switch {
	case errors.Is(err, errorz.UserNotFound):
		return nil, errorz.UserNotFound
	case err != nil:
		return nil, err
	}

	resp := dto.UpdateEachUserResponse(toDTO(updatedUser))
	return &resp, nil
}
