package storage

import "time"

type Account struct {
	ID                uint      `gorm:"AUTO_INCREMENT;primary_key" json:"id"`
	ConfirmationToken string    `json:"confirmation_token"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type AccountStorage interface {
	FindById(id uint) (*Account, error)
	Save(acc *Account) error
}
