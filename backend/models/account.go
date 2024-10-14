package models

import (
	"backend/db"
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	min = 1000
	max = 10000
)

// Account model
type Account struct {
	Model

	Confirmed        uint `json:"is_confirmed"`
	ConfirmationCode uint `json:"confirmation_code"`
}

// GetByID finds account by ID.
// If account with such ID does not exist returns nil
func (a *Account) GetByID(id int) (*Account, error) {
	var account Account
	err := db.GetDB().First(&account, 1).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("could not find account with ID %d: %w", id, err)
	}

	return &account, nil
}

// SaveAccount creates a new account record or updates existed one
func SaveAccount(a *Account) error {
	var err error

	if db.GetDB().NewRecord(a) {
		err = db.GetDB().Create(a).Error
	} else {
		err = db.GetDB().Save(a).Error
	}

	if err != nil {
		return err
	}

	return nil
}

// IsConfirmed checks if account is confirmed
func (a *Account) IsConfirmed() bool {
	return a.Confirmed != 0
}

// GenerateConfirmationCode generates confirmation code for the account
func (a *Account) GenerateConfirmationCode() {
	rand.Seed(time.Now().UnixNano())
	a.ConfirmationCode = uint(rand.Intn(max-min) + min)
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	a.GenerateConfirmationCode()
	return nil
}

// VerifyCode that provided confirmation code is valid
func (a *Account) VerifyCode(code uint) bool {
	return a.ConfirmationCode == code
}
