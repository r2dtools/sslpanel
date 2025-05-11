package service

import (
	"backend/config"
	accountStorage "backend/internal/app/panel/account/storage"
	"backend/internal/app/panel/auth/account"
	userStorage "backend/internal/app/panel/user/storage"
	"backend/internal/pkg/logger"
	"backend/internal/pkg/notification"
	"backend/internal/pkg/token"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidConfirmationToken = errors.New("confirmation token is invalid or expired")
	ErrAccountAlreadyExists     = errors.New("account with such email already exists")
	ErrUserNotFound             = errors.New("user not found")
)

const (
	confirmationTokenExpiration = 10 * time.Minute
	tokenLength                 = 8
)

type AuthService struct {
	config            *config.Config
	userStorage       userStorage.UserStorage
	accountStorage    accountStorage.AccountStorage
	accountCreater    account.AccountCreater
	emailNotification notification.EmailNotificationService
	logger            logger.Logger
}

func (a AuthService) Register(email, password string) error {
	user, err := a.userStorage.FindByEmail(email)

	if err != nil {
		return err
	}

	confirmationToken := token.GenerateRandomToken(tokenLength)

	if user == nil {
		user, err = a.accountCreater.Create(email, password, confirmationToken)

		if err != nil {
			return err
		}
	} else {
		account := &user.Account

		if !user.IsAccountOwner() || account.IsConfirmed() {
			return ErrAccountAlreadyExists
		}

		account.ConfirmationToken = confirmationToken

		if err = a.accountStorage.Save(account); err != nil {
			return err
		}
	}

	linkToken := generateLinkToken(confirmationToken, int(user.ID), time.Now())
	link := fmt.Sprintf("%s/auth/confirm?token=%s", a.config.PanelHost, linkToken)
	data := struct{ Link string }{link}

	err = a.emailNotification.CreateAndSendHtmlNotification(
		"confirmEmail",
		"signup-confirm-email-template",
		email,
		"Email confirmation",
		data,
	)

	if err != nil {
		a.logger.Error("failed to send confirmation notification for email %s: %v", email, err.Error())
	}

	return nil
}

func (a AuthService) ConfirmEmail(token string) error {
	confirmationToken, createdAt, userId, ok := parseLinkToken(token)

	if !ok {
		return ErrInvalidConfirmationToken
	}

	user, err := a.userStorage.FindById(userId)

	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	account := user.Account

	if account.ConfirmationToken != confirmationToken ||
		time.Now().After(createdAt.Add(confirmationTokenExpiration)) {
		return ErrInvalidConfirmationToken
	}

	user.Active = 1

	if err = a.userStorage.Save(user); err != nil {
		return err
	}

	account.Confirmed = 1

	return a.accountStorage.Save(&account)
}

func (a AuthService) RecoverPassword(email string) error {
	user, err := a.userStorage.FindByEmail(email)

	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	confirmationToken := token.GenerateRandomToken(tokenLength)
	user.ConfirmationToken = confirmationToken

	if err = a.userStorage.Save(user); err != nil {
		return err
	}

	linkToken := generateLinkToken(confirmationToken, int(user.ID), time.Now())
	link := fmt.Sprintf("%s/auth/reset?token=%s", a.config.PanelHost, linkToken)
	data := struct{ Link string }{link}

	err = a.emailNotification.CreateAndSendPlainNotification(
		"confirmEmail",
		"recover-confirm-email-template",
		email,
		"Email confirmation",
		data,
	)

	if err != nil {
		a.logger.Error("failed to send email notification for %s on password recover: %v", email, err)
	}

	return nil
}

func (a AuthService) ResetPassword(token string) error {
	confirmationToken, createdAt, userId, ok := parseLinkToken(token)

	if !ok {
		return ErrInvalidConfirmationToken
	}

	user, err := a.userStorage.FindById(userId)

	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	if user.ConfirmationToken != confirmationToken ||
		time.Now().After(createdAt.Add(confirmationTokenExpiration)) {
		return ErrInvalidConfirmationToken
	}

	password := generatePassword()

	if err = user.SetPassword(password); err != nil {
		return err
	}

	if err = a.userStorage.Save(user); err != nil {
		return err
	}

	tplData := struct{ Password string }{password}

	err = a.emailNotification.CreateAndSendPlainNotification(
		"passwordReset",
		"reset-password-email-template",
		user.Email,
		"Password reset",
		tplData,
	)

	if err != nil {
		a.logger.Error("failed to send email notification for user with ID %d: %v", userId, err)
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

func generateLinkToken(confirmationToken string, userId int, now time.Time) string {
	tokenRaw := fmt.Sprintf("%s_%d_%d", confirmationToken, now.Unix(), userId)

	return base64.RawStdEncoding.EncodeToString([]byte(tokenRaw))
}

func parseLinkToken(token string) (confirmationToken string, createdAt time.Time, userId int, ok bool) {
	decoded, err := base64.RawStdEncoding.DecodeString(token)

	if err != nil {
		return
	}

	parts := strings.Split(string(decoded), "_")

	if len(parts) != 3 {
		return
	}

	confirmationToken = parts[0]
	createdAtTimestamp, err := strconv.Atoi(parts[1])

	if err != nil {
		return
	}

	createdAt = time.Unix(int64(createdAtTimestamp), 0)
	userId, err = strconv.Atoi(parts[2])

	if err != nil {
		return
	}

	return confirmationToken, createdAt, userId, true
}

func NewAuthService(
	config *config.Config,
	userStorage userStorage.UserStorage,
	accountStorage accountStorage.AccountStorage,
	accountCreater account.AccountCreater,
	emailNotification notification.EmailNotificationService,
	logger logger.Logger,
) AuthService {
	return AuthService{
		config:            config,
		userStorage:       userStorage,
		accountStorage:    accountStorage,
		accountCreater:    accountCreater,
		emailNotification: emailNotification,
		logger:            logger,
	}
}
