package postgres

import (
	"backend/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

type UserRepo struct {
	pool *pgxpool.Pool
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, is_active, created_at, last_login_at
		FROM auth.users
		WHERE email = $1;
	`

	var user domain.User
	var lastLoginAt *time.Time

	err := u.pool.QueryRow(ctx, q, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&lastLoginAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user.LastLoginAt = lastLoginAt
	return &user, nil
}

func (u *UserRepo) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	const q = `
		UPDATE auth.users
		SET last_login_at = now()
		WHERE id = $1;
	`

	_, err := u.pool.Exec(ctx, q, userID)
	return err
}
