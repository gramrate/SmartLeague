package user

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (r *Repo) Create(ctx context.Context, u model.User) (*model.User, error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	row := r.db.QueryRowContext(ctx, `
INSERT INTO profiles (id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
RETURNING id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role
`,
		u.ID,
		u.Nickname,
		u.Name,
		u.ShowName,
		u.Description,
		normalizeEmail(u.Email),
		u.PasswordHash,
		u.ClubID,
		int16(u.ClubState),
		int16(u.Role),
	)

	var out model.User
	var role, clubState int16
	if err := row.Scan(&out.ID, &out.Nickname, &out.Name, &out.ShowName, &out.Description, &out.Email, &out.PasswordHash, &out.ClubID, &clubState, &role); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, errorz.EmailAlreadyExist
		}
		return nil, err
	}
	out.Role = types.Role(role)
	out.ClubState = types.ClubState(clubState)
	return &out, nil
}

func (r *Repo) GetById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role FROM profiles WHERE id=$1`, id)

	var out model.User
	var role, clubState int16
	if err := row.Scan(&out.ID, &out.Nickname, &out.Name, &out.ShowName, &out.Description, &out.Email, &out.PasswordHash, &out.ClubID, &clubState, &role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.UserNotFound
		}
		return nil, err
	}
	out.Role = types.Role(role)
	out.ClubState = types.ClubState(clubState)
	return &out, nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role FROM profiles WHERE lower(email)=lower($1)`, normalizeEmail(email))

	var out model.User
	var role, clubState int16
	if err := row.Scan(&out.ID, &out.Nickname, &out.Name, &out.ShowName, &out.Description, &out.Email, &out.PasswordHash, &out.ClubID, &clubState, &role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.UserNotFound
		}
		return nil, err
	}
	out.Role = types.Role(role)
	out.ClubState = types.ClubState(clubState)
	return &out, nil
}

func (r *Repo) GetAllByFilter(
	ctx context.Context,
	limit, offset int,
	role *types.Role,
	query, emailPrefix *string,
) ([]*model.User, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	where := "1=1"
	args := []any{}
	argN := 1

	if role != nil {
		where += " AND role=$" + itoa(argN)
		args = append(args, int16(*role))
		argN++
	}
	if emailPrefix != nil {
		where += " AND lower(email) LIKE lower($" + itoa(argN) + ")"
		args = append(args, normalizeEmail(*emailPrefix)+"%")
		argN++
	}
	if query != nil {
		q := strings.TrimSpace(*query)
		where += " AND (lower(email) LIKE lower($" + itoa(argN) + ") OR lower(name) LIKE lower($" + itoa(argN) + ") OR lower(nickname) LIKE lower($" + itoa(argN) + "))"
		args = append(args, "%"+strings.ToLower(q)+"%")
		argN++
	}

	var total int
	countSQL := "SELECT count(*) FROM profiles WHERE " + where
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listSQL := "SELECT id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role FROM profiles WHERE " + where + " ORDER BY created_at DESC LIMIT $" + itoa(argN) + " OFFSET $" + itoa(argN+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.User
	for rows.Next() {
		var u model.User
		var rrole, clubState int16
		if err := rows.Scan(&u.ID, &u.Nickname, &u.Name, &u.ShowName, &u.Description, &u.Email, &u.PasswordHash, &u.ClubID, &clubState, &rrole); err != nil {
			return nil, 0, err
		}
		u.Role = types.Role(rrole)
		u.ClubState = types.ClubState(clubState)
		out = append(out, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *Repo) GetGamesByProfileID(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]*model.Game, []string, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `
SELECT count(*)
FROM game_participants gp
WHERE gp.profile_id=$1
`, profileID).Scan(&total); err != nil {
		return nil, nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT g.id, g.series_id, s.name, g.name, g.number, g.description, g.host_id, g.status, g.created_at, g.updated_at
FROM game_participants gp
JOIN games g ON g.id = gp.game_id
JOIN series s ON s.id = g.series_id
WHERE gp.profile_id=$1 AND g.deleted_at IS NULL AND s.deleted_at IS NULL
ORDER BY g.created_at DESC
LIMIT $2 OFFSET $3
`, profileID, limit, offset)
	if err != nil {
		return nil, nil, 0, err
	}
	defer rows.Close()

	var out []*model.Game
	var seriesNames []string
	for rows.Next() {
		var g model.Game
		var seriesName string
		var hostID sql.NullString
		var status int16
		var description sql.NullString
		if err := rows.Scan(&g.ID, &g.SeriesID, &seriesName, &g.Name, &g.Number, &description, &hostID, &status, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, nil, 0, err
		}
		g.Description = nullStringToPtr(description)
		g.HostID = nullStringToUUIDPtr(hostID)
		g.Status = types.GameStatus(status)
		out = append(out, &g)
		seriesNames = append(seriesNames, seriesName)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, 0, err
	}

	return out, seriesNames, total, nil
}

func (r *Repo) GetSeriesByProfileID(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]*model.Series, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `
SELECT count(*)
FROM series_participants sp
WHERE sp.profile_id=$1
`, profileID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT s.id, s.club_id, s.creator_id, s.name, s.scoring_rules, s.start_at, s.end_at, s.price_rub, s.is_closed, s.game_type, s.status, s.created_at, s.updated_at
FROM series_participants sp
JOIN series s ON s.id = sp.series_id
WHERE sp.profile_id=$1 AND s.deleted_at IS NULL
ORDER BY s.start_at DESC
LIMIT $2 OFFSET $3
`, profileID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.Series
	for rows.Next() {
		var s model.Series
		var gameType int16
		var status int16
		if err := rows.Scan(&s.ID, &s.ClubID, &s.CreatorID, &s.Name, &s.Description, &s.StartAt, &s.EndAt, &s.PriceRub, &s.IsClosed, &gameType, &status, &s.CreatedAt, &s.UpdatedAt); err != nil {
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

func (r *Repo) Update(ctx context.Context, u model.User) (*model.User, error) {
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
RETURNING id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role
`,
		u.ID,
		u.Nickname,
		u.Name,
		u.ShowName,
		u.Description,
		normalizeEmail(u.Email),
		u.PasswordHash,
		u.ClubID,
		int16(u.ClubState),
		int16(u.Role),
	)

	var out model.User
	var role, clubState int16
	if err := row.Scan(&out.ID, &out.Nickname, &out.Name, &out.ShowName, &out.Description, &out.Email, &out.PasswordHash, &out.ClubID, &clubState, &role); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, errorz.EmailAlreadyExist
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.UserNotFound
		}
		return nil, err
	}
	out.Role = types.Role(role)
	out.ClubState = types.ClubState(clubState)
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
		return errorz.UserNotFound
	}
	return nil
}

func itoa(n int) string {
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	var buf [16]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = digits[n%10]
		n /= 10
	}
	return string(buf[i:])
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
