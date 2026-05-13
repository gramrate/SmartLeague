package user

import (
	"SmartLeague/internal/domain/types"
	"SmartLeague/pkg/ent"
	"context"
	"github.com/google/uuid"
)

type userRepo interface {
	Create(ctx context.Context, userEntity ent.User) (*ent.User, error)
	GetById(ctx context.Context, id uuid.UUID) (*ent.User, error)
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	GetAllByFilter(
		ctx context.Context,
		limit, offset int,
		role *types.Role,
		query, emailPrefix *string,
	) ([]*ent.User, int, error)
	Update(ctx context.Context, userEntity ent.User) (*ent.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type tokenService interface {
	RevokeAccessToken(ctx context.Context, userID uuid.UUID) error
	GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
	LogoutAllSessions(ctx context.Context, userID uuid.UUID) error
}

type serverConfig interface {
	DevMode() bool
}

// TODO forgot password
type userService struct {
	userRepo     userRepo
	tokenService tokenService
	serverConfig serverConfig
}

func NewUserService(userRepo userRepo, tokenService tokenService, serverConfig serverConfig) *userService {
	return &userService{
		userRepo:     userRepo,
		tokenService: tokenService,
		serverConfig: serverConfig,
	}
}
