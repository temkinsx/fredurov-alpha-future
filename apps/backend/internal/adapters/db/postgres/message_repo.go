package postgres

import (
	"backend/internal/domain"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepo struct {
	pool *pgxpool.Pool
}

func NewMessageRepo(pool *pgxpool.Pool) *MessageRepo {
	return &MessageRepo{pool: pool}
}

func (m *MessageRepo) Append(ctx context.Context, msg *domain.Message) error {
	const q = `
	INSERT INTO app.messages (id, chat_id, role, content, created_at)
	VALUES ($1, $2, $3, $4, now())
	RETURNING created_at;
	`

	return m.pool.QueryRow(ctx, q, msg.ID, msg.ChatID, msg.Role, msg.Content).Scan(&msg.CreatedAt)
}

func (m *MessageRepo) GetLastN(ctx context.Context, chatID uuid.UUID, n int) ([]*domain.Message, error) {
	const q = `
	SELECT id, chat_id, role, content, created_at
	FROM (
		SELECT id, chat_id, role, content, created_at
		FROM app.messages
		WHERE chat_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	) AS latest
	ORDER BY created_at ASC;
	`

	rows, err := m.pool.Query(ctx, q, chatID, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message

	for rows.Next() {
		var msg domain.Message
		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.Role, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, err
}

func (m *MessageRepo) ListByChat(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*domain.Message, error) {
	const q = `
	SELECT id, chat_id, role, content, created_at 
	FROM app.messages
	WHERE chat_id = $1
	ORDER BY created_at ASC
	LIMIT $2
	OFFSET $3
	`

	rows, err := m.pool.Query(ctx, q, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var msg domain.Message
		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.Role, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, err
}
