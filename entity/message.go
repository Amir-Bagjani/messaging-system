package entity

import "time"

type Message struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}