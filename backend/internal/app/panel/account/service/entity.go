package service

import "time"

type Account struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
