package service

import (
	accountStorage "backend/internal/app/panel/account/storage"
	userService "backend/internal/app/panel/user/service"
	userStorage "backend/internal/app/panel/user/storage"
	"backend/internal/pkg/logger"
	"backend/internal/pkg/notification"
	"errors"
	"math/rand"
	"strings"
)

var (
	ErrInvalidConfirmationCode = errors.New("invalid confirmation code")
	ErrAccountAlreadyExists    = errors.New("account with such email already exists")
	ErrUserNotFound            = errors.New("user not found")
)

type AuthService struct {
	userStorage       userStorage.UserStorage
	accountStorage    accountStorage.AccountStorage
	emailNotification notification.EmailNotificationService
	logger            logger.Logger
}

func (a AuthService) Register(email, password string) error {
	user, err := a.userStorage.FindByEmail(email)

	if err != nil {
		return err
	}

	var account *accountStorage.Account

	if user == nil {
		account = new(accountStorage.Account)

		if err := a.accountStorage.Save(account); err != nil {
			return err
		}

		user = new(userStorage.User)
		user.AccountID = account.ID
		user.AccountOwner = 1
		user.Email = email
		user.SetPassword(password)
		user.Account = *account

		if err = a.userStorage.Save(user); err != nil {
			return err
		}
	} else {
		account = &user.Account

		if !user.IsAccountOwner() || account.IsConfirmed() {
			return ErrAccountAlreadyExists
		}

		account.ConfirmationCode = uint(generateConfirmationCode())

		if err = a.accountStorage.Save(account); err != nil {
			return err
		}
	}

	data := struct{ Code uint }{account.ConfirmationCode}

	err = a.emailNotification.CreateAndSendPlainNotification(
		"confirmEmail",
		"signup-confirm-email-template",
		user.Email,
		"Email confirmation",
		data,
	)

	if err != nil {
		a.logger.Error("failed to send confirmation notification for email %s: %v", email, err.Error())
	}

	return nil
}

func (a AuthService) ConfirmEmail(userID, code int) error {
	user, err := a.userStorage.FindById(userID)

	if err != nil {
		return err
	}

	if err == nil {
		return ErrUserNotFound
	}

	account := user.Account

	if !account.VerifyCode(code) {
		return ErrInvalidConfirmationCode
	}

	user.Active = 1

	if err = a.userStorage.Save(user); err != nil {
		return err
	}

	account.Confirmed = 1

	return a.accountStorage.Save(&account)
}

func (a AuthService) RecoverPassword(email string) (userService.User, error) {
	user, err := a.userStorage.FindByEmail(email)

	if err != nil {
		return userService.User{}, err
	}

	if user == nil {
		return userService.User{}, ErrUserNotFound
	}

	user.ConfirmationCode = uint(generateConfirmationCode())

	if err = a.userStorage.Save(user); err != nil {
		return userService.User{}, err
	}

	data := struct{ Code uint }{user.ConfirmationCode}

	err = a.emailNotification.CreateAndSendPlainNotification(
		"confirmEmail",
		"recover-confirm-email-template",
		user.Email,
		"Email confirmation",
		data,
	)

	if err != nil {
		a.logger.Error("failed to send email notification for %s on password recover: %v", email, err)
	}

	return *userService.CreateUser(user), nil
}

func (a AuthService) ResetPassword(userID, code int) error {
	userModel, err := a.userStorage.FindById(userID)

	if err != nil {
		return err
	}

	if userModel == nil {
		return ErrUserNotFound
	}

	if !userModel.VerifyCode(code) {
		return ErrInvalidConfirmationCode
	}

	password := generatePassword()

	if err = userModel.SetPassword(password); err != nil {
		return err
	}

	if err = a.userStorage.Save(userModel); err != nil {
		return err
	}

	tplData := struct{ Password string }{password}

	err = a.emailNotification.CreateAndSendPlainNotification(
		"passwordReset",
		"reset-password-email-template",
		userModel.Email,
		"Password reset",
		tplData,
	)

	if err != nil {
		a.logger.Error("failed to send email notification for user with ID %d: %v", userID, err)
	}

	return nil
}

func generatePassword() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder

	for i := 0; i < 8; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}

func generateConfirmationCode() int {
	return rand.Intn(10000-1000) + 1000
}

func NewAuthService(
	userStorage userStorage.UserStorage,
	emailNotification notification.EmailNotificationService,
) AuthService {
	return AuthService{
		userStorage:       userStorage,
		emailNotification: emailNotification,
	}
}
