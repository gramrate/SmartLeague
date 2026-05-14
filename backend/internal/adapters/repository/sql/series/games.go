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

func (r *Repo) CreateGame(ctx context.Context, g model.Game) (*model.Game, error) {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	if g.Name == "" {
		g.Name = fmt.Sprintf("Игра - %d", g.Number)
	}

	row := r.db.QueryRowContext(ctx, `
INSERT INTO games (id, series_id, name, number, description, host_id, status)
VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING id, series_id, name, number, description, host_id, status, created_at, updated_at
`, g.ID, g.SeriesID, g.Name, g.Number, ptrToNullString(g.Description), ptrToNullUUID(g.HostID), int16(g.Status))

	var out model.Game
	var desc sql.NullString
	var hostID sql.NullString
	var status int16
	if err := row.Scan(&out.ID, &out.SeriesID, &out.Name, &out.Number, &desc, &hostID, &status, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	out.HostID = nullStringToUUIDPtr(hostID)
	out.Status = types.GameStatus(status)
	return &out, nil
}

func (r *Repo) GetGameByID(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, series_id, name, number, description, host_id, status, created_at, updated_at
FROM games WHERE id=$1
`, id)

	var out model.Game
	var desc sql.NullString
	var hostID sql.NullString
	var status int16
	if err := row.Scan(&out.ID, &out.SeriesID, &out.Name, &out.Number, &desc, &hostID, &status, &out.CreatedAt, &out.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.GameNotFound
		}
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	out.HostID = nullStringToUUIDPtr(hostID)
	out.Status = types.GameStatus(status)
	return &out, nil
}

func (r *Repo) ListGamesBySeries(ctx context.Context, seriesID uuid.UUID, limit, offset int) ([]*model.Game, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM games WHERE series_id=$1`, seriesID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT id, series_id, name, number, description, host_id, status, created_at, updated_at
FROM games
WHERE series_id=$1
ORDER BY number ASC
LIMIT $2 OFFSET $3
`, seriesID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.Game
	for rows.Next() {
		var g model.Game
		var desc sql.NullString
		var hostID sql.NullString
		var status int16
		if err := rows.Scan(&g.ID, &g.SeriesID, &g.Name, &g.Number, &desc, &hostID, &status, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, 0, err
		}
		g.Description = nullStringToPtr(desc)
		g.HostID = nullStringToUUIDPtr(hostID)
		g.Status = types.GameStatus(status)
		out = append(out, &g)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *Repo) UpdateGame(ctx context.Context, id uuid.UUID, patch model.GameUpdatePatch) (*model.Game, error) {
	current, err := r.GetGameByID(ctx, id)
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
	if patch.HostID != nil {
		next.HostID = patch.HostID
	}
	if patch.Status != nil {
		next.Status = *patch.Status
	}

	row := r.db.QueryRowContext(ctx, `
UPDATE games
SET name=$2,
    description=$3,
    host_id=$4,
    status=$5,
    updated_at=now()
WHERE id=$1
RETURNING id, series_id, name, number, description, host_id, status, created_at, updated_at
`, id, next.Name, ptrToNullString(next.Description), ptrToNullUUID(next.HostID), int16(next.Status))

	var out model.Game
	var desc sql.NullString
	var hostID sql.NullString
	var status int16
	if err := row.Scan(&out.ID, &out.SeriesID, &out.Name, &out.Number, &desc, &hostID, &status, &out.CreatedAt, &out.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.GameNotFound
		}
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	out.HostID = nullStringToUUIDPtr(hostID)
	out.Status = types.GameStatus(status)
	return &out, nil
}

func (r *Repo) ReplaceGameParticipants(ctx context.Context, gameID uuid.UUID, participantIDs []uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM game_participants WHERE game_id=$1`, gameID); err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO game_participants (game_id, profile_id) VALUES ($1,$2)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, pid := range participantIDs {
		if _, err := stmt.ExecContext(ctx, gameID, pid); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repo) UpsertGameResults(ctx context.Context, gameID uuid.UUID, rows []model.GameResultRow) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.PrepareContext(ctx, `
INSERT INTO game_results (game_id, profile_id, place, role, best_move, first_killed, compensation, yellow_cards, removed, extra_points, total_points)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
ON CONFLICT (game_id, profile_id) DO UPDATE SET
  place=excluded.place,
  role=excluded.role,
  best_move=excluded.best_move,
  first_killed=excluded.first_killed,
  compensation=excluded.compensation,
  yellow_cards=excluded.yellow_cards,
  removed=excluded.removed,
  extra_points=excluded.extra_points,
  total_points=excluded.total_points,
  updated_at=now()
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, rrow := range rows {
		if _, err := stmt.ExecContext(
			ctx,
			gameID,
			rrow.ProfileID,
			ptrToNullInt(rrow.Place),
			mafiaRoleToNullString(rrow.Role),
			ptrToNullString(rrow.BestMove),
			rrow.FirstKilled,
			rrow.Compensation,
			rrow.YellowCards,
			rrow.Removed,
			rrow.ExtraPoints,
			rrow.TotalPoints,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repo) ListSeriesLeaderboard(ctx context.Context, seriesID uuid.UUID, limit, offset int) ([]*model.LeaderboardRow, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `
SELECT count(DISTINCT gr.profile_id)
FROM game_results gr
JOIN games g ON g.id=gr.game_id
WHERE g.series_id=$1
`, seriesID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT gr.profile_id, sum(gr.total_points) AS points
FROM game_results gr
JOIN games g ON g.id=gr.game_id
WHERE g.series_id=$1
GROUP BY gr.profile_id
ORDER BY points DESC, gr.profile_id ASC
LIMIT $2 OFFSET $3
`, seriesID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.LeaderboardRow
	for rows.Next() {
		var row model.LeaderboardRow
		if err := rows.Scan(&row.ProfileID, &row.Points); err != nil {
			return nil, 0, err
		}
		out = append(out, &row)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *Repo) ListGameParticipants(ctx context.Context, gameID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT profile_id FROM game_participants WHERE game_id=$1 ORDER BY created_at ASC`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repo) ListGameResults(ctx context.Context, gameID uuid.UUID) ([]model.GameResultRow, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT game_id, profile_id, place, role, best_move, first_killed, compensation, yellow_cards, removed, extra_points, total_points
FROM game_results
WHERE game_id=$1
ORDER BY profile_id ASC
`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.GameResultRow
	for rows.Next() {
		var row model.GameResultRow
		var place sql.NullInt64
		var role sql.NullString
		var bestMove sql.NullString
		if err := rows.Scan(
			&row.GameID,
			&row.ProfileID,
			&place,
			&role,
			&bestMove,
			&row.FirstKilled,
			&row.Compensation,
			&row.YellowCards,
			&row.Removed,
			&row.ExtraPoints,
			&row.TotalPoints,
		); err != nil {
			return nil, err
		}
		if place.Valid {
			p := int(place.Int64)
			row.Place = &p
		}
		if role.Valid {
			rv := types.MafiaRole(role.String)
			row.Role = &rv
		}
		if bestMove.Valid {
			bm := bestMove.String
			row.BestMove = &bm
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func mafiaRoleToNullString(role *types.MafiaRole) sql.NullString {
	if role == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(*role), Valid: true}
}

func (r *Repo) DeleteGame(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM games WHERE id=$1`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errorz.GameNotFound
	}
	return nil
}

func ptrToNullUUID(p *uuid.UUID) any {
	if p == nil || *p == uuid.Nil {
		return nil
	}
	return p.String()
}

func ptrToNullInt(p *int) any {
	if p == nil {
		return nil
	}
	return *p
}
