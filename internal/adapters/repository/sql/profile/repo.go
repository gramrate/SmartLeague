package profile

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) (*Repo, error) {
	return &Repo{db: db}, nil
}

type Profile struct {
	model.Profile
}

func (r *Repo) Create(ctx context.Context, p model.Profile) (*model.Profile, error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	row := r.db.QueryRowContext(ctx, `
INSERT INTO profiles (id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
RETURNING id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role, created_at, updated_at
`,
		p.ID,
		p.Nickname,
		p.Name,
		p.ShowName,
		ptrToNullString(p.Description),
		p.Email,
		p.PasswordHash,
		ptrToNullUUID(p.ClubID),
		int16(p.ClubState),
		int16(p.Role),
	)

	var desc sql.NullString
	var clubID sql.NullString
	var clubState int16
	var role int16
	var out model.Profile
	if err := row.Scan(
		&out.ID,
		&out.Nickname,
		&out.Name,
		&out.ShowName,
		&desc,
		&out.Email,
		&out.PasswordHash,
		&clubID,
		&clubState,
		&role,
		&out.CreatedAt,
		&out.UpdatedAt,
	); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, errorz.EmailAlreadyExist
		}
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	out.ClubID = nullStringToUUIDPtr(clubID)
	out.ClubState = types.ClubState(clubState)
	out.Role = types.Role(role)
	return &out, nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*model.Profile, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role, created_at, updated_at
FROM profiles
WHERE id = $1
`, id)

	var out model.Profile
	var desc sql.NullString
	var clubID sql.NullString
	var clubState int16
	var role int16
	if err := row.Scan(
		&out.ID,
		&out.Nickname,
		&out.Name,
		&out.ShowName,
		&desc,
		&out.Email,
		&out.PasswordHash,
		&clubID,
		&clubState,
		&role,
		&out.CreatedAt,
		&out.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.ProfileNotFound
		}
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	out.ClubID = nullStringToUUIDPtr(clubID)
	out.ClubState = types.ClubState(clubState)
	out.Role = types.Role(role)
	return &out, nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (*model.Profile, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role, created_at, updated_at
FROM profiles
WHERE lower(email) = lower($1)
`, email)

	var out model.Profile
	var desc sql.NullString
	var clubID sql.NullString
	var clubState int16
	var role int16
	if err := row.Scan(
		&out.ID,
		&out.Nickname,
		&out.Name,
		&out.ShowName,
		&desc,
		&out.Email,
		&out.PasswordHash,
		&clubID,
		&clubState,
		&role,
		&out.CreatedAt,
		&out.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.ProfileNotFound
		}
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	out.ClubID = nullStringToUUIDPtr(clubID)
	out.ClubState = types.ClubState(clubState)
	out.Role = types.Role(role)
	return &out, nil
}

func (r *Repo) List(ctx context.Context, limit, offset int) ([]*model.Profile, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM profiles`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role, created_at, updated_at
FROM profiles
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.Profile
	for rows.Next() {
		var p model.Profile
		var desc sql.NullString
		var clubID sql.NullString
		var clubState int16
		var role int16
		if err := rows.Scan(
			&p.ID,
			&p.Nickname,
			&p.Name,
			&p.ShowName,
			&desc,
			&p.Email,
			&p.PasswordHash,
			&clubID,
			&clubState,
			&role,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		p.Description = nullStringToPtr(desc)
		p.ClubID = nullStringToUUIDPtr(clubID)
		p.ClubState = types.ClubState(clubState)
		p.Role = types.Role(role)
		out = append(out, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

type UpdatePatch struct {
	model.ProfileUpdatePatch
}

func (r *Repo) Update(ctx context.Context, id uuid.UUID, patch model.ProfileUpdatePatch) (*model.Profile, error) {
	// Read current row for COALESCE behavior.
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	next := *current
	if patch.Nickname != nil {
		next.Nickname = *patch.Nickname
	}
	if patch.Name != nil {
		next.Name = *patch.Name
	}
	if patch.ShowName != nil {
		next.ShowName = *patch.ShowName
	}
	if patch.Description != nil {
		next.Description = patch.Description
	}
	if patch.ClubID != nil {
		next.ClubID = patch.ClubID
	}
	if patch.ClubState != nil {
		next.ClubState = *patch.ClubState
	}
	if patch.Email != nil {
		next.Email = *patch.Email
	}
	if patch.PasswordHash != nil {
		next.PasswordHash = *patch.PasswordHash
	}
	if patch.Role != nil {
		next.Role = *patch.Role
	}

	row := r.db.QueryRowContext(ctx, `
UPDATE profiles
SET nickname=$2,
    name=$3,
    show_name=$4,
    description=$5,
    email=$6,
    password_hash=$7,
    club_id=$8,
    club_state=$9,
    role=$10,
    updated_at=now()
WHERE id=$1
RETURNING id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role, created_at, updated_at
`,
		id,
		next.Nickname,
		next.Name,
		next.ShowName,
		ptrToNullString(next.Description),
		next.Email,
		next.PasswordHash,
		ptrToNullUUID(next.ClubID),
		int16(next.ClubState),
		int16(next.Role),
	)

	var out model.Profile
	var desc sql.NullString
	var clubID sql.NullString
	var clubState int16
	var role int16
	if err := row.Scan(
		&out.ID,
		&out.Nickname,
		&out.Name,
		&out.ShowName,
		&desc,
		&out.Email,
		&out.PasswordHash,
		&clubID,
		&clubState,
		&role,
		&out.CreatedAt,
		&out.UpdatedAt,
	); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, errorz.EmailAlreadyExist
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.ProfileNotFound
		}
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	out.ClubID = nullStringToUUIDPtr(clubID)
	out.ClubState = types.ClubState(clubState)
	out.Role = types.Role(role)
	return &out, nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM profiles WHERE id=$1`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errorz.ProfileNotFound
	}
	return nil
}

func ptrToNullString(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *p, Valid: true}
}

func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

func ptrToNullUUID(p *uuid.UUID) sql.NullString {
	if p == nil || *p == uuid.Nil {
		return sql.NullString{}
	}
	return sql.NullString{String: p.String(), Valid: true}
}

func nullStringToUUIDPtr(ns sql.NullString) *uuid.UUID {
	if !ns.Valid {
		return nil
	}
	id, err := uuid.Parse(ns.String)
	if err != nil {
		return nil
	}
	return &id
}
