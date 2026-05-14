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
