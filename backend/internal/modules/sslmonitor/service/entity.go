package service

import "time"

type Domain struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Data      string    `json:"data"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type AddDomainRequest struct {
	URL  string `json:"url"`
	Data string `json:"data"`
}
