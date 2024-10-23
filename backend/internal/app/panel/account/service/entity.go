package service

import "time"

type Account struct {
	ID        uint      `json:"id"`
	Confirmed uint      `json:"is_confirmed"`
	CreatedAt time.Time `json:"created_at"`
}
