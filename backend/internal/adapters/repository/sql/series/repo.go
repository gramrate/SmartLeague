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
	return &Repo{db: db}, nil
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
VALUES ($1,$2,$3,$4,$5,$6,$7,NULL,$8,$9,$10,$11)
RETURNING id, club_id, creator_id, name, scoring_rules, start_at, end_at, price_rub, is_closed, game_type, status, created_at, updated_at
`,
		s.ID, s.ClubID, s.CreatorID, s.Name, s.Description, s.StartAt, s.EndAt,
		s.PriceRub, s.IsClosed, int16(s.GameType), int16(s.Status),
	)

	var out model.Series
	var gameType int16
	var status int16
	if err := row.Scan(
		&out.ID,
		&out.ClubID,
		&out.CreatorID,
		&out.Name,
		&out.Description,
		&out.StartAt,
		&out.EndAt,
		&out.PriceRub,
		&out.IsClosed,
		&gameType,
		&status,
		&out.CreatedAt,
		&out.UpdatedAt,
	); err != nil {
		return nil, err
	}
	out.GameType = types.GameType(gameType)
	out.Status = types.SeriesStatus(status)
	return &out, nil
}

func (r *Repo) GetSeriesByID(ctx context.Context, id uuid.UUID) (*model.Series, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, club_id, creator_id, name, scoring_rules, start_at, end_at, price_rub, is_closed, game_type, status, created_at, updated_at
FROM series
WHERE id=$1 AND deleted_at IS NULL
`, id)

	var out model.Series
	var gameType int16
	var status int16
	if err := row.Scan(
		&out.ID,
		&out.ClubID,
		&out.CreatorID,
		&out.Name,
		&out.Description,
		&out.StartAt,
		&out.EndAt,
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
	where = where + " AND deleted_at IS NULL"

	var total int
	countQuery := fmt.Sprintf("SELECT count(*) FROM series WHERE %s", where)
	if err := r.db.QueryRowContext(ctx, countQuery, clubID).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := fmt.Sprintf(`
SELECT id, club_id, creator_id, name, scoring_rules, start_at, end_at, price_rub, is_closed, game_type, status, created_at, updated_at
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
		var gameType int16
		var status int16
		if err := rows.Scan(
			&s.ID,
			&s.ClubID,
			&s.CreatorID,
			&s.Name,
			&s.Description,
			&s.StartAt,
			&s.EndAt,
			&s.PriceRub,
			&s.IsClosed,
			&gameType,
			&status,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		s.GameType = types.GameType(gameType)
		s.Status = types.SeriesStatus(status)
		out = append(out, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *Repo) ListAllSeries(ctx context.Context, limit, offset int, showPast, showClosed bool) ([]*model.SeriesListItem, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	where := "s.deleted_at IS NULL"
	if !showClosed {
		where += " AND s.is_closed = false"
	}
	if !showPast {
		where += " AND s.end_at >= now()"
	}

	var total int
	countQuery := fmt.Sprintf(`
SELECT count(*)
FROM series s
WHERE %s
`, where)
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := fmt.Sprintf(`
SELECT
  s.id,
  s.club_id,
  c.name,
  s.name,
  s.scoring_rules,
  s.start_at,
  s.end_at,
  s.is_closed,
  (
    SELECT count(*)
    FROM games g
    WHERE g.series_id = s.id AND g.deleted_at IS NULL AND g.status <> 0
  ) AS games_count
FROM series s
JOIN clubs c ON c.id = s.club_id
WHERE %s
ORDER BY s.start_at DESC
LIMIT $1 OFFSET $2
`, where)
	rows, err := r.db.QueryContext(ctx, listQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]*model.SeriesListItem, 0, limit)
	for rows.Next() {
		var item model.SeriesListItem
		if err := rows.Scan(
			&item.ID,
			&item.ClubID,
			&item.ClubName,
			&item.Name,
			&item.Description,
			&item.StartAt,
			&item.EndAt,
			&item.IsClosed,
			&item.GamesCount,
		); err != nil {
			return nil, 0, err
		}
		out = append(out, &item)
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
	if patch.Description != nil {
		next.Description = *patch.Description
	}
	if patch.StartAt != nil {
		next.StartAt = *patch.StartAt
	}
	if patch.EndAt != nil {
		next.EndAt = *patch.EndAt
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
    description=NULL,
    price_rub=$6,
    is_closed=$7,
    game_type=$8,
    status=$9,
    updated_at=now()
WHERE id=$1
RETURNING id, club_id, creator_id, name, scoring_rules, start_at, end_at, price_rub, is_closed, game_type, status, created_at, updated_at
`,
		id,
		next.Name,
		next.Description,
		next.StartAt,
		next.EndAt,
		next.PriceRub,
		next.IsClosed,
		int16(next.GameType),
		int16(next.Status),
	)

	var out model.Series
	var gameType int16
	var status int16
	if err := row.Scan(
		&out.ID,
		&out.ClubID,
		&out.CreatorID,
		&out.Name,
		&out.Description,
		&out.StartAt,
		&out.EndAt,
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
	out.GameType = types.GameType(gameType)
	out.Status = types.SeriesStatus(status)
	return &out, nil
}

func (r *Repo) DeleteSeries(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `UPDATE series SET deleted_at=now(), updated_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
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
