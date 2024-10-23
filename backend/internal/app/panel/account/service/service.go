package service

import "backend/internal/app/panel/account/storage"

type AccountService struct {
	accountStorage storage.AccountStorage
}

func (s *AccountService) FindAccount(id uint) (*Account, error) {
	modelAccount, err := s.accountStorage.FindById(id)

	if err != nil {
		return nil, err
	}

	if modelAccount == nil {
		return nil, nil
	}

	return &Account{
		ID:        modelAccount.ID,
		Confirmed: modelAccount.Confirmed,
		CreatedAt: modelAccount.CreatedAt,
	}, nil
}

func NewAccountService(accountStorage storage.AccountStorage) AccountService {
	return AccountService{
		accountStorage: accountStorage,
	}
}
