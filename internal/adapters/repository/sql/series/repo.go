package series

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"context"
	"database/sql"
	"errors"
	"fmt"

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
CREATE TABLE IF NOT EXISTS series (
	id uuid PRIMARY KEY,
	club_id uuid NOT NULL,
	creator_id uuid NOT NULL,
	name text NOT NULL,
	scoring_rules text NOT NULL,
	start_at timestamptz NOT NULL,
	end_at timestamptz NOT NULL,
	description text NULL,
	price_rub integer NOT NULL DEFAULT 0,
	is_closed boolean NOT NULL DEFAULT false,
	game_type smallint NOT NULL DEFAULT 0,
	status smallint NOT NULL DEFAULT 0,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now()
);

DO $$
BEGIN
	ALTER TABLE series
		ADD CONSTRAINT series_club_id_fk
		FOREIGN KEY (club_id) REFERENCES clubs(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
	ALTER TABLE series
		ADD CONSTRAINT series_creator_id_fk
		FOREIGN KEY (creator_id) REFERENCES profiles(id) ON DELETE RESTRICT;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

CREATE INDEX IF NOT EXISTS series_club_id_idx ON series (club_id);

CREATE TABLE IF NOT EXISTS series_participants (
	series_id uuid NOT NULL,
	profile_id uuid NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (series_id, profile_id)
);

DO $$
BEGIN
	ALTER TABLE series_participants
		ADD CONSTRAINT series_participants_series_fk
		FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
	ALTER TABLE series_participants
		ADD CONSTRAINT series_participants_profile_fk
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS games (
	id uuid PRIMARY KEY,
	series_id uuid NOT NULL,
	name text NOT NULL,
	number integer NOT NULL,
	description text NULL,
	host_id uuid NULL,
	status smallint NOT NULL DEFAULT 0,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	UNIQUE (series_id, number)
);

DO $$
BEGIN
	ALTER TABLE games
		ADD CONSTRAINT games_series_fk
		FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
	ALTER TABLE games
		ADD CONSTRAINT games_host_fk
		FOREIGN KEY (host_id) REFERENCES profiles(id) ON DELETE SET NULL;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

CREATE INDEX IF NOT EXISTS games_series_idx ON games (series_id);

CREATE TABLE IF NOT EXISTS game_participants (
	game_id uuid NOT NULL,
	profile_id uuid NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (game_id, profile_id)
);

DO $$
BEGIN
	ALTER TABLE game_participants
		ADD CONSTRAINT game_participants_game_fk
		FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
	ALTER TABLE game_participants
		ADD CONSTRAINT game_participants_profile_fk
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS game_results (
	game_id uuid NOT NULL,
	profile_id uuid NOT NULL,
	place integer NULL,
	role text NULL,
	best_move boolean NOT NULL DEFAULT false,
	first_killed boolean NOT NULL DEFAULT false,
	compensation integer NOT NULL DEFAULT 0,
	yellow_cards integer NOT NULL DEFAULT 0,
	removed boolean NOT NULL DEFAULT false,
	extra_points integer NOT NULL DEFAULT 0,
	total_points integer NOT NULL DEFAULT 0,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (game_id, profile_id)
);

DO $$
BEGIN
	ALTER TABLE game_results
		ADD CONSTRAINT game_results_game_fk
		FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
	ALTER TABLE game_results
		ADD CONSTRAINT game_results_profile_fk
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE;
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;
`)
	return err
}

func (r *Repo) GetProfileClubState(ctx context.Context, profileID uuid.UUID) (clubID *uuid.UUID, state types.ClubState, err error) {
	row := r.db.QueryRowContext(ctx, `SELECT club_id, club_state FROM profiles WHERE id=$1`, profileID)
	var clubIDRaw sql.NullString
	var clubState int16
	if err := row.Scan(&clubIDRaw, &clubState); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ClubStateNone, errorz.UserNotFound
		}
		return nil, types.ClubStateNone, err
	}
	return nullStringToUUIDPtr(clubIDRaw), types.ClubState(clubState), nil
}

func (r *Repo) CreateSeries(ctx context.Context, s model.Series) (*model.Series, error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	row := r.db.QueryRowContext(ctx, `
INSERT INTO series (id, club_id, creator_id, name, scoring_rules, start_at, end_at, description, price_rub, is_closed, game_type, status)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
RETURNING id, club_id, creator_id, name, scoring_rules, start_at, end_at, description, price_rub, is_closed, game_type, status, created_at, updated_at
`,
		s.ID, s.ClubID, s.CreatorID, s.Name, s.ScoringRules, s.StartAt, s.EndAt,
		ptrToNullString(s.Description), s.PriceRub, s.IsClosed, int16(s.GameType), int16(s.Status),
	)

	var out model.Series
	var desc sql.NullString
	var gameType int16
	var status int16
	if err := row.Scan(
		&out.ID,
		&out.ClubID,
		&out.CreatorID,
		&out.Name,
		&out.ScoringRules,
		&out.StartAt,
		&out.EndAt,
		&desc,
		&out.PriceRub,
		&out.IsClosed,
		&gameType,
		&status,
		&out.CreatedAt,
		&out.UpdatedAt,
	); err != nil {
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	out.GameType = types.GameType(gameType)
	out.Status = types.SeriesStatus(status)
	return &out, nil
}

func (r *Repo) GetSeriesByID(ctx context.Context, id uuid.UUID) (*model.Series, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, club_id, creator_id, name, scoring_rules, start_at, end_at, description, price_rub, is_closed, game_type, status, created_at, updated_at
FROM series
WHERE id=$1
`, id)

	var out model.Series
	var desc sql.NullString
	var gameType int16
	var status int16
	if err := row.Scan(
		&out.ID,
		&out.ClubID,
		&out.CreatorID,
		&out.Name,
		&out.ScoringRules,
		&out.StartAt,
		&out.EndAt,
		&desc,
		&out.PriceRub,
		&out.IsClosed,
		&gameType,
		&status,
		&out.CreatedAt,
		&out.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.SeriesNotFound
		}
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	out.GameType = types.GameType(gameType)
	out.Status = types.SeriesStatus(status)
	return &out, nil
}

