package account

import (
	accountStorage "backend/internal/app/panel/account/storage"
	userStorage "backend/internal/app/panel/user/storage"

	"gorm.io/gorm"
)

type AccountCreator interface {
	Create(email string, password string, confirmationToken string) (*userStorage.User, error)
}

type accountCreator struct {
	db *gorm.DB
}

func (a accountCreator) Create(email string, password string, confirmationToken string) (*userStorage.User, error) {
	db := a.db.Begin()

	account := new(accountStorage.Account)
	account.ConfirmationToken = confirmationToken

	err := db.Create(account).Error

	if err != nil {
		db.Rollback()

		return nil, err
	}

	user := new(userStorage.User)
	user.AccountID = account.ID
	user.AccountOwner = 1
	user.Email = email

	err = user.SetPassword(password)

	if err != nil {
		return nil, err
	}

	user.Account = *account

	err = db.Create(user).Error

	if err != nil {
		db.Rollback()

		return nil, err
	}

	err = db.Commit().Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewAccountCreator(db *gorm.DB) AccountCreator {
	return accountCreator{
		db: db,
	}
}
