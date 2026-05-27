package service_provider

import (
	clubRepo "SmartLeague/internal/adapters/repository/sql/club"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/service/club"
	"SmartLeague/internal/domain/types"
	"context"

	"github.com/google/uuid"
)

type clubService interface {
	Create(ctx context.Context, creatorID uuid.UUID, req *dto.CreateClubRequest) (*dto.CreateClubResponse, error)
	GetByID(ctx context.Context, req *dto.GetClubRequest) (*dto.GetClubResponse, error)
	GetAll(ctx context.Context, req *dto.GetAllClubsRequest) (*dto.GetAllClubsResponse, error)
	Update(ctx context.Context, req *dto.UpdateClubRequest) (*dto.UpdateClubResponse, error)
	UpdateByManager(ctx context.Context, requesterID uuid.UUID, req *dto.UpdateClubRequest) (*dto.UpdateClubResponse, error)
	Delete(ctx context.Context, req *dto.DeleteClubRequest) error
	DeleteByManager(ctx context.Context, requesterID uuid.UUID, req *dto.DeleteClubRequest) error
	GetMembers(ctx context.Context, req *dto.GetClubMembersRequest) (*dto.GetClubMembersResponse, error)
	GetGames(ctx context.Context, req *dto.GetClubGamesRequest) (*dto.GetClubGamesResponse, error)
	GetBans(ctx context.Context, requesterID uuid.UUID, req *dto.GetClubBansRequest) (*dto.GetClubBansResponse, error)
	Join(ctx context.Context, req *dto.JoinClubRequest) error
	Leave(ctx context.Context, req *dto.LeaveClubRequest) error
	SetLeader(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, memberID uuid.UUID) error
	SetMemberRole(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, memberID uuid.UUID, state types.ClubState) error
	KickMember(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, memberID uuid.UUID) error
	BlockMember(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, memberID uuid.UUID) error
	UnbanMember(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, memberID uuid.UUID) error
	BlockProfile(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, profileID uuid.UUID) error
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
