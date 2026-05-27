package series

import (
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (r *Repo) ListSeriesParticipants(ctx context.Context, seriesID uuid.UUID, limit, offset int, query *string) ([]*model.User, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	where := "sp.series_id=$1"
	countArgs := []any{seriesID}
	listArgs := []any{seriesID}
	if query != nil {
		qv := strings.TrimSpace(*query)
		if qv != "" {
			where += " AND (p.nickname ILIKE $2 OR p.name ILIKE $2 OR p.email ILIKE $2)"
			like := "%" + qv + "%"
			countArgs = append(countArgs, like)
			listArgs = append(listArgs, like)
		}
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM series_participants sp JOIN profiles p ON p.id = sp.profile_id WHERE `+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT p.id, p.nickname, p.name, p.show_name, p.description, p.email, p.password_hash, p.club_id, p.club_state, p.role
FROM series_participants sp
JOIN profiles p ON p.id = sp.profile_id
WHERE `+where+`
ORDER BY sp.created_at ASC
LIMIT $`+fmt.Sprintf("%d", len(listArgs)+1)+` OFFSET $`+fmt.Sprintf("%d", len(listArgs)+2)+`
`, append(listArgs, limit, offset)...)
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

func (r *Repo) ListPaidSeriesParticipants(ctx context.Context, seriesID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT profile_id
FROM series_paid_participants
WHERE series_id=$1
ORDER BY paid_at ASC
`, seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]uuid.UUID, 0)
	for rows.Next() {
		var profileID uuid.UUID
		if err := rows.Scan(&profileID); err != nil {
			return nil, err
		}
		out = append(out, profileID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repo) SetSeriesParticipantPaid(ctx context.Context, seriesID uuid.UUID, profileID uuid.UUID, paid bool) error {
	if paid {
		_, err := r.db.ExecContext(ctx, `
INSERT INTO series_paid_participants (series_id, profile_id)
VALUES ($1,$2)
ON CONFLICT (series_id, profile_id) DO UPDATE SET paid_at=now()
`, seriesID, profileID)
		return err
	}
	_, err := r.db.ExecContext(ctx, `
DELETE FROM series_paid_participants
WHERE series_id=$1 AND profile_id=$2
`, seriesID, profileID)
	return err
}
