package user

import (
	"SmartLeague/internal/domain/common/errorz"
	"SmartLeague/internal/domain/model"
	"SmartLeague/internal/domain/types"
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	clubState *types.ClubState,
	clubQuery *string,
	query *string,
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
		where += " AND p.role=$" + itoa(argN)
		args = append(args, int16(*role))
		argN++
	}
	if clubState != nil {
		if *clubState == types.ClubStateLeader {
			where += " AND p.club_state IN ($" + itoa(argN) + ",$" + itoa(argN+1) + ")"
			args = append(args, int16(types.ClubStateLeader), int16(types.ClubStatePresident))
			argN += 2
		} else {
			where += " AND p.club_state=$" + itoa(argN)
			args = append(args, int16(*clubState))
			argN++
		}
	}
	if clubQuery != nil {
		cq := strings.TrimSpace(*clubQuery)
		if cq != "" {
			where += " AND c.name ILIKE $" + itoa(argN)
			args = append(args, "%"+cq+"%")
			argN++
		}
	}
	if query != nil {
		q := strings.TrimSpace(*query)
		if q != "" {
			where += " AND p.nickname ILIKE $" + itoa(argN)
			args = append(args, "%"+q+"%")
			argN++
		}
	}

	var total int
	countSQL := "SELECT count(*) FROM profiles p LEFT JOIN clubs c ON c.id = p.club_id WHERE " + where
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listSQL := "SELECT p.id, p.nickname, p.name, p.show_name, p.description, p.email, p.password_hash, p.club_id, p.club_state, p.role FROM profiles p LEFT JOIN clubs c ON c.id = p.club_id WHERE " + where + " ORDER BY p.created_at DESC LIMIT $" + itoa(argN) + " OFFSET $" + itoa(argN+1)
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

func (r *Repo) GetGamesByProfileID(ctx context.Context, profileID uuid.UUID, limit, offset int, viewerClubID *uuid.UUID) ([]*model.Game, []string, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	visibilityFilter := " AND (s.is_club_only = false OR s.show_to_all = true)"
	countArgs := []any{profileID}
	listArgs := []any{profileID}
	if viewerClubID != nil {
		visibilityFilter = " AND (s.is_club_only = false OR s.show_to_all = true OR s.club_id = $2)"
		countArgs = append(countArgs, *viewerClubID)
		listArgs = append(listArgs, *viewerClubID)
	}
	limitArg := len(listArgs) + 1
	offsetArg := len(listArgs) + 2

	var total int
	if err := r.db.QueryRowContext(ctx, `
SELECT count(*)
FROM game_participants gp
JOIN games g ON g.id = gp.game_id
JOIN series s ON s.id = g.series_id
WHERE gp.profile_id=$1 AND g.deleted_at IS NULL AND s.deleted_at IS NULL`+visibilityFilter,
		countArgs...).Scan(&total); err != nil {
		return nil, nil, 0, err
	}

	listArgs = append(listArgs, limit, offset)
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
SELECT g.id, g.series_id, s.name, g.name, g.number, g.description, g.host_id, g.status, g.created_at, g.updated_at
FROM game_participants gp
JOIN games g ON g.id = gp.game_id
JOIN series s ON s.id = g.series_id
WHERE gp.profile_id=$1 AND g.deleted_at IS NULL AND s.deleted_at IS NULL%s
ORDER BY g.created_at DESC
LIMIT $%d OFFSET $%d
`, visibilityFilter, limitArg, offsetArg), listArgs...)
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

func (r *Repo) GetSeriesByProfileID(
	ctx context.Context,
	profileID uuid.UUID,
	limit, offset int,
	query, from, to *string,
	isRating *bool,
	showPast, showClosed bool,
	viewerClubID *uuid.UUID,
) ([]*model.PlayerSeriesItem, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	// Filters on series columns only (no sp/sj prefix needed for s.* columns)
	filters := "s.deleted_at IS NULL"
	args := []any{profileID}
	argN := 2
	if !showPast {
		filters += " AND s.end_at::date >= CURRENT_DATE"
	}
	if !showClosed {
		filters += " AND s.is_closed = false"
	}
	if isRating != nil {
		filters += " AND s.is_rating=$" + itoa(argN)
		args = append(args, *isRating)
		argN++
	}
	if query != nil && *query != "" {
		filters += " AND s.name ILIKE $" + itoa(argN)
		args = append(args, "%"+strings.TrimSpace(*query)+"%")
		argN++
	}
	if from != nil && *from != "" {
		filters += " AND s.end_at >= $" + itoa(argN) + "::date"
		args = append(args, *from)
		argN++
	}
	if to != nil && *to != "" {
		filters += " AND s.start_at < ($" + itoa(argN) + "::date + interval '1 day')"
		args = append(args, *to)
		argN++
	}
	if viewerClubID != nil {
		filters += " AND (s.is_club_only = false OR s.show_to_all = true OR s.club_id = $" + itoa(argN) + ")"
		args = append(args, *viewerClubID)
		argN++
	} else {
		filters += " AND (s.is_club_only = false OR s.show_to_all = true)"
	}

	// JOIN both participants and judges; the user appears if in either table
	baseFrom := `FROM series s
LEFT JOIN series_participants sp ON sp.series_id = s.id AND sp.profile_id = $1
LEFT JOIN series_judges sj ON sj.series_id = s.id AND sj.profile_id = $1`
	baseWhere := "(sp.profile_id IS NOT NULL OR sj.profile_id IS NOT NULL) AND " + filters

	var total int
	countSQL := "SELECT count(*) " + baseFrom + " WHERE " + baseWhere
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listSQL := `SELECT s.id, s.club_id, s.creator_id, s.name, s.description, s.start_at, s.end_at,
       s.price_rub, s.is_rating, s.is_club_only, s.show_to_all, s.is_closed, s.is_tournament, s.game_type,
       s.created_at, s.updated_at,
       (sj.profile_id IS NOT NULL) AS is_judge, sj.role AS judge_role ` +
		baseFrom + " WHERE " + baseWhere +
		" ORDER BY s.start_at DESC LIMIT $" + itoa(argN) + " OFFSET $" + itoa(argN+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.PlayerSeriesItem
	for rows.Next() {
		var item model.PlayerSeriesItem
		var desc sql.NullString
		var gameType int16
		var judgeRole sql.NullInt16
		if err := rows.Scan(
			&item.ID, &item.ClubID, &item.CreatorID, &item.Name, &desc,
			&item.StartAt, &item.EndAt, &item.PriceRub, &item.IsRating, &item.IsClubOnly, &item.ShowToAll,
			&item.IsClosed, &item.IsTournament, &gameType,
			&item.CreatedAt, &item.UpdatedAt,
			&item.IsJudge, &judgeRole,
		); err != nil {
			return nil, 0, err
		}
		item.Description = desc.String
		item.GameType = types.GameType(gameType)
		if judgeRole.Valid {
			r := model.JudgeRole(judgeRole.Int16)
			item.JudgeRole = &r
		}
		out = append(out, &item)
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
