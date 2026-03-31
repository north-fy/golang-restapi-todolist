package models

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          int          `json:"id"`
	UserID      int          `json:"user_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Completed   bool         `json:"completed"`
	CreatedAt   time.Time    `json:"created_at"`
	CompletedAt sql.NullTime `json:"completed_at"`
}
