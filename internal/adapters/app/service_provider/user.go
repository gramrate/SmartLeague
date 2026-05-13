package service_provider

import (
	userMem "SmartLeague/internal/adapters/repository/memory/user"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/service/user"
	"SmartLeague/internal/domain/types"
	"context"
	"github.com/google/uuid"
)

type userService interface {
	Register(ctx context.Context, req *dto.RegisterUserRequest) (*dto.RegisterUserResponse, error)
	Login(ctx context.Context, req *dto.LoginUserRequest) (*dto.LoginUserResponse, error)
	ChangePassword(ctx context.Context, req *dto.ChangePasswordRequest) (*dto.ChangePasswordResponse, error)
	GetByID(ctx context.Context, req *dto.GetUserRequest) (*dto.GetUserResponse, error)
	GetRoleByID(ctx context.Context, userID uuid.UUID) (types.Role, error)
	GetAllByFilter(ctx context.Context, req *dto.GetAllByFilterUsersRequest) (*dto.GetAllByFilterUsersResponse, error)
	UpdateCurrent(ctx context.Context, req *dto.UpdateCurrentUserRequest) (*dto.UpdateCurrentUserResponse, error)
	UpdateEach(ctx context.Context, req *dto.UpdateEachUserRequest) (*dto.UpdateEachUserResponse, error)
	Delete(ctx context.Context, req *dto.DeleteUserRequest) error
	Logout(ctx context.Context, req *dto.LogoutRequest) error
}

func (s *ServiceProvider) UserService() userService {
	if s.userService == nil {
		s.userService = user.NewUserService(userMem.New(), s.TokenService(), s.ServerConfig())
	}
	return s.userService
}
