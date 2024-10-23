package storage

import "time"

type Account struct {
	ID               uint      `gorm:"AUTO_INCREMENT" gorm:"primary_key" json:"id"`
	Confirmed        uint      `json:"is_confirmed"`
	ConfirmationCode uint      `json:"confirmation_code"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (a *Account) IsConfirmed() bool {
	return a.Confirmed != 0
}

func (a *Account) VerifyCode(code int) bool {
	return int(a.ConfirmationCode) == code
}

type AccountStorage interface {
	FindById(id uint) (*Account, error)
	Save(acc *Account) error
}
