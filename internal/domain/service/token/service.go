package token

import (
	"SmartLeague/internal/domain/model"
	"context"
	"github.com/google/uuid"
	"time"
)

type accessTokenRepo interface {
	Set(ctx context.Context, userID uuid.UUID, value string, exp time.Time) error
	Get(ctx context.Context, userID uuid.UUID) (string, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}

type refreshTokenRepo interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.RefreshToken, error)
	Update(ctx context.Context, entity model.RefreshToken) (*model.RefreshToken, error)
	Upsert(ctx context.Context, entity model.RefreshToken) (*model.RefreshToken, error)
}

type jwtService interface {
	GenerateToken(userID uuid.UUID, ttl time.Duration) (token string, jti string, err error)
	ParseToken(tokenString string) (userID uuid.UUID, jti string, err error)
}

type jwtConfig interface {
	RefreshTokenExpires() time.Duration
	AccessTokenExpires() time.Duration
}

type tokenService struct {
	refreshTokenRepo refreshTokenRepo
	accessTokenRepo  accessTokenRepo
	jwtService       jwtService
	jwtConfig        jwtConfig
}

func NewTokenService(refreshTokenRepo refreshTokenRepo, accessTokenRepo accessTokenRepo, jwtService jwtService, jwtConfig jwtConfig) *tokenService {
	return &tokenService{
		refreshTokenRepo: refreshTokenRepo,
		accessTokenRepo:  accessTokenRepo,
		jwtService:       jwtService,
		jwtConfig:        jwtConfig,
	}
}
