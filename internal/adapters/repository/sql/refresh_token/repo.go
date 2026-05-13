package refresh_token

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/pkg/ent"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetByUserID(ctx context.Context, userID uuid.UUID) (*ent.RefreshToken, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, jti FROM refresh_tokens WHERE user_id=$1`, userID)
	var out ent.RefreshToken
	if err := row.Scan(&out.ID, &out.Jti); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.TokenNotFound
		}
		return nil, err
	}
	out.Edges.User = &ent.User{ID: userID}
	return &out, nil
}

func (r *Repo) Update(ctx context.Context, entity ent.RefreshToken) (*ent.RefreshToken, error) {
	userID := uuid.Nil
	if entity.Edges.User != nil {
		userID = entity.Edges.User.ID
	}
	if userID == uuid.Nil {
		return nil, errorz.TokenNotFound
	}

	res, err := r.db.ExecContext(ctx, `UPDATE refresh_tokens SET jti=$2, updated_at=now() WHERE user_id=$1`, userID, entity.Jti)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, errorz.TokenNotFound
	}
	entity.Edges.User = &ent.User{ID: userID}
	return &entity, nil
}

func (r *Repo) Upsert(ctx context.Context, entity ent.RefreshToken) (*ent.RefreshToken, error) {
	userID := uuid.Nil
	if entity.Edges.User != nil {
		userID = entity.Edges.User.ID
	}
	if userID == uuid.Nil {
		return nil, errorz.TokenNotFound
	}
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}

	_, err := r.db.ExecContext(ctx, `
INSERT INTO refresh_tokens (id, user_id, jti)
VALUES ($1,$2,$3)
ON CONFLICT (user_id) DO UPDATE SET
  jti=excluded.jti,
  updated_at=now()
`, entity.ID, userID, entity.Jti)
	if err != nil {
		return nil, err
	}
	entity.Edges.User = &ent.User{ID: userID}
	return &entity, nil
}

