package postgres

import (
	"alpha_future_fredurov/apps/backend/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
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

// TODO: транзакция (удаление чата из app.chats и удаление сообщений с id этого чата  до 17.11.2025
func (c *ChatRepo) Delete(ctx context.Context, chatID uuid.UUID) error {
	const q = `
	DELETE 
	FROM app.chats
	WHERE id = $1
	`

	_, err := c.pool.Exec(ctx, q, chatID)
	if err != nil {
		return err
	}

	return nil
}
