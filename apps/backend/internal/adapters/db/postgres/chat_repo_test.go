package postgres

import (
	"alpha_future_fredurov/apps/backend/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestChatRepo_Create_Success(t *testing.T) {
	ctx := context.Background()
	repo := &ChatRepo{pool: testPool}

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	id := uuid.New()
	userID := insertTestUser(t, ctx)

	chat := &domain.Chat{
		ID:            id,
		Title:         "first_chat",
		UserID:        userID,
		LastMessageAt: time.Now(),
	}

	err = repo.Create(ctx, chat)
	require.NoError(t, err)

	var (
		dbID        uuid.UUID
		dbTitle     string
		dbUserID    uuid.UUID
		dbCreatedAt time.Time
		dbUpdatedAt time.Time
	)

	err = testPool.QueryRow(ctx,
		`SELECT id, title, user_id, created_at, updated_at
         FROM app.chats WHERE id = $1`,
		chat.ID,
	).Scan(&dbID, &dbTitle, &dbUserID, &dbCreatedAt, &dbUpdatedAt)
	require.NoError(t, err)

	require.Equal(t, chat.ID, dbID)
	require.Equal(t, chat.Title, dbTitle)
	require.Equal(t, chat.UserID, dbUserID)
	require.False(t, dbCreatedAt.IsZero())
	require.False(t, dbUpdatedAt.IsZero())
}

func TestChatRepo_Create_DuplicateID_Error(t *testing.T) {
	ctx := context.Background()
	repo := &ChatRepo{pool: testPool}

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	id := uuid.New()
	userID := insertTestUser(t, ctx)

	chat1 := &domain.Chat{
		ID:     id,
		Title:  "first_chat",
		UserID: userID,
	}
	chat2 := &domain.Chat{
		ID:     id,
		Title:  "second_chat",
		UserID: userID,
	}

	err = repo.Create(ctx, chat1)
	require.NoError(t, err)

	err = repo.Create(ctx, chat2)
	require.Error(t, err)
}

