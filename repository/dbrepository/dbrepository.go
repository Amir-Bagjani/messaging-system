package dbrepository

import (
	"database/sql"
	"messaging-system/entity"
	"time"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (repo *MessageRepository) StoreMessage(userID, message string) error {
	_, err := repo.db.Exec("INSERT INTO messages (user_id, message, timestamp) VALUES ($1, $2, $3)",
		userID, message, time.Now())
	return err
}

func (repo *MessageRepository) GetMessagesByUserID(userID string) ([]entity.Message, error) {
	rows, err := repo.db.Query("SELECT id, user_id, message, timestamp FROM messages WHERE user_id = $1 ORDER BY timestamp DESC",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Message, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
