package repository

import (
	"context"
	"fmt"
	"time"

	"chat/chat_server/internal/repository"
	"common/database/client"
	"common/database/transaction"
)

type chatRepository struct {
	db client.Client
}

func NewChatRepository(db client.Client) repository.ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) CreateChat(ctx context.Context, usernames []string) (int64, error) {
	var chatID int64

	txManager := transaction.NewTransactionManager(r.db.DB())

	err := txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		now := time.Now()

		q1 := client.Query{
			Name:     "chat_repository.CreateChat.InsertChat",
			QueryRaw: `INSERT INTO chats (created_at, updated_at) VALUES ($1,$1) RETURNING id`,
		}

		if err := r.db.DB().QueryRowContext(ctx, q1, now).Scan(&chatID); err != nil {
			return fmt.Errorf("insert chat: %w", err)
		}

		for _, u := range usernames {
			q2 := client.Query{
				Name:     "chat_repository.CreateChat.InsertUser",
				QueryRaw: `INSERT INTO chat_users (chat_id, username, created_at) VALUES ($1,$2,$3)`,
			}

			if _, err := r.db.DB().ExecContext(ctx, q2, chatID, u, now); err != nil {
				return fmt.Errorf("insert user %s: %w", u, err)
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (r *chatRepository) DeleteChat(ctx context.Context, chatID int64) error {
	q := client.Query{
		Name:     "chat_repository.DeleteChat",
		QueryRaw: `DELETE FROM chats WHERE id=$1`,
	}

	cmd, err := r.db.DB().ExecContext(ctx, q, chatID)
	if err != nil {
		return fmt.Errorf("delete chat: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("chat not found")
	}
	return nil
}

func (r *chatRepository) SendMessage(ctx context.Context, fromUser, text string, ts time.Time) error {
	q := client.Query{
		Name:     "chat_repository.SendMessage",
		QueryRaw: `INSERT INTO messages (chat_id, from_user, text, timestamp, created_at) VALUES ($1,$2,$3,$4,$5)`,
	}

	_, err := r.db.DB().ExecContext(ctx, q, 1, fromUser, text, ts, time.Now())
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	return nil
}

func (r *chatRepository) GetChatUsers(ctx context.Context, chatID int64) ([]string, error) {
	q := client.Query{
		Name:     "chat_repository.GetChatUsers",
		QueryRaw: `SELECT username FROM chat_users WHERE chat_id=$1 ORDER BY created_at`,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, chatID)
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

func (r *chatRepository) ChatExists(ctx context.Context, chatID int64) (bool, error) {
	q := client.Query{
		Name:     "chat_repository.ChatExists",
		QueryRaw: `SELECT EXISTS(SELECT 1 FROM chats WHERE id=$1)`,
	}

	var exists bool
	if err := r.db.DB().QueryRowContext(ctx, q, chatID).Scan(&exists); err != nil {
		return false, fmt.Errorf("exists: %w", err)
	}
	return exists, nil
}
