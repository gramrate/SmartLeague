package series

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/model"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

func (r *Repo) ListSeriesJudges(ctx context.Context, seriesID uuid.UUID) ([]*model.SeriesJudge, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT sj.profile_id, sj.role,
       p.nickname, p.name, p.show_name
FROM series_judges sj
JOIN profiles p ON p.id = sj.profile_id
WHERE sj.series_id = $1
ORDER BY sj.role DESC, sj.created_at ASC
`, seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model.SeriesJudge
	for rows.Next() {
		j := &model.SeriesJudge{SeriesID: seriesID, User: &model.User{}}
		var role int16
		var nickname, name sql.NullString
		if err := rows.Scan(&j.ProfileID, &role, &nickname, &name, &j.User.ShowName); err != nil {
			return nil, err
		}
		j.Role = model.JudgeRole(role)
		j.User.ID = j.ProfileID
		if nickname.Valid {
			j.User.Nickname = nickname.String
		}
		if name.Valid {
			j.User.Name = name.String
		}
		out = append(out, j)
	}
	return out, rows.Err()
}

func (r *Repo) SetSeriesJudge(ctx context.Context, seriesID, profileID uuid.UUID, role model.JudgeRole) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO series_judges (series_id, profile_id, role)
VALUES ($1, $2, $3)
ON CONFLICT (series_id, profile_id) DO UPDATE SET role = EXCLUDED.role
`, seriesID, profileID, int16(role))
	return err
}

func (r *Repo) RemoveSeriesJudge(ctx context.Context, seriesID, profileID uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `
DELETE FROM series_judges WHERE series_id=$1 AND profile_id=$2
`, seriesID, profileID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errorz.InvalidRequest
	}
	return nil
}

func (r *Repo) IsSeriesJudge(ctx context.Context, seriesID, profileID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
SELECT EXISTS(SELECT 1 FROM series_judges WHERE series_id=$1 AND profile_id=$2)
`, seriesID, profileID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (r *Repo) ClearAllSeriesJudges(ctx context.Context, seriesID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM series_judges WHERE series_id=$1`, seriesID)
	return err
}

func (r *Repo) ClearAllGameJudgesForSeries(ctx context.Context, seriesID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE games SET game_judge_id=NULL, game_judge_confirmed=false
WHERE series_id=$1 AND deleted_at IS NULL
`, seriesID)
	return err
}

func (r *Repo) ClearAllGameDraftJudgesForSeries(ctx context.Context, seriesID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE game_drafts SET judge_id=NULL, judge_confirmed=false
WHERE game_id IN (SELECT id FROM games WHERE series_id=$1 AND deleted_at IS NULL)
`, seriesID)
	return err
}
