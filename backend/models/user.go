package models

import (
	"backend/db"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func generatePasswordHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func comparePasswordHash(hashedPassword, givenPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, givenPassword)
	return err == nil
}

// User model
type User struct {
	Model

	Email            string  `gorm:"size:255" json:"email"`
	Password         []byte  `gorm:"size:512" json:"-"`
	Active           uint    `json:"is_active"`
	AccountOwner     uint    `json:"account_owner"`
	AccountID        uint    `gorm:"index:UserAccountID" json:"account_id"`
	Account          Account `json:"-"`
	ConfirmationCode uint    `json:"confirmation_code"`
}

// SetPassword generates hashed user password
func (u *User) SetPassword(password string) error {
	hashed, err := generatePasswordHash([]byte(password))

	if err != nil {
		return err
	}

	u.Password = hashed

	return nil
}

// CheckPassword checks if user password is valid
func (u *User) CheckPassword(password string) bool {
	return comparePasswordHash(u.Password, []byte(password))
}

// GetUserByEmail finds user by Email.
// If user with such ID does not exist returns nil
func GetUserByEmail(email string) (*User, error) {
	var user User
	err := db.GetDB().First(&user, "email = ?", email).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("could not find user with email %s: %w", email, err)
	}

	return &user, nil
}

// GetUserByID finds user by ID.
// If user with such ID does not exist returns nil
func GetUserByID(id uint) (*User, error) {
	var user User
	err := db.GetDB().First(&user, id).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("could not find user with ID %d: %w", id, err)
	}

	return &user, nil
}

// GetAccount returns user account
func (u *User) GetAccount() (*Account, error) {
	var account Account
	assoc := db.GetDB().Model(u).Association("Account").Find(&account)

	if assoc.Error != nil {
		return nil, assoc.Error
	}

	return &account, nil
}

// GetAllUsers loads all users
func GetAllUsers() ([]*User, error) {
	var users []*User

	if err := db.GetDB().Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// SaveUser creates a new user or updates existed one
func SaveUser(u *User) error {
	var err error

	if db.GetDB().NewRecord(u) {
		err = db.GetDB().Create(u).Error
	} else {
		err = db.GetDB().Save(u).Error
	}

	if err != nil {
		return err
	}

	return nil
}

// IsAccountOwner checks if the user is an account owner
func (u *User) IsAccountOwner() bool {
	return u.AccountOwner != 0
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Active != 0
}

// GenerateConfirmationCode generates confirmation code for the user
func (u *User) GenerateConfirmationCode() {
	rand.Seed(time.Now().UnixNano())
	u.ConfirmationCode = uint(rand.Intn(max-min) + min)
}

// VerifyCode that provided confirmation code is valid
func (u *User) VerifyCode(code uint) bool {
	return u.ConfirmationCode == code
}

// GeneratePassword generates random password
func (u *User) GeneratePassword() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder

	for i := 0; i < 8; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	password := b.String()
	u.SetPassword(password)

	return password
}
