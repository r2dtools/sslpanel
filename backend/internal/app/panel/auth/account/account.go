package account

import (
	accountStorage "backend/internal/app/panel/account/storage"
	userStorage "backend/internal/app/panel/user/storage"

	"gorm.io/gorm"
)

type AccountCreater interface {
	Create(email string, password string, confirmationToken string) (*userStorage.User, error)
}

type accountCreater struct {
	db *gorm.DB
}

func (a accountCreater) Create(email string, password string, confirmationToken string) (*userStorage.User, error) {
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
	user.SetPassword(password)
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

func NewAccountCreater(db *gorm.DB) AccountCreater {
	return accountCreater{
		db: db,
	}
}
