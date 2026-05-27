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
	List(ctx context.Context, query *string, limit, offset int) ([]*model.Club, int, error)
	Update(ctx context.Context, id uuid.UUID, patch model.ClubUpdatePatch) (*model.Club, error)
	Delete(ctx context.Context, id uuid.UUID) error

	SetProfileClub(ctx context.Context, profileID uuid.UUID, clubID *uuid.UUID, state types.ClubState) error
	ListMembers(ctx context.Context, clubID uuid.UUID, query *string, clubState *types.ClubState, limit, offset int) ([]*model.User, int, error)
	ListGames(ctx context.Context, clubID uuid.UUID, limit, offset int) ([]*model.Game, []string, int, error)
	ListBannedProfiles(ctx context.Context, clubID uuid.UUID, query *string, limit, offset int) ([]*model.User, int, error)
	UnbanProfileInClub(ctx context.Context, profileID uuid.UUID, clubID uuid.UUID) error

	GetProfileClubState(ctx context.Context, profileID uuid.UUID) (clubID *uuid.UUID, state types.ClubState, err error)
	SetMemberState(ctx context.Context, profileID uuid.UUID, clubID uuid.UUID, state types.ClubState) error
	TransferPresidency(ctx context.Context, clubID uuid.UUID, fromPresidentID uuid.UUID, toPresidentID uuid.UUID) error
	IsProfileBannedInClub(ctx context.Context, profileID uuid.UUID, clubID uuid.UUID) (bool, error)
	BanProfileInClub(ctx context.Context, profileID uuid.UUID, clubID uuid.UUID) error
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

	items, total, err := s.repo.List(ctx, req.Query, limit, offset)
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

func canManageClub(state types.ClubState) bool {
	return state == types.ClubStateLeader || state == types.ClubStatePresident
}

func canLeaderManageTarget(targetState types.ClubState) bool {
	return targetState == types.ClubStateMember || targetState == types.ClubStateResident
}

func (s *service) UpdateByManager(ctx context.Context, requesterID uuid.UUID, req *dto.UpdateClubRequest) (*dto.UpdateClubResponse, error) {
	clubID, clubState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if clubID == nil || *clubID != req.ID || !canManageClub(clubState) {
		return nil, errorz.Unauthorized
	}
	return s.Update(ctx, req)
}

func (s *service) Delete(ctx context.Context, req *dto.DeleteClubRequest) error {
	return s.repo.Delete(ctx, req.ID)
}

func (s *service) DeleteByManager(ctx context.Context, requesterID uuid.UUID, req *dto.DeleteClubRequest) error {
	requesterClubID, requesterState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return err
	}
	if requesterClubID == nil || *requesterClubID != req.ID || requesterState != types.ClubStatePresident {
		return errorz.Unauthorized
	}
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

	items, total, err := s.repo.ListMembers(ctx, req.ClubID, req.Query, req.ClubState, limit, offset)
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

