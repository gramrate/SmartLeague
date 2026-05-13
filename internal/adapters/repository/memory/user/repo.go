package user

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/types"
	"SmartLeague/pkg/ent"
	"context"
	"github.com/google/uuid"
	"strings"
	"sync"
)

type Repo struct {
	mu      sync.RWMutex
	byID    map[uuid.UUID]*ent.User
	byEmail map[string]uuid.UUID
}

func New() *Repo {
	return &Repo{
		byID:    make(map[uuid.UUID]*ent.User),
		byEmail: make(map[string]uuid.UUID),
	}
}

func (r *Repo) Create(ctx context.Context, userEntity ent.User) (*ent.User, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	emailKey := strings.ToLower(strings.TrimSpace(userEntity.Email))
	if _, exists := r.byEmail[emailKey]; exists {
		return nil, errorz.EmailAlreadyExist
	}

	if userEntity.ID == uuid.Nil {
		userEntity.ID = uuid.New()
	}
	userCopy := userEntity
	r.byID[userCopy.ID] = &userCopy
	r.byEmail[emailKey] = userCopy.ID

	return &userCopy, nil
}

func (r *Repo) GetById(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.byID[id]
	if !ok {
		return nil, errorz.UserNotFound
	}
	userCopy := *u
	return &userCopy, nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	emailKey := strings.ToLower(strings.TrimSpace(email))
	id, ok := r.byEmail[emailKey]
	if !ok {
		return nil, errorz.UserNotFound
	}
	u := r.byID[id]
	userCopy := *u
	return &userCopy, nil
}

func (r *Repo) GetAllByFilter(
	ctx context.Context,
	limit, offset int,
	role *types.Role,
	query, emailPrefix *string,
) ([]*ent.User, int, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 50
	}

	var out []*ent.User
	for _, u := range r.byID {
		if role != nil && u.Role != *role {
			continue
		}
		if emailPrefix != nil && !strings.HasPrefix(strings.ToLower(u.Email), strings.ToLower(*emailPrefix)) {
			continue
		}
		if query != nil {
			q := strings.ToLower(strings.TrimSpace(*query))
			if q != "" {
				hay := strings.ToLower(u.Email + " " + u.Name + " " + u.Surname)
				if !strings.Contains(hay, q) {
					continue
				}
			}
		}
		userCopy := *u
		out = append(out, &userCopy)
	}

	total := len(out)
	if offset >= total {
		return []*ent.User{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}
	return out[offset:end], total, nil
}

func (r *Repo) Update(ctx context.Context, userEntity ent.User) (*ent.User, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if userEntity.ID == uuid.Nil {
		return nil, errorz.UserNotFound
	}
	existing, ok := r.byID[userEntity.ID]
	if !ok {
		return nil, errorz.UserNotFound
	}

	// handle email change + uniqueness
	newEmailKey := strings.ToLower(strings.TrimSpace(userEntity.Email))
	oldEmailKey := strings.ToLower(strings.TrimSpace(existing.Email))
	if newEmailKey != oldEmailKey {
		if _, exists := r.byEmail[newEmailKey]; exists {
			return nil, errorz.EmailAlreadyExist
		}
		delete(r.byEmail, oldEmailKey)
		r.byEmail[newEmailKey] = userEntity.ID
	}

	userCopy := userEntity
	r.byID[userCopy.ID] = &userCopy
	return &userCopy, nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.byID[id]
	if !ok {
		return errorz.UserNotFound
	}
	delete(r.byID, id)
	delete(r.byEmail, strings.ToLower(strings.TrimSpace(u.Email)))
	return nil
}

