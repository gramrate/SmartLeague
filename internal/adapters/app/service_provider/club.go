package service_provider

import (
	clubRepo "SmartLeague/internal/adapters/repository/sql/club"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/service/club"
	"context"

	"github.com/google/uuid"
)

type clubService interface {
	Create(ctx context.Context, creatorID uuid.UUID, req *dto.CreateClubRequest) (*dto.CreateClubResponse, error)
	GetByID(ctx context.Context, req *dto.GetClubRequest) (*dto.GetClubResponse, error)
	GetAll(ctx context.Context, req *dto.GetAllClubsRequest) (*dto.GetAllClubsResponse, error)
	Update(ctx context.Context, req *dto.UpdateClubRequest) (*dto.UpdateClubResponse, error)
	Delete(ctx context.Context, req *dto.DeleteClubRequest) error
	GetMembers(ctx context.Context, req *dto.GetClubMembersRequest) (*dto.GetClubMembersResponse, error)
	Join(ctx context.Context, req *dto.JoinClubRequest) error
	Leave(ctx context.Context, req *dto.LeaveClubRequest) error
	SetLeader(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, memberID uuid.UUID) error
}

func (s *ServiceProvider) ClubService() clubService {
	if s.clubService == nil {
		repo, err := clubRepo.NewRepo(s.SQLDB())
		if err != nil {
			s.Logger().Panicf("failed to init club repo: %v", err)
		}
		s.clubService = club.NewService(repo)
	}
	return s.clubService
}