func (r *Repo) ListSeriesByClub(ctx context.Context, clubID uuid.UUID, includeClosed bool, limit, offset int) ([]*model.Series, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	where := "club_id=$1"
	if !includeClosed {
		where = where + " AND is_closed=false"
	}

	var total int
	countQuery := fmt.Sprintf("SELECT count(*) FROM series WHERE %s", where)
	if err := r.db.QueryRowContext(ctx, countQuery, clubID).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := fmt.Sprintf(`
SELECT id, club_id, creator_id, name, scoring_rules, start_at, end_at, description, price_rub, is_closed, game_type, status, created_at, updated_at
FROM series
WHERE %s
ORDER BY start_at DESC
LIMIT $2 OFFSET $3
`, where)

	rows, err := r.db.QueryContext(ctx, listQuery, clubID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.Series
	for rows.Next() {
		var s model.Series
		var desc sql.NullString
		var gameType int16
		var status int16
		if err := rows.Scan(
			&s.ID,
			&s.ClubID,
			&s.CreatorID,
			&s.Name,
			&s.ScoringRules,
			&s.StartAt,
			&s.EndAt,
			&desc,
			&s.PriceRub,
			&s.IsClosed,
			&gameType,
			&status,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		s.Description = nullStringToPtr(desc)
		s.GameType = types.GameType(gameType)
		s.Status = types.SeriesStatus(status)
		out = append(out, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *Repo) UpdateSeries(ctx context.Context, id uuid.UUID, patch model.SeriesUpdatePatch) (*model.Series, error) {
	current, err := r.GetSeriesByID(ctx, id)
	if err != nil {
		return nil, err
	}
	next := *current
	if patch.Name != nil {
		next.Name = *patch.Name
	}
	if patch.ScoringRules != nil {
		next.ScoringRules = *patch.ScoringRules
	}
	if patch.StartAt != nil {
		next.StartAt = *patch.StartAt
	}
	if patch.EndAt != nil {
		next.EndAt = *patch.EndAt
	}
	if patch.Description != nil {
		next.Description = patch.Description
	}
	if patch.PriceRub != nil {
		next.PriceRub = *patch.PriceRub
	}
	if patch.IsClosed != nil {
		next.IsClosed = *patch.IsClosed
	}
	if patch.GameType != nil {
		next.GameType = *patch.GameType
	}
	if patch.Status != nil {
		next.Status = *patch.Status
	}

	row := r.db.QueryRowContext(ctx, `
UPDATE series
SET name=$2,
    scoring_rules=$3,
    start_at=$4,
    end_at=$5,
    description=$6,
    price_rub=$7,
    is_closed=$8,
    game_type=$9,
    status=$10,
    updated_at=now()
WHERE id=$1
RETURNING id, club_id, creator_id, name, scoring_rules, start_at, end_at, description, price_rub, is_closed, game_type, status, created_at, updated_at
`,
		id,
		next.Name,
		next.ScoringRules,
		next.StartAt,
		next.EndAt,
		ptrToNullString(next.Description),
		next.PriceRub,
		next.IsClosed,
		int16(next.GameType),
		int16(next.Status),
	)

	var out model.Series
	var desc sql.NullString
	var gameType int16
	var status int16
	if err := row.Scan(
		&out.ID,
		&out.ClubID,
		&out.CreatorID,
		&out.Name,
		&out.ScoringRules,
		&out.StartAt,
		&out.EndAt,
		&desc,
		&out.PriceRub,
		&out.IsClosed,
		&gameType,
		&status,
		&out.CreatedAt,
		&out.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.SeriesNotFound
		}
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	out.GameType = types.GameType(gameType)
	out.Status = types.SeriesStatus(status)
	return &out, nil
}

func (r *Repo) DeleteSeries(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM series WHERE id=$1`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errorz.SeriesNotFound
	}
	return nil
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
