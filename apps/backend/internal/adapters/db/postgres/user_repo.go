package postgres

import (
	"backend/internal/domain"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	const q = `
INSERT INTO app.users (id, email, created_at)
VALUES ($1, $2, now())
RETURNING created_at;
`
	return r.pool.QueryRow(ctx, q,
		u.ID,
		u.Email,
	).Scan(&u.CreatedAt)
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const q = `
SELECT id, email, created_at
FROM app.users
WHERE id = $1;
`
	var u domain.User
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&u.ID,
		&u.Email,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
SELECT id, email, created_at
FROM app.users
WHERE email = $1;
`
	var u domain.User
	err := r.pool.QueryRow(ctx, q, email).Scan(
		&u.ID,
		&u.Email,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `
DELETE FROM app.users
WHERE id = $1;
`
	_, err := r.pool.Exec(ctx, q, id)
	return err
}
