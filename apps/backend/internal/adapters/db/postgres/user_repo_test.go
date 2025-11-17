package postgres

import (
	"backend/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestUserRepo_Create_Success(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepo(testPool)

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	u := &domain.User{
		ID:    uuid.New(),
		Email: "user@example.com",
	}

	err = repo.Create(ctx, u)
	require.NoError(t, err)
	require.False(t, u.CreatedAt.IsZero())

	var (
		dbID        uuid.UUID
		dbEmail     string
		dbCreatedAt time.Time
	)

	err = testPool.QueryRow(ctx, `
		SELECT id, email, created_at
		FROM app.users
		WHERE id = $1;
	`, u.ID).Scan(&dbID, &dbEmail, &dbCreatedAt)
	require.NoError(t, err)

	require.Equal(t, u.ID, dbID)
	require.Equal(t, u.Email, dbEmail)
	require.False(t, dbCreatedAt.IsZero())
}

func TestUserRepo_Create_DuplicateID_Error(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepo(testPool)

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	id := uuid.New()

	u1 := &domain.User{
		ID:    id,
		Email: "user1@example.com",
	}
	u2 := &domain.User{
		ID:    id, // тот же ID
		Email: "user2@example.com",
	}

	err = repo.Create(ctx, u1)
	require.NoError(t, err)

	err = repo.Create(ctx, u2)
	require.Error(t, err)
}

func TestUserRepo_Create_DuplicateEmail_Error(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepo(testPool)

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	email := "dup@example.com"

	u1 := &domain.User{
		ID:    uuid.New(),
		Email: email,
	}
	u2 := &domain.User{
		ID:    uuid.New(),
		Email: email, // тот же email → UNIQUE
	}

	err = repo.Create(ctx, u1)
	require.NoError(t, err)

	err = repo.Create(ctx, u2)
	require.Error(t, err)
}

func TestUserRepo_GetByID_Success(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepo(testPool)

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	id := uuid.New()
	email := "getbyid@example.com"

	_, err = testPool.Exec(ctx, `
		INSERT INTO app.users (id, email, created_at)
		VALUES ($1, $2, now());
	`, id, email)
	require.NoError(t, err)

	u, err := repo.GetByID(ctx, id)
	require.NoError(t, err)

	require.Equal(t, id, u.ID)
	require.Equal(t, email, u.Email)
	require.False(t, u.CreatedAt.IsZero())
}

func TestUserRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepo(testPool)

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestUserRepo_GetByEmail_Success(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepo(testPool)

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	id := uuid.New()
	email := "byemail@example.com"

	_, err = testPool.Exec(ctx, `
		INSERT INTO app.users (id, email, created_at)
		VALUES ($1, $2, now());
	`, id, email)
	require.NoError(t, err)

	u, err := repo.GetByEmail(ctx, email)
	require.NoError(t, err)

	require.Equal(t, id, u.ID)
	require.Equal(t, email, u.Email)
	require.False(t, u.CreatedAt.IsZero())
}

func TestUserRepo_GetByEmail_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepo(testPool)

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	_, err = repo.GetByEmail(ctx, "nope@example.com")
	require.Error(t, err)
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestUserRepo_Delete_Success(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepo(testPool)

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	id := uuid.New()
	email := "delete@example.com"

	_, err = testPool.Exec(ctx, `
		INSERT INTO app.users (id, email, created_at)
		VALUES ($1, $2, now());
	`, id, email)
	require.NoError(t, err)

	err = repo.Delete(ctx, id)
	require.NoError(t, err)

	var dbID uuid.UUID
	err = testPool.QueryRow(ctx, `
		SELECT id FROM app.users WHERE id = $1;
	`, id).Scan(&dbID)

	require.Error(t, err)
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestUserRepo_Delete_NoErrorOnMissing(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepo(testPool)

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	// удаляем несуществующего — по текущему контракту это OK
	err = repo.Delete(ctx, uuid.New())
	require.NoError(t, err)
}

func TestNewUserRepo(t *testing.T) {
	repo := NewUserRepo(testPool)
	require.NotNil(t, repo)
	require.Equal(t, testPool, repo.pool)
}