func TestChatRepo_Delete_Success(t *testing.T) {
	ctx := context.Background()
	repo := &ChatRepo{pool: testPool}

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	id := uuid.New()
	userID := insertTestUser(t, ctx)

	_, err = testPool.Exec(ctx, `
									INSERT INTO app.chats (id, title, user_id, created_at, updated_at) 
									VALUES ($1, $2, $3, now(), now())
									`, id, "to_delete", userID)
	require.NoError(t, err)

	require.NoError(t, repo.Delete(ctx, id))
	var dbID uuid.UUID
	err = testPool.QueryRow(ctx,
		`SELECT id FROM app.chats WHERE id = $1`,
		id,
	).Scan(&dbID)

	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestChatRepo_GetByID_Success(t *testing.T) {
	ctx := context.Background()
	repo := &ChatRepo{pool: testPool}
	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	id := uuid.New()
	userID := insertTestUser(t, ctx)

	_, err = testPool.Exec(ctx, `
									INSERT INTO app.chats (id, title, user_id, created_at, updated_at) 
									VALUES ($1, $2, $3, now(), now())
									`, id, "test_chat", userID)
	require.NoError(t, err)

	chat, err := repo.GetByID(ctx, id)
	require.NoError(t, err)

	var (
		dbID        uuid.UUID
		dbTitle     string
		dbUserID    uuid.UUID
		dbCreatedAt time.Time
		dbUpdatedAt time.Time
	)

	err = testPool.QueryRow(ctx,
		`SELECT id, title, user_id, created_at, updated_at
         FROM app.chats WHERE id = $1`,
		chat.ID,
	).Scan(&dbID, &dbTitle, &dbUserID, &dbCreatedAt, &dbUpdatedAt)

	require.Equal(t, chat.ID, dbID)
	require.Equal(t, "test_chat", dbTitle)
	require.Equal(t, chat.UserID, dbUserID)
	require.False(t, dbCreatedAt.IsZero())
	require.False(t, dbUpdatedAt.IsZero())
}

func TestChatRepo_GetByID_Error(t *testing.T) {
	ctx := context.Background()
	repo := &ChatRepo{pool: testPool}

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	id := uuid.New()

	_, err = repo.GetByID(ctx, id)
	require.Error(t, err)
}

func TestChatRepo_ListByUser_Success(t *testing.T) {
	ctx := context.Background()
	repo := &ChatRepo{pool: testPool}

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	userID := insertTestUser(t, ctx)
	otherUserID := insertTestUser(t, ctx)

	chat1 := &domain.Chat{
		ID:     uuid.New(),
		Title:  "chat_1",
		UserID: userID,
	}
	chat2 := &domain.Chat{
		ID:     uuid.New(),
		Title:  "chat_2",
		UserID: userID,
	}
	chatOther := &domain.Chat{
		ID:     uuid.New(),
		Title:  "foreign_chat",
		UserID: otherUserID,
	}

	_, err = testPool.Exec(ctx, `
		INSERT INTO app.chats (id, title, user_id, created_at, updated_at)
		VALUES
		    ($1, $2, $3, now(), now()),
		    ($4, $5, $6, now(), now()),
		    ($7, $8, $9, now(), now());
	`,
		chat1.ID, chat1.Title, chat1.UserID,
		chat2.ID, chat2.Title, chat2.UserID,
		chatOther.ID, chatOther.Title, chatOther.UserID,
	)
	require.NoError(t, err)

	chats, err := repo.ListByUser(ctx, userID)
	require.NoError(t, err)

	require.Len(t, chats, 2)

	gotIDs := map[uuid.UUID]domain.Chat{}
	for _, ch := range chats {
		gotIDs[ch.ID] = *ch
	}

	require.Contains(t, gotIDs, chat1.ID)
	require.Contains(t, gotIDs, chat2.ID)

	require.Equal(t, chat1.Title, gotIDs[chat1.ID].Title)
	require.Equal(t, chat2.Title, gotIDs[chat2.ID].Title)
	require.Equal(t, userID, gotIDs[chat1.ID].UserID)
	require.Equal(t, userID, gotIDs[chat2.ID].UserID)

	require.False(t, gotIDs[chat1.ID].CreatedAt.IsZero())
	require.False(t, gotIDs[chat2.ID].CreatedAt.IsZero())
}

func TestChatRepo_ListByUser_Empty(t *testing.T) {
	ctx := context.Background()
	repo := &ChatRepo{pool: testPool}

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	userID := insertTestUser(t, ctx)

	chats, err := repo.ListByUser(ctx, userID)
	require.NoError(t, err)
	require.Len(t, chats, 0)
}

func TestChatRepo_Touch_Success(t *testing.T) {
	ctx := context.Background()
	repo := &ChatRepo{pool: testPool}

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	chatID := uuid.New()
	userID := insertTestUser(t, ctx)

	_, err = testPool.Exec(ctx, `
		INSERT INTO app.chats (id, title, user_id, created_at, updated_at, last_message_at)
		VALUES ($1, $2, $3, now(), now(), now())
	`, chatID, "touch_me", userID)
	require.NoError(t, err)

	touchTime := time.Date(2025, 1, 2, 3, 4, 5, 0, time.Local)

	err = repo.Touch(ctx, chatID, touchTime)
	require.NoError(t, err)

	var (
		dbLastMsgAt time.Time
		dbUpdatedAt time.Time
	)

	err = testPool.QueryRow(ctx, `
		SELECT last_message_at, updated_at
		FROM app.chats
		WHERE id = $1
	`, chatID).Scan(&dbLastMsgAt, &dbUpdatedAt)
	require.NoError(t, err)

	require.Equal(t, touchTime, dbLastMsgAt)
	require.False(t, dbUpdatedAt.IsZero())
}

func TestChatRepo_UpdateTitle_Success(t *testing.T) {
	ctx := context.Background()
	repo := &ChatRepo{pool: testPool}

	_, err := testPool.Exec(ctx, "TRUNCATE app.chats CASCADE")
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "TRUNCATE app.users CASCADE")
	require.NoError(t, err)

	chatID := uuid.New()
	userID := insertTestUser(t, ctx)

	_, err = testPool.Exec(ctx, `
		INSERT INTO app.chats (id, title, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
	`, chatID, "old_title", userID)
	require.NoError(t, err)

	chat := &domain.Chat{
		ID:     chatID,
		Title:  "new_title",
		UserID: userID,
	}

	err = repo.UpdateTitle(ctx, chat)
	require.NoError(t, err)

	var (
		dbTitle     string
		dbUserID    uuid.UUID
		dbCreatedAt time.Time
		dbUpdatedAt time.Time
	)

	err = testPool.QueryRow(ctx, `
		SELECT title, user_id, created_at, updated_at
		FROM app.chats
		WHERE id = $1;
	`, chatID).Scan(&dbTitle, &dbUserID, &dbCreatedAt, &dbUpdatedAt)
	require.NoError(t, err)

	require.Equal(t, "new_title", dbTitle)
	require.Equal(t, userID, dbUserID)
	require.False(t, dbCreatedAt.IsZero())
	require.False(t, dbUpdatedAt.IsZero())

	if !chat.UpdatedAt.IsZero() {
		require.WithinDuration(t, dbUpdatedAt, chat.UpdatedAt, time.Second)
	}
}

func TestNewChatRepo(t *testing.T) {
	repo := NewChatRepo(testPool)
	require.NotNil(t, repo)
	require.Equal(t, testPool, repo.pool)
}

func insertTestUser(t *testing.T, ctx context.Context) uuid.UUID {
	t.Helper()
	id := uuid.New()
	_, err := testPool.Exec(ctx,
		`INSERT INTO app.users (id, email, created_at)
         VALUES ($1, $2, now())`,
		id, "user_"+id.String()+"@example.com",
	)
	require.NoError(t, err)
	return id
}
