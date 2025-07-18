package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ChatRepository struct {
	pool *pgxpool.Pool
}

func NewChatRepository(pool *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{pool: pool}
}

func (r *ChatRepository) CreateChat(ctx context.Context, usernames []string) (int64, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var chatID int64
	now := time.Now()
	if err := tx.QueryRow(ctx, `INSERT INTO chats (created_at, updated_at) VALUES ($1,$1) RETURNING id`, now).Scan(&chatID); err != nil {
		return 0, fmt.Errorf("insert chat: %w", err)
	}

	for _, u := range usernames {
		if _, err := tx.Exec(ctx, `INSERT INTO chat_users (chat_id, username, created_at) VALUES ($1,$2,$3)`, chatID, u, now); err != nil {
			return 0, fmt.Errorf("insert user %s: %w", u, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return chatID, nil
}

func (r *ChatRepository) DeleteChat(ctx context.Context, chatID int64) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM chats WHERE id=$1`, chatID)
	if err != nil {
		return fmt.Errorf("delete chat: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("chat not found")
	}
	return nil
}

func (r *ChatRepository) SendMessage(ctx context.Context, fromUser, text string, ts time.Time) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO messages (chat_id, from_user, text, timestamp, created_at) VALUES ($1,$2,$3,$4,$5)`, 1, fromUser, text, ts, time.Now())
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	return nil
}

func (r *ChatRepository) GetChatUsers(ctx context.Context, chatID int64) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT username FROM chat_users WHERE chat_id=$1 ORDER BY created_at`, chatID)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()
	var res []string
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

func (r *ChatRepository) ChatExists(ctx context.Context, chatID int64) (bool, error) {
	var exists bool
	if err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM chats WHERE id=$1)`, chatID).Scan(&exists); err != nil {
		return false, fmt.Errorf("exists: %w", err)
	}
	return exists, nil
}
