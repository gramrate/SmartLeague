package profile

import (
	"SmartLeague/internal/domain/dto"
	"context"
)

func (s *service) Delete(ctx context.Context, req *dto.DeleteProfileRequest) error {
	return s.repo.Delete(ctx, req.ID)
}

