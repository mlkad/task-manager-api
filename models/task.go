package models

import "time"

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	Priority  string    `json:"priority"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
