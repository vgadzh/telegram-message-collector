package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vgadzh/telegram-message-collector/internal/domain/user"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	const query = `
SELECT EXISTS(
    SELECT 1
    FROM users
    WHERE login=$1
)
`
	var exists bool
	err := r.db.QueryRow(ctx, query, login).Scan(&exists)

	return exists, err
}

func (r *UserRepo) Create(ctx context.Context, user *user.User) (int64, error) {
	const query = `
INSERT INTO users(login, password_hash, role) 
VALUES ($1,$2,$3)
RETURNING id
`
	var id int64
	err := r.db.QueryRow(ctx, query, user.Login, user.PasswordHash, user.Role).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}
	return id, nil
}

func (r *UserRepo) GetByLogin(ctx context.Context, login string) (*user.User, error) {
	const query = `
SELECT
	id,
	login,
	password_hash,
	role,
	created_at
FROM users
WHERE login=$1
`
	var u user.User
	err := r.db.QueryRow(ctx, query, login).Scan(
		&u.ID,
		&u.Login,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("user repo - GetByLogin: %w", err)
	}
	return &u, nil
}
