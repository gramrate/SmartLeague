package profile

import (
	"SmartLeague/internal/adapters/controller/api/middleware/auth"
	"SmartLeague/internal/adapters/controller/api/middleware/role"
	"SmartLeague/internal/adapters/controller/api/validator"
	"SmartLeague/internal/domain/dto"
	"SmartLeague/internal/domain/types"
	"context"

	"github.com/go-playground/form"
	"github.com/labstack/echo/v4"
)

type profileService interface {
	Create(ctx context.Context, req *dto.CreateProfileRequest) (*dto.CreateProfileResponse, error)
	GetByID(ctx context.Context, req *dto.GetProfileRequest) (*dto.GetProfileResponse, error)
	GetAll(ctx context.Context, req *dto.GetAllProfilesRequest) (*dto.GetAllProfilesResponse, error)
	UpdateCurrent(ctx context.Context, req *dto.UpdateCurrentProfileRequest) (*dto.UpdateCurrentProfileResponse, error)
	UpdateEach(ctx context.Context, req *dto.UpdateEachProfileRequest) (*dto.UpdateEachProfileResponse, error)
	Delete(ctx context.Context, req *dto.DeleteProfileRequest) error
}

type handler struct {
	profileService profileService
	authMiddleware *auth.Middleware
	roleMiddleware *role.Middleware
	validator      *validator.Validator
	formDecoder    *form.Decoder
}

func NewHandler(
	profileService profileService,
	authMiddleware *auth.Middleware,
	roleMiddleware *role.Middleware,
	validator *validator.Validator,
	formDecoder *form.Decoder,
) *handler {
	return &handler{
		profileService: profileService,
		authMiddleware: authMiddleware,
		roleMiddleware: roleMiddleware,
		validator:      validator,
		formDecoder:    formDecoder,
	}
}

func (h *handler) Setup(router *echo.Group) {
	router.POST("/profile", h.Create)

	router.GET("/profile", h.GetMe, h.authMiddleware.RequireAuth)
	router.PATCH("/profile", h.UpdateMe, h.authMiddleware.RequireAuth)
	router.DELETE("/profile", h.DeleteMe, h.authMiddleware.RequireAuth)

	router.GET("/profile/:id", h.GetByID, h.authMiddleware.RequireAuth, h.roleMiddleware.RequireRole(types.RoleAdmin))
	router.GET("/profile/all", h.GetAll, h.authMiddleware.RequireAuth, h.roleMiddleware.RequireRole(types.RoleAdmin))
	router.PATCH("/profile/:id", h.UpdateEach, h.authMiddleware.RequireAuth, h.roleMiddleware.RequireRole(types.RoleAdmin))
	router.DELETE("/profile/:id", h.DeleteEach, h.authMiddleware.RequireAuth, h.roleMiddleware.RequireRole(types.RoleAdmin))
}

