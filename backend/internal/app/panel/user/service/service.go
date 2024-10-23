package service

import (
	accountService "backend/internal/app/panel/account/service"
	"backend/internal/app/panel/user/storage"
	"errors"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid current password")
)

type UserService struct {
	userStorage storage.UserStorage
}

func (s UserService) FindUserByID(id int) (*User, error) {
	userModel, err := s.userStorage.FindById(id)

	if err != nil {
		return nil, err
	}

	if userModel == nil {
		return nil, nil
	}

	return CreateUser(userModel), nil
}

func (s UserService) FindUserByEmail(email string) (*User, error) {
	userModel, err := s.userStorage.FindByEmail(email)

	if err != nil {
		return nil, err
	}

	if userModel == nil {
		return nil, nil
	}

	return CreateUser(userModel), nil
}

func (s UserService) ChangePassword(userID int, password, newPassword string) error {
	user, err := s.userStorage.FindById(userID)

	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	if !user.CheckPassword(password) {
		return ErrInvalidPassword
	}

	if err = user.SetPassword(newPassword); err != nil {
		return err
	}

	return s.userStorage.Save(user)
}

func NewUserService(userStorage storage.UserStorage) UserService {
	return UserService{
		userStorage: userStorage,
	}
}

func CreateUser(userModel *storage.User) *User {
	return &User{
		ID:             int(userModel.ID),
		Email:          userModel.Email,
		Active:         userModel.IsActive(),
		AccountID:      int(userModel.AccountID),
		CreatedAt:      userModel.CreatedAt,
		IsAccountOwner: userModel.IsAccountOwner(),
		Account: accountService.Account{
			ID:        userModel.Account.ID,
			Confirmed: userModel.Account.Confirmed,
			CreatedAt: userModel.Account.CreatedAt,
		},
	}
}
