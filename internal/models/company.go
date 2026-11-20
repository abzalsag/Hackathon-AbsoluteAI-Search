package models

import "time"

type Company struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Website     string    `json:"website"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