func (s *service) GetGames(ctx context.Context, req *dto.GetClubGamesRequest) (*dto.GetClubGamesResponse, error) {
	limit := 10
	offset := 0
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}

	games, seriesNames, totalItems, err := s.repo.ListGames(ctx, req.ClubID, limit, offset)
	if err != nil {
		return nil, err
	}

	items := make([]*dto.PlayerGame, 0, len(games))
	for i, g := range games {
		items = append(items, &dto.PlayerGame{
			ID:         g.ID,
			SeriesID:   g.SeriesID,
			SeriesName: seriesNames[i],
			Name:       g.Name,
			Number:     g.Number,
			Status:     g.Status,
			CreatedAt:  g.CreatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	currentPage := (offset / limit) + 1
	if totalPages == 0 {
		totalPages = 1
		currentPage = 1
	}

	return &dto.GetClubGamesResponse{
		Items: items,
		Pagination: dto.PaginationInfo{
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: currentPage,
			HasNext:     offset+limit < totalItems,
			HasPrevious: offset > 0,
		},
	}, nil
}

func (s *service) GetBans(ctx context.Context, requesterID uuid.UUID, req *dto.GetClubBansRequest) (*dto.GetClubBansResponse, error) {
	requesterClubID, requesterState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if requesterClubID == nil || *requesterClubID != req.ClubID || !canManageClub(requesterState) {
		return nil, errorz.Unauthorized
	}

	limit := 10
	offset := 0
	if req.Limit != nil {
		limit = *req.Limit
	}
	if req.Offset != nil {
		offset = *req.Offset
	}

	items, total, err := s.repo.ListBannedProfiles(ctx, req.ClubID, req.Query, limit, offset)
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
	return &dto.GetClubBansResponse{
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

func (s *service) UnbanMember(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, memberID uuid.UUID) error {
	requesterClubID, requesterState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return err
	}
	if requesterClubID == nil || *requesterClubID != clubID || !canManageClub(requesterState) {
		return errorz.Unauthorized
	}
	return s.repo.UnbanProfileInClub(ctx, memberID, clubID)
}

func (s *service) Join(ctx context.Context, req *dto.JoinClubRequest) error {
	banned, err := s.repo.IsProfileBannedInClub(ctx, req.ProfileID, req.ClubID)
	if err != nil {
		return err
	}
	if banned {
		return errorz.ClubBanned
	}
	currentClubID, currentState, err := s.repo.GetProfileClubState(ctx, req.ProfileID)
	if err != nil {
		return err
	}
	if currentClubID != nil && *currentClubID == req.ClubID && currentState != types.ClubStateNone {
		return errorz.AlreadyInThisClub
	}
	if currentClubID != nil && *currentClubID != req.ClubID && currentState != types.ClubStateNone {
		return errorz.AlreadyInOtherClub
	}
	return s.repo.SetProfileClub(ctx, req.ProfileID, &req.ClubID, types.ClubStateMember)
}

func (s *service) Leave(ctx context.Context, req *dto.LeaveClubRequest) error {
	currentClubID, currentState, err := s.repo.GetProfileClubState(ctx, req.ProfileID)
	if err != nil {
		return err
	}
	if currentClubID == nil || currentState == types.ClubStateNone {
		return errorz.InvalidRequest
	}
	if currentState == types.ClubStateLeader || currentState == types.ClubStatePresident {
		return errorz.ClubSelfAction
	}
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
	if requesterID == memberID {
		return errorz.ClubSelfAction
	}
	memberClubID, memberState, err := s.repo.GetProfileClubState(ctx, memberID)
	if err != nil {
		return err
	}
	if memberClubID == nil || *memberClubID != clubID || memberState == types.ClubStateNone {
		return errorz.InvalidRequest
	}
	// Transfer presidency: requester becomes leader, selected member becomes president.
	return s.repo.TransferPresidency(ctx, clubID, requesterID, memberID)
}

func (s *service) SetMemberRole(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, memberID uuid.UUID, state types.ClubState) error {
	requesterClubID, requesterState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return err
	}
	if requesterClubID == nil || *requesterClubID != clubID || !canManageClub(requesterState) {
		return errorz.Unauthorized
	}
	if requesterID == memberID {
		return errorz.ClubSelfAction
	}
	if state != types.ClubStateMember && state != types.ClubStateLeader && state != types.ClubStateResident {
		return errorz.InvalidRequest
	}
	memberClubID, memberState, err := s.repo.GetProfileClubState(ctx, memberID)
	if err != nil {
		return err
	}
	if memberClubID == nil || *memberClubID != clubID || memberState == types.ClubStateNone {
		return errorz.InvalidRequest
	}
	if requesterState == types.ClubStateLeader {
		if !canLeaderManageTarget(memberState) {
			return errorz.ClubRoleRestricted
		}
		if state == types.ClubStateLeader {
			return errorz.ClubRoleRestricted
		}
	}
	return s.repo.SetMemberState(ctx, memberID, clubID, state)
}

func (s *service) KickFromClub(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, memberID uuid.UUID) error {
	requesterClubID, requesterState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return err
	}
	if requesterClubID == nil || *requesterClubID != clubID || !canManageClub(requesterState) {
		return errorz.Unauthorized
	}
	if requesterID == memberID {
		return errorz.ClubSelfAction
	}
	memberClubID, memberState, err := s.repo.GetProfileClubState(ctx, memberID)
	if err != nil {
		return err
	}
	if memberClubID == nil || *memberClubID != clubID || memberState == types.ClubStateNone {
		return errorz.InvalidRequest
	}
	if requesterState == types.ClubStateLeader && !canLeaderManageTarget(memberState) {
		return errorz.ClubRoleRestricted
	}
	return s.repo.SetProfileClub(ctx, memberID, nil, types.ClubStateNone)
}

func (s *service) KickMember(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, memberID uuid.UUID) error {
	return s.KickFromClub(ctx, requesterID, clubID, memberID)
}

func (s *service) BlockMember(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, memberID uuid.UUID) error {
	if err := s.KickFromClub(ctx, requesterID, clubID, memberID); err != nil {
		return err
	}
	return s.repo.BanProfileInClub(ctx, memberID, clubID)
}

func (s *service) BlockProfile(ctx context.Context, requesterID uuid.UUID, clubID uuid.UUID, profileID uuid.UUID) error {
	requesterClubID, requesterState, err := s.repo.GetProfileClubState(ctx, requesterID)
	if err != nil {
		return err
	}
	if requesterClubID == nil || *requesterClubID != clubID || !canManageClub(requesterState) {
		return errorz.Unauthorized
	}
	targetClubID, targetState, err := s.repo.GetProfileClubState(ctx, profileID)
	if err != nil {
		return err
	}
	if targetClubID != nil && *targetClubID == clubID && targetState != types.ClubStateNone {
		if requesterID == profileID {
			return errorz.ClubSelfAction
		}
		if requesterState == types.ClubStateLeader && !canLeaderManageTarget(targetState) {
			return errorz.ClubRoleRestricted
		}
		if err := s.repo.SetProfileClub(ctx, profileID, nil, types.ClubStateNone); err != nil {
			return err
		}
	}
	return s.repo.BanProfileInClub(ctx, profileID, clubID)
}
