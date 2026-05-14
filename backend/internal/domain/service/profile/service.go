package profile

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"SmartLeague/internal/domain/utils/password"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type profileRepo interface {
	Create(ctx context.Context, p model.Profile) (*model.Profile, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Profile, error)
	List(ctx context.Context, limit, offset int) ([]*model.Profile, int, error)
	Update(ctx context.Context, id uuid.UUID, patch model.ProfileUpdatePatch) (*model.Profile, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type serverConfig interface {
	DevMode() bool
}

type service struct {
	repo         profileRepo
	serverConfig serverConfig
}

func NewService(repo profileRepo, serverConfig serverConfig) *service {
	return &service{repo: repo, serverConfig: serverConfig}
}

func toDTO(p *model.Profile) *dto.Profile {
	return &dto.Profile{
		ID:          p.ID,
		Nickname:    p.Nickname,
		Name:        p.Name,
		ShowName:    p.ShowName,
		Description: p.Description,
		Email:       p.Email,
		ClubID:      p.ClubID,
		ClubState:   p.ClubState,
		Role:        p.Role,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (s *service) Create(ctx context.Context, req *dto.CreateProfileRequest) (*dto.CreateProfileResponse, error) {
	email := normalizeEmail(req.Email)
	passwordHash, err := password.PasswordHash(req.Password)
	if err != nil {
		return nil, err
	}

	showName := true
	if req.ShowName != nil {
		showName = *req.ShowName
	}
	nickname := ""
	if req.Nickname != nil {
		nickname = strings.TrimSpace(*req.Nickname)
	}

	role := types.RoleUser
	if s.serverConfig.DevMode() && req.Role != nil {
		role = *req.Role
	}

	created, err := s.repo.Create(ctx, model.Profile{
		ID:           uuid.New(),
		Nickname:     nickname,
		Name:         strings.TrimSpace(req.Name),
		ShowName:     showName,
		Description:  req.Description,
		Email:        email,
		PasswordHash: passwordHash,
		ClubID:       req.ClubID,
		Role:         role,
	})
	switch {
	case errors.Is(err, errorz.EmailAlreadyExist):
		return nil, errorz.EmailAlreadyExist
	case err != nil:
		return nil, err
	}

	resp := dto.CreateProfileResponse(*toDTO(created))
	return &resp, nil
}
