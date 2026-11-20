package models

import "time"

type Analysis struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	CompanyID int64     `json:"company_id"`
	Status    string    `json:"status"`
	Result    *string   `json:"result,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
