package club

import (
	"SmartLeague/internal/adapters/controller/api/middleware/auth"
	"SmartLeague/internal/adapters/controller/api/middleware/role"
	"SmartLeague/internal/adapters/controller/api/validator"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/types"
	"context"

	"github.com/go-playground/form"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

type handler struct {
	clubService    clubService
	authMiddleware *auth.Middleware
	roleMiddleware *role.Middleware
	validator      *validator.Validator
	formDecoder    *form.Decoder
}

func NewHandler(
	clubService clubService,
	authMiddleware *auth.Middleware,
	roleMiddleware *role.Middleware,
	validator *validator.Validator,
	formDecoder *form.Decoder,
) *handler {
	return &handler{
		clubService:    clubService,
		authMiddleware: authMiddleware,
		roleMiddleware: roleMiddleware,
		validator:      validator,
		formDecoder:    formDecoder,
	}
}

func (h *handler) Setup(router *echo.Group) {
	router.POST("/club", h.Create, h.authMiddleware.RequireAuth)
	router.GET("/club/:id", h.GetByID)
	router.GET("/club/all", h.GetAll)

	router.PATCH("/club/:id", h.Update, h.authMiddleware.RequireAuth, h.roleMiddleware.RequireRole(types.RoleAdmin))
	router.DELETE("/club/:id", h.Delete, h.authMiddleware.RequireAuth, h.roleMiddleware.RequireRole(types.RoleAdmin))

	router.GET("/club/:id/members", h.GetMembers)
	router.POST("/club/:id/join", h.Join, h.authMiddleware.RequireAuth)
	router.POST("/club/leave", h.Leave, h.authMiddleware.RequireAuth)
	router.POST("/club/:id/leader/:member_id", h.SetLeader, h.authMiddleware.RequireAuth)
}
