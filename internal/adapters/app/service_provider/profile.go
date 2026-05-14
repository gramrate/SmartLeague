package service_provider

import (
	profileRepo "SmartLeague/internal/adapters/repository/sql/profile"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/service/profile"
	"context"
)

type profileService interface {
	Create(ctx context.Context, req *dto.CreateProfileRequest) (*dto.CreateProfileResponse, error)
	GetByID(ctx context.Context, req *dto.GetProfileRequest) (*dto.GetProfileResponse, error)
	GetAll(ctx context.Context, req *dto.GetAllProfilesRequest) (*dto.GetAllProfilesResponse, error)
	UpdateCurrent(ctx context.Context, req *dto.UpdateCurrentProfileRequest) (*dto.UpdateCurrentProfileResponse, error)
	UpdateEach(ctx context.Context, req *dto.UpdateEachProfileRequest) (*dto.UpdateEachProfileResponse, error)
	Delete(ctx context.Context, req *dto.DeleteProfileRequest) error
}

func (s *ServiceProvider) ProfileService() profileService {
	if s.profileService == nil {
		repo, err := profileRepo.NewRepo(s.SQLDB())
		if err != nil {
			s.Logger().Panicf("failed to init profile repo: %v", err)
		}
		s.profileService = profile.NewService(repo, s.ServerConfig())
	}
	return s.profileService
}
