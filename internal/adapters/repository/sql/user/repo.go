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
INSERT INTO profiles (id, nickname, name, show_name, description, email, password_hash, club_id, club_state, role, surname)
VALUES ($1,'',$2,true,NULL,$3,$4,NULL,0,$5,$6)
RETURNING id, email, password_hash, name, surname, role
`,
		u.ID,
		u.Name,
		normalizeEmail(u.Email),
		u.PasswordHash,
		int16(u.Role),
		u.Surname,
	)

	var out model.User
	var role int16
	if err := row.Scan(&out.ID, &out.Email, &out.PasswordHash, &out.Name, &out.Surname, &role); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, errorz.EmailAlreadyExist
		}
		return nil, err
	}
	out.Role = types.Role(role)
	return &out, nil
}

func (r *Repo) GetById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, email, password_hash, name, surname, role FROM profiles WHERE id=$1`, id)

	var out model.User
	var role int16
	if err := row.Scan(&out.ID, &out.Email, &out.PasswordHash, &out.Name, &out.Surname, &role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.UserNotFound
		}
		return nil, err
	}
	out.Role = types.Role(role)
	return &out, nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, email, password_hash, name, surname, role FROM profiles WHERE lower(email)=lower($1)`, normalizeEmail(email))

	var out model.User
	var role int16
	if err := row.Scan(&out.ID, &out.Email, &out.PasswordHash, &out.Name, &out.Surname, &role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorz.UserNotFound
		}
		return nil, err
	}
	out.Role = types.Role(role)
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
		where += " AND (lower(email) LIKE lower($" + itoa(argN) + ") OR lower(name) LIKE lower($" + itoa(argN) + ") OR lower(surname) LIKE lower($" + itoa(argN) + "))"
		args = append(args, "%"+strings.ToLower(q)+"%")
		argN++
	}

	var total int
	countSQL := "SELECT count(*) FROM profiles WHERE " + where
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listSQL := "SELECT id, email, password_hash, name, surname, role FROM profiles WHERE " + where + " ORDER BY created_at DESC LIMIT $" + itoa(argN) + " OFFSET $" + itoa(argN+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*model.User
	for rows.Next() {
		var u model.User
		var rrole int16
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Surname, &rrole); err != nil {
			return nil, 0, err
		}
		u.Role = types.Role(rrole)
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
SET email=$2,
    password_hash=$3,
    name=$4,
    surname=$5,
    role=$6,
    updated_at=now()
WHERE id=$1
RETURNING id, email, password_hash, name, surname, role
`,
		u.ID,
		normalizeEmail(u.Email),
		u.PasswordHash,
		u.Name,
		u.Surname,
		int16(u.Role),
	)

	var out model.User
	var role int16
	if err := row.Scan(&out.ID, &out.Email, &out.PasswordHash, &out.Name, &out.Surname, &role); err != nil {
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
