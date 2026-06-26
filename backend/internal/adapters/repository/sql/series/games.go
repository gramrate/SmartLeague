package series

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
)

func (r *Repo) CreateGame(ctx context.Context, g model.Game) (*model.Game, error) {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}

	var newID uuid.UUID
	if err := r.db.QueryRowContext(ctx, `
INSERT INTO games (id, series_id, name, number, description, host_id, status)
VALUES ($1,$2,$3,(SELECT COALESCE(MAX(number), 0) + 1 FROM games WHERE series_id=$2),$4,$5,$6)
RETURNING id
`, g.ID, g.SeriesID, g.Name, ptrToNullString(g.Description), ptrToNullUUID(g.HostID), int16(g.Status),
	).Scan(&newID); err != nil {
		return nil, err
	}
	return r.GetGameByID(ctx, newID)
}

func (r *Repo) GetGameByID(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT g.id, g.series_id, g.name, g.number, g.description, g.host_id, g.status,
       s.is_tournament,
       g.game_judge_id, g.game_judge_confirmed,
       COALESCE(p.nickname,''), COALESCE(p.name,''), COALESCE(p.show_name, false),
       g.created_at, g.updated_at
FROM games g
JOIN series s ON s.id = g.series_id
LEFT JOIN profiles p ON p.id = g.game_judge_id
WHERE g.id=$1 AND g.deleted_at IS NULL
`, id)

	var out model.Game
	var desc sql.NullString
	var hostID sql.NullString
	var judgeID sql.NullString
	var status int16
	if err := row.Scan(
		&out.ID, &out.SeriesID, &out.Name, &out.Number, &desc, &hostID, &status,
		&out.IsTournament,
		&judgeID, &out.JudgeConfirmed,
		&out.JudgeNickname, &out.JudgeName, &out.JudgeShowName,
		&out.CreatedAt, &out.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.GameNotFound
		}
		return nil, err
	}
	out.Description = nullStringToPtr(desc)
	out.HostID = nullStringToUUIDPtr(hostID)
	out.JudgeID = nullStringToUUIDPtr(judgeID)
	out.Status = types.GameStatus(status)
	return &out, nil
}

func (r *Repo) SetGameJudge(ctx context.Context, gameID uuid.UUID, judgeID *uuid.UUID, confirmed bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE games SET game_judge_id=$2, game_judge_confirmed=$3, updated_at=now() WHERE id=$1`,
		gameID, ptrToNullUUID(judgeID), confirmed,
	)
	return err
}

func (r *Repo) UpsertGameDraft(ctx context.Context, d *model.GameDraft) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO game_drafts (game_id, rows, judge_id, judge_confirmed, updated_at)
VALUES ($1, $2, $3, $4, now())
ON CONFLICT (game_id) DO UPDATE
SET rows=$2, judge_id=$3, judge_confirmed=$4, updated_at=now()
`, d.GameID, d.Rows, ptrToNullUUID(d.JudgeID), d.JudgeConfirmed)
	return err
}

func (r *Repo) GetGameDraft(ctx context.Context, gameID uuid.UUID) (*model.GameDraft, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT game_id, rows, judge_id, judge_confirmed, updated_at FROM game_drafts WHERE game_id=$1`,
		gameID,
	)
	var d model.GameDraft
	var judgeID sql.NullString
	if err := row.Scan(&d.GameID, &d.Rows, &judgeID, &d.JudgeConfirmed, &d.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	d.JudgeID = nullStringToUUIDPtr(judgeID)
	return &d, nil
}

