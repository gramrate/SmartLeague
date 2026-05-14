package series

import (
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

func (r *Repo) AddSeriesParticipant(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO series_participants (series_id, profile_id)
VALUES ($1,$2)
ON CONFLICT DO NOTHING
`, seriesID, profileID)
	return err
}

func (r *Repo) RemoveSeriesParticipant(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM series_participants WHERE series_id=$1 AND profile_id=$2`, seriesID, profileID)
	return err
}

func (r *Repo) CountSeriesParticipants(ctx context.Context, seriesID uuid.UUID) (int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM series_participants WHERE series_id=$1`, seriesID).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *Repo) IsSeriesParticipant(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, `
SELECT exists(
  SELECT 1 FROM series_participants WHERE series_id=$1 AND profile_id=$2
)
`, seriesID, profileID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repo) ListSeriesParticipants(ctx context.Context, seriesID uuid.UUID, limit, offset int) ([]*model.User, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM series_participants WHERE series_id=$1`, seriesID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT p.id, p.nickname, p.name, p.show_name, p.description, p.email, p.password_hash, p.club_id, p.club_state, p.role
FROM series_participants sp
JOIN profiles p ON p.id = sp.profile_id
WHERE sp.series_id=$1
ORDER BY sp.created_at ASC
LIMIT $2 OFFSET $3
`, seriesID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.User
	for rows.Next() {
		var p model.User
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
