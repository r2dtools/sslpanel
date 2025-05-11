package storage

import "time"

type Account struct {
	ID                uint      `gorm:"AUTO_INCREMENT" gorm:"primary_key" json:"id"`
	Confirmed         uint      `json:"is_confirmed"`
	ConfirmationToken string    `json:"confirmation_token"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (a *Account) IsConfirmed() bool {
	return a.Confirmed != 0
}

type AccountStorage interface {
	FindById(id uint) (*Account, error)
	Save(acc *Account) error
}
