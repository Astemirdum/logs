package models

import "time"

type Log struct {
	ID        int       `db:"id"`
	Raw       string    `db:"raw"`
	CreatedAt time.Time `db:"created_at"`
}

type WriteLogResponse struct {
	ID int64 `json:"id"`
}

type ReadLogResponse struct {
	ID        int       `json:"id"`
	Raw       string    `json:"raw"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateLogRequest struct {
	Raw string `json:"raw" validate:"required"`
}
