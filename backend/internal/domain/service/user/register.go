package user

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"SmartLeague/internal/domain/utils/password"
	"context"
	"errors"
)

// Register returns registered user with token.
func (s *userService) Register(ctx context.Context, req *dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
	passwordHash, err := password.PasswordHash(req.Password)
	if err != nil {
		return nil, err
	}
	user := model.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Name:         req.Name,
		Surname:      req.Surname,
		Role:         types.RoleUser,
	}
	if s.serverConfig.DevMode() && req.Role != nil {
		user.Role = *req.Role
	}
	u, err := s.userRepo.Create(ctx, user)
	switch {
	case errors.Is(err, errorz.EmailAlreadyExist):
		return nil, errorz.EmailAlreadyExist
	case err != nil:
		return nil, err
	}

	token, err := s.tokenService.GenerateRefreshToken(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	err = s.tokenService.RevokeAccessToken(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterUserResponse{
		RefreshToken: token,
		User: dto.User{
			ID:      u.ID,
			Email:   u.Email,
			Name:    u.Name,
			Surname: u.Surname,
			Role:    u.Role,
		},
	}, nil
}
