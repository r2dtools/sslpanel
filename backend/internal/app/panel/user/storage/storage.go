package storage

import (
	accountStorage "backend/internal/app/panel/account/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func generatePasswordHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func comparePasswordHash(hashedPassword, givenPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, givenPassword)

	return err == nil
}

type User struct {
	ID                uint                   `gorm:"AUTO_INCREMENT" gorm:"primary_key" json:"id"`
	Email             string                 `gorm:"size:255" json:"email"`
	Password          []byte                 `gorm:"size:512" json:"-"`
	Active            uint                   `json:"is_active"`
	AccountOwner      uint                   `json:"account_owner"`
	AccountID         uint                   `gorm:"index:UserAccountID" json:"account_id"`
	Account           accountStorage.Account `json:"-"`
	ConfirmationToken string                 `json:"confirmation_token"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

func (u *User) IsAccountOwner() bool {
	return u.AccountOwner == 1
}

func (u *User) IsActive() bool {
	return u.Active != 0
}

func (u *User) SetPassword(password string) error {
	hashed, err := generatePasswordHash([]byte(password))

	if err != nil {
		return err
	}

	u.Password = hashed

	return nil
}

func (u *User) CheckPassword(password string) bool {
	return comparePasswordHash(u.Password, []byte(password))
}

type UserStorage interface {
	FindAll() ([]*User, error)
	FindById(id int) (*User, error)
	FindByEmail(email string) (*User, error)
	Save(user *User) error
}