func (r *Repo) ListGamesBySeries(ctx context.Context, seriesID uuid.UUID, limit, offset int, includeDrafts bool) ([]*model.Game, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	countQuery := `SELECT count(*) FROM games WHERE series_id=$1 AND deleted_at IS NULL`
	if !includeDrafts {
		countQuery += ` AND status <> 0`
	}
	if err := r.db.QueryRowContext(ctx, countQuery, seriesID).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := `
SELECT g.id, g.series_id, g.name, g.number, g.description, g.host_id, g.status,
       s.is_tournament,
       g.game_judge_id, g.game_judge_confirmed,
       COALESCE(p.nickname,''), COALESCE(p.name,''), COALESCE(p.show_name, false),
       g.created_at, g.updated_at
FROM games g
JOIN series s ON s.id = g.series_id
LEFT JOIN profiles p ON p.id = g.game_judge_id
WHERE g.series_id=$1 AND g.deleted_at IS NULL`
	if !includeDrafts {
		listQuery += `
  AND g.status <> 0`
	}
	listQuery += `
ORDER BY g.number ASC
LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, listQuery, seriesID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.Game
	for rows.Next() {
		var g model.Game
		var desc sql.NullString
		var hostID sql.NullString
		var judgeID sql.NullString
		var status int16
		if err := rows.Scan(
			&g.ID, &g.SeriesID, &g.Name, &g.Number, &desc, &hostID, &status,
			&g.IsTournament,
			&judgeID, &g.JudgeConfirmed,
			&g.JudgeNickname, &g.JudgeName, &g.JudgeShowName,
			&g.CreatedAt, &g.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		g.Description = nullStringToPtr(desc)
		g.HostID = nullStringToUUIDPtr(hostID)
		g.JudgeID = nullStringToUUIDPtr(judgeID)
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

// UpsertGameResults writes results for registered profiles only (legacy endpoint).
// SaveDraft/Publish use ClearGameResults + UpsertAllGameResults instead.
func (r *Repo) UpsertGameResults(ctx context.Context, gameID uuid.UUID, rows []model.GameResultRow) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	profileStmt, err := tx.PrepareContext(ctx, `
INSERT INTO game_results (game_id, profile_id, place, role, best_move, first_killed, compensation, yellow_cards, removed, victory_points, extra_points, total_points)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
ON CONFLICT (game_id, profile_id) WHERE profile_id IS NOT NULL DO UPDATE SET
  place=excluded.place,
  role=excluded.role,
  best_move=excluded.best_move,
  first_killed=excluded.first_killed,
  compensation=excluded.compensation,
  yellow_cards=excluded.yellow_cards,
  removed=excluded.removed,
  victory_points=excluded.victory_points,
  extra_points=excluded.extra_points,
  total_points=excluded.total_points,
  updated_at=now()
`)
	if err != nil {
		return err
	}
	defer profileStmt.Close()

	guestStmt, err := tx.PrepareContext(ctx, `
INSERT INTO game_results (game_id, guest_id, place, role, best_move, first_killed, compensation, yellow_cards, removed, victory_points, extra_points, total_points)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
ON CONFLICT (game_id, guest_id) WHERE guest_id IS NOT NULL DO UPDATE SET
  place=excluded.place,
  role=excluded.role,
  best_move=excluded.best_move,
  first_killed=excluded.first_killed,
  compensation=excluded.compensation,
  yellow_cards=excluded.yellow_cards,
  removed=excluded.removed,
  victory_points=excluded.victory_points,
  extra_points=excluded.extra_points,
  total_points=excluded.total_points,
  updated_at=now()
`)
	if err != nil {
		return err
	}
	defer guestStmt.Close()

	for _, rrow := range rows {
		if rrow.ProfileID != nil {
			if _, err := profileStmt.ExecContext(
				ctx,
				gameID, *rrow.ProfileID,
				ptrToNullInt(rrow.Place),
				mafiaRoleToNullString(rrow.Role),
				ptrToNullString(rrow.BestMove),
				rrow.FirstKilled,
				rrow.Compensation, rrow.YellowCards, rrow.Removed,
				rrow.VictoryPoints, rrow.ExtraPoints, rrow.TotalPoints,
			); err != nil {
				return err
			}
		} else if rrow.GuestID != nil {
			if _, err := guestStmt.ExecContext(
				ctx,
				gameID, *rrow.GuestID,
				ptrToNullInt(rrow.Place),
				mafiaRoleToNullString(rrow.Role),
				ptrToNullString(rrow.BestMove),
				rrow.FirstKilled,
				rrow.Compensation, rrow.YellowCards, rrow.Removed,
				rrow.VictoryPoints, rrow.ExtraPoints, rrow.TotalPoints,
			); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *Repo) ClearGameResults(ctx context.Context, gameID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM game_results WHERE game_id=$1`, gameID)
	return err
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
SELECT count(*) FROM (
  SELECT COALESCE(gr.profile_id::text, gr.guest_id::text) AS player_key
  FROM game_results gr
  JOIN games g ON g.id = gr.game_id
  WHERE g.series_id = $1 AND g.deleted_at IS NULL AND g.status = 2
  GROUP BY player_key
) sub
`, seriesID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT
  gr.profile_id,
  gr.guest_id,
  sg.nickname AS guest_nickname,
  sum(gr.total_points) AS points
FROM game_results gr
JOIN games g ON g.id = gr.game_id
LEFT JOIN series_guests sg ON sg.id = gr.guest_id
WHERE g.series_id = $1 AND g.deleted_at IS NULL AND g.status = 2
GROUP BY gr.profile_id, gr.guest_id, sg.nickname
ORDER BY points DESC
LIMIT $2 OFFSET $3
`, seriesID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.LeaderboardRow
	for rows.Next() {
		var row model.LeaderboardRow
		var profileIDRaw sql.NullString
		var guestIDRaw sql.NullString
		var guestNickname sql.NullString
		if err := rows.Scan(&profileIDRaw, &guestIDRaw, &guestNickname, &row.Points); err != nil {
			return nil, 0, err
		}
		row.ProfileID = nullStringToUUIDPtr(profileIDRaw)
		row.GuestID = nullStringToUUIDPtr(guestIDRaw)
		row.GuestNickname = nullStringToPtr(guestNickname)
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
SELECT
  gr.game_id,
  gr.profile_id,
  gr.guest_id,
  sg.nickname AS guest_nickname,
  gr.place,
  gr.role,
  gr.best_move,
  gr.first_killed,
  gr.compensation,
  gr.yellow_cards,
  gr.removed,
  gr.victory_points,
  gr.extra_points,
  gr.total_points
FROM game_results gr
LEFT JOIN series_guests sg ON sg.id = gr.guest_id
WHERE gr.game_id = $1
ORDER BY gr.place ASC NULLS LAST
`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.GameResultRow
	for rows.Next() {
		var row model.GameResultRow
		var profileIDRaw sql.NullString
		var guestIDRaw sql.NullString
		var guestNickname sql.NullString
		var place sql.NullInt64
		var role sql.NullString
		var bestMove sql.NullString
		if err := rows.Scan(
			&row.GameID,
			&profileIDRaw,
			&guestIDRaw,
			&guestNickname,
			&place,
			&role,
			&bestMove,
			&row.FirstKilled,
			&row.Compensation,
			&row.YellowCards,
			&row.Removed,
			&row.VictoryPoints,
			&row.ExtraPoints,
			&row.TotalPoints,
		); err != nil {
			return nil, err
		}
		row.ProfileID = nullStringToUUIDPtr(profileIDRaw)
		row.GuestID = nullStringToUUIDPtr(guestIDRaw)
		row.GuestNickname = nullStringToPtr(guestNickname)
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

// FindOrCreateGuest returns an existing guest (matched by series + case-insensitive nickname)
// or creates a new one.
func (r *Repo) FindOrCreateGuest(ctx context.Context, seriesID uuid.UUID, nickname string) (*model.Guest, error) {
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return nil, errorz.InvalidRequest
	}

	var g model.Guest
	err := r.db.QueryRowContext(ctx, `
SELECT id, series_id, nickname, created_at
FROM series_guests
WHERE series_id = $1 AND lower(nickname) = lower($2)
`, seriesID, nickname).Scan(&g.ID, &g.SeriesID, &g.Nickname, &g.CreatedAt)
	if err == nil {
		return &g, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Create new guest
	g.ID = uuid.New()
	g.SeriesID = seriesID
	g.Nickname = nickname
	if err := r.db.QueryRowContext(ctx, `
INSERT INTO series_guests (id, series_id, nickname)
VALUES ($1, $2, $3)
ON CONFLICT (series_id, lower(nickname)) DO UPDATE SET nickname = series_guests.nickname
RETURNING id, series_id, nickname, created_at
`, g.ID, seriesID, nickname).Scan(&g.ID, &g.SeriesID, &g.Nickname, &g.CreatedAt); err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *Repo) DeleteGame(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `UPDATE games SET deleted_at=now(), updated_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
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

func mafiaRoleToNullString(role *types.MafiaRole) sql.NullString {
	if role == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(*role), Valid: true}
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
