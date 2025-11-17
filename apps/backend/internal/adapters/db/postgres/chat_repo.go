package postgres

import (
	"backend/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepo struct {
	pool *pgxpool.Pool
}

func NewChatRepo(pool *pgxpool.Pool) *ChatRepo {
	return &ChatRepo{pool: pool}
}

func (c *ChatRepo) Create(ctx context.Context, chat *domain.Chat) error {
	const q = `
    INSERT INTO app.chats (id, title, user_id, created_at, updated_at)
    VALUES ($1, $2, $3, now(), now())
    RETURNING created_at, updated_at;
    `

	return c.pool.QueryRow(ctx, q, chat.ID, chat.Title, chat.UserID).
		Scan(&chat.CreatedAt, &chat.UpdatedAt)
}

func (c *ChatRepo) GetByID(ctx context.Context, chatID uuid.UUID) (*domain.Chat, error) {
	const q = `
    SELECT id, title, user_id, created_at, updated_at
    FROM app.chats
    WHERE id = $1;
    `

	var chat domain.Chat
	err := c.pool.QueryRow(ctx, q, chatID).Scan(&chat.ID, &chat.Title, &chat.UserID, &chat.CreatedAt, &chat.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (c *ChatRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Chat, error) {
	const q = `
	SELECT * 
	FROM app.chats
	WHERE user_id = $1
	ORDER BY updated_at DESC;
	`

	rows, err := c.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*domain.Chat
	for rows.Next() {
		var chat domain.Chat
		err := rows.Scan(&chat.ID, &chat.Title, &chat.UserID, &chat.CreatedAt, &chat.UpdatedAt, nil)
		if err != nil {
			return nil, err
		}

		chats = append(chats, &chat)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return chats, nil
}

func (c *ChatRepo) UpdateTitle(ctx context.Context, chat *domain.Chat) error {
	const q = `
	UPDATE app.chats
	SET title = $1,
	    updated_at = now()
	WHERE id = $2
	RETURNING title, updated_at;
	`

	return c.pool.QueryRow(ctx, q, chat.Title, chat.ID).Scan(&chat.Title, &chat.UpdatedAt)
}

func (c *ChatRepo) Touch(ctx context.Context, chatID uuid.UUID, t time.Time) error {
	const q = `
	UPDATE app.chats
	SET last_message_at = $1,
    	updated_at       = now()
	WHERE id = $2;
	`
	_, err := c.pool.Exec(ctx, q, t, chatID)
	return err
}

// Delete удаляет чат и все связанные с ним сообщения в одной транзакции.
func (c *ChatRepo) Delete(ctx context.Context, chatID uuid.UUID) error {
	tx, err := c.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	// Если где-то ошибка — откат
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// 1. Удаляем сообщения чата
	const deleteMsgs = `
		DELETE FROM app.messages
		WHERE chat_id = $1;
	`

	if _, err = tx.Exec(ctx, deleteMsgs, chatID); err != nil {
		return fmt.Errorf("delete messages: %w", err)
	}

	const deleteChat = `
		DELETE FROM app.chats
		WHERE id = $1;
	`

	if _, err = tx.Exec(ctx, deleteChat, chatID); err != nil {
		return fmt.Errorf("delete chat: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
