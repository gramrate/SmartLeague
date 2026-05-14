package club

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"context"
	"math"

	"github.com/google/uuid"
)

type clubRepo interface {
	Create(ctx context.Context, c model.Club) (*model.Club, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Club, error)
	List(ctx context.Context, limit, offset int) ([]*model.Club, int, error)
	Update(ctx context.Context, id uuid.UUID, patch model.ClubUpdatePatch) (*model.Club, error)
	Delete(ctx context.Context, id uuid.UUID) error

	SetProfileClub(ctx context.Context, profileID uuid.UUID, clubID *uuid.UUID, state types.ClubState) error
	ListMembers(ctx context.Context, clubID uuid.UUID, limit, offset int) ([]*model.User, int, error)

	GetProfileClubState(ctx context.Context, profileID uuid.UUID) (clubID *uuid.UUID, state types.ClubState, err error)
	SetMemberState(ctx context.Context, profileID uuid.UUID, clubID uuid.UUID, state types.ClubState) error
}

type service struct {
	repo clubRepo
}

func NewService(repo clubRepo) *service {
	return &service{repo: repo}
}

func toDTO(c *model.Club) *dto.Club {
	return &dto.Club{
		ID:          c.ID,
		CreatorID:   c.CreatorID,
		Name:        c.Name,
		Description: c.Description,
	}
}

func (s *service) Create(ctx context.Context, creatorID uuid.UUID, req *dto.CreateClubRequest) (*dto.CreateClubResponse, error) {
	created, err := s.repo.Create(ctx, model.Club{
		ID:          uuid.New(),
		CreatorID:   creatorID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	// creator becomes a member
	if err := s.repo.SetProfileClub(ctx, creatorID, &created.ID, types.ClubStatePresident); err != nil {
		return nil, err
	}

	resp := dto.CreateClubResponse(*toDTO(created))
	return &resp, nil
}

func (s *service) GetByID(ctx context.Context, req *dto.GetClubRequest) (*dto.GetClubResponse, error) {
	c, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	resp := dto.GetClubResponse(*toDTO(c))
	return &resp, nil
}

func (s *service) GetAll(ctx context.Context, req *dto.GetAllClubsRequest) (*dto.GetAllClubsResponse, error) {
	limit := 10
	offset := 0
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}

	items, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	outItems := make([]*dto.Club, 0, len(items))
	for _, it := range items {
		outItems = append(outItems, toDTO(it))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	currentPage := (offset / limit) + 1
	if totalPages == 0 {
		totalPages = 1
		currentPage = 1
	}

	return &dto.GetAllClubsResponse{
		Items: outItems,
		Pagination: dto.PaginationInfo{
			TotalItems:  total,
			TotalPages:  totalPages,
			CurrentPage: currentPage,
			HasNext:     offset+limit < total,
			HasPrevious: offset > 0,
		},
	}, nil
}

func (s *service) Update(ctx context.Context, req *dto.UpdateClubRequest) (*dto.UpdateClubResponse, error) {
	updated, err := s.repo.Update(ctx, req.ID, model.ClubUpdatePatch{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	resp := dto.UpdateClubResponse(*toDTO(updated))
	return &resp, nil
}

func (s *service) Delete(ctx context.Context, req *dto.DeleteClubRequest) error {
	return s.repo.Delete(ctx, req.ID)
}

func (s *service) GetMembers(ctx context.Context, req *dto.GetClubMembersRequest) (*dto.GetClubMembersResponse, error) {
	limit := 10
	offset := 0
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}

	items, total, err := s.repo.ListMembers(ctx, req.ClubID, limit, offset)
	if err != nil {
		return nil, err
	}

	outItems := make([]*dto.User, 0, len(items))
	for _, it := range items {
		outItems = append(outItems, &dto.User{
			ID:          it.ID,
			Nickname:    it.Nickname,
			Name:        it.Name,
			ShowName:    it.ShowName,
			Description: it.Description,
			Email:       it.Email,
			ClubID:      it.ClubID,
			ClubState:   it.ClubState,
			Role:        it.Role,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	currentPage := (offset / limit) + 1
	if totalPages == 0 {
		totalPages = 1
		currentPage = 1
	}

	return &dto.GetClubMembersResponse{
		Items: outItems,
		Pagination: dto.PaginationInfo{
			TotalItems:  total,
			TotalPages:  totalPages,
			CurrentPage: currentPage,
			HasNext:     offset+limit < total,
			HasPrevious: offset > 0,
		},
	}, nil
}

func (s *service) Join(ctx context.Context, req *dto.JoinClubRequest) error {
	return s.repo.SetProfileClub(ctx, req.ProfileID, &req.ClubID, types.ClubStateMember)
}

func (s *service) Leave(ctx context.Context, req *dto.LeaveClubRequest) error {
	return s.repo.SetProfileClub(ctx, req.ProfileID, nil, types.ClubStateNone)
}

func (s *service) SetLeader(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, memberID uuid.UUID) error {
	requesterClubID, requesterState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return err
	}
	if requesterClubID == nil || *requesterClubID != clubID || requesterState != types.ClubStatePresident {
		return errorz.Unauthorized
	}
	return s.repo.SetMemberState(ctx, memberID, clubID, types.ClubStateLeader)
}
