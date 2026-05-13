package refresh_token

import (
	"SmartLeague/pkg/ent"
	"context"
	"github.com/google/uuid"
	"sync"
)

type Repo struct {
	mu      sync.RWMutex
	byUser  map[uuid.UUID]*ent.RefreshToken
}

func New() *Repo {
	return &Repo{
		byUser: make(map[uuid.UUID]*ent.RefreshToken),
	}
}

func (r *Repo) GetByUserID(ctx context.Context, userID uuid.UUID) (*ent.RefreshToken, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	tok := r.byUser[userID]
	if tok == nil {
		return nil, nil
	}
	copyTok := *tok
	return &copyTok, nil
}

func (r *Repo) Update(ctx context.Context, entity ent.RefreshToken) (*ent.RefreshToken, error) {
	return r.Upsert(ctx, entity)
}

func (r *Repo) Upsert(ctx context.Context, entity ent.RefreshToken) (*ent.RefreshToken, error) {
	_ = ctx

	if entity.Edges.User == nil {
		return &entity, nil
	}
	userID := entity.Edges.User.ID

	r.mu.Lock()
	defer r.mu.Unlock()

	copyTok := entity
	r.byUser[userID] = &copyTok
	return &copyTok, nil
}

