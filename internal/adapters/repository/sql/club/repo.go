package club

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) (*Repo, error) {
	r := &Repo{db: db}
	if err := r.ensureSchema(context.Background()); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Repo) ensureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS clubs (
	id uuid PRIMARY KEY,
	creator_id uuid NOT NULL,
	name text NOT NULL,
	description text NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now()
);

-- FK to profiles for creator_id (profiles must exist)
DO $$
BEGIN
	ALTER TABLE clubs
		ADD CONSTRAINT clubs_creator_id_fk
		FOREIGN KEY (creator_id) REFERENCES profiles(id) ON DELETE RESTRICT;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

ALTER TABLE profiles ADD COLUMN IF NOT EXISTS club_id uuid NULL;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS club_state smallint NOT NULL DEFAULT 0;
DO $$
BEGIN
	ALTER TABLE profiles
		ADD CONSTRAINT profiles_club_id_fk
		FOREIGN KEY (club_id) REFERENCES clubs(id) ON DELETE SET NULL;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;
`)
	return err
}

func (r *Repo) Create(ctx context.Context, c model.Club) (*model.Club, error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	row := r.db.QueryRowContext(ctx, `
INSERT INTO clubs (id, creator_id, name, description)
VALUES ($1,$2,$3,$4)
RETURNING id, creator_id, name, description, created_at, updated_at
`, c.ID, c.CreatorID, c.Name, ptrToNullString(c.Description))

	var out model.Club
	var desc sql.NullString
	if err := row.Scan(&out.ID, &out.CreatorID, &out.Name, &desc, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	return &out, nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*model.Club, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, creator_id, name, description, created_at, updated_at
FROM clubs
WHERE id=$1
`, id)

	var out model.Club
	var desc sql.NullString
	if err := row.Scan(&out.ID, &out.CreatorID, &out.Name, &desc, &out.CreatedAt, &out.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.ClubNotFound
		}
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	return &out, nil
}

func (r *Repo) List(ctx context.Context, limit, offset int) ([]*model.Club, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM clubs`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT id, creator_id, name, description, created_at, updated_at
FROM clubs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.Club
	for rows.Next() {
		var c model.Club
		var desc sql.NullString
		if err := rows.Scan(&c.ID, &c.CreatorID, &c.Name, &desc, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		c.Description = nullStringToPtr(desc)
		out = append(out, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *Repo) Update(ctx context.Context, id uuid.UUID, patch model.ClubUpdatePatch) (*model.Club, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	next := *current
	if patch.Name != nil {
		next.Name = *patch.Name
	}
	if patch.Description != nil {
		next.Description = patch.Description
	}

	row := r.db.QueryRowContext(ctx, `
UPDATE clubs
SET name=$2,
    description=$3,
    updated_at=now()
WHERE id=$1
RETURNING id, creator_id, name, description, created_at, updated_at
`, id, next.Name, ptrToNullString(next.Description))

	var out model.Club
	var desc sql.NullString
	if err := row.Scan(&out.ID, &out.CreatorID, &out.Name, &desc, &out.CreatedAt, &out.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.ClubNotFound
		}
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	return &out, nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM clubs WHERE id=$1`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errorz.ClubNotFound
	}
	return nil
}

func (r *Repo) SetProfileClub(ctx context.Context, profileID uuid.UUID, clubID *uuid.UUID, state types.ClubState) error {
	if clubID == nil || *clubID == uuid.Nil {
		state = types.ClubStateNone
	}
	_, err := r.db.ExecContext(ctx, `UPDATE profiles SET club_id=$2, club_state=$3, updated_at=now() WHERE id=$1`, profileID, ptrToNullUUID(clubID), int16(state))
	return err
}

func (r *Repo) ListMembers(ctx context.Context, clubID uuid.UUID, limit, offset int) ([]*model.Profile, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM profiles WHERE club_id=$1`, clubID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role, created_at, updated_at
FROM profiles
WHERE club_id=$1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`, clubID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.Profile
	for rows.Next() {
		var p model.Profile
		var desc sql.NullString
		var clubIDRaw sql.NullString
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
			&clubIDRaw,
			&clubState,
			&role,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		p.Description = nullStringToPtr(desc)
		p.ClubID = nullStringToUUIDPtr(clubIDRaw)
		p.ClubState = types.ClubState(clubState)
		p.Role = types.Role(role)
		out = append(out, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *Repo) GetProfileClubState(ctx context.Context, profileID uuid.UUID) (clubID *uuid.UUID, state types.ClubState, err error) {
	row := r.db.QueryRowContext(ctx, `SELECT club_id, club_state FROM profiles WHERE id=$1`, profileID)
	var clubIDRaw sql.NullString
	var clubState int16
	if err := row.Scan(&clubIDRaw, &clubState); err != nil {
		return nil, types.ClubStateNone, err
	}
	return nullStringToUUIDPtr(clubIDRaw), types.ClubState(clubState), nil
}

func (r *Repo) SetMemberState(ctx context.Context, profileID uuid.UUID, clubID uuid.UUID, state types.ClubState) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE profiles
SET club_state=$3, updated_at=now()
WHERE id=$1 AND club_id=$2
`, profileID, clubID, int16(state))
	return err
}

func ptrToNullString(p *string) any {
	if p == nil {
		return nil
	}
	return *p
}

func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

func ptrToNullUUID(p *uuid.UUID) any {
	if p == nil || *p == uuid.Nil {
		return nil
	}
	return p.String()
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
