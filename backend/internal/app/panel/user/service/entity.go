package service

import (
	accountService "backend/internal/app/panel/account/service"
	"time"
)

type User struct {
	ID             int                    `json:"id"`
	Email          string                 `json:"email"`
	Active         bool                   `json:"is_active"`
	IsAccountOwner bool                   `json:"is_account_owner"`
	AccountID      int                    `json:"account_id"`
	Account        accountService.Account `json:"account"`
	CreatedAt      time.Time              `json:"created_at"`
}
