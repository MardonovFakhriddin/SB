package service

import (
	"SB/internal/errs"
	"SB/internal/models"
	"SB/internal/repository"
	"database/sql"
	"errors"
)

// Валидные валюты ISO 4217 (пример)
var validCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"RUB": true,
}

// Проверка валидности валюты
func IsValidCurrency(currency string) bool {
	return validCurrencies[currency]
}

// Создать аккаунт
func CreateAccount(account *models.Account) error {
	// Проверка, что пользователь существует и активен
	_, err := repository.GetUserByID(account.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrNotFound
		}
		return err
	}

	// Валидация валюты
	if !IsValidCurrency(account.Currency) {
		return errs.ErrInvalidCurrency
	}

	// Проверка лимита аккаунтов
	accountExists, err := repository.GetAccountsByUserID(account.UserID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	if accountExists.ID > 0 {
		return errs.ErrAccountAlreadyExists
	}

	// Создаем аккаунт через репозиторий
	return repository.CreateAccount(account)
}

// Обновить аккаунт
func UpdateAccount(account *models.UpdateAccount) (*models.Account, error) {
	existing, err := repository.GetAccountByID(account.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	if existing.UserID != account.UserID && account.UserID != AdminID {
		return nil, errs.ErrFraud
	}

	if account.Currency != nil {
		if !IsValidCurrency(*account.Currency) {
			return nil, errs.ErrInvalidCurrency
		}
		existing.Currency = *account.Currency
	}

	if account.PhoneNumber != nil {
		existing.PhoneNumber = *account.PhoneNumber
	}

	if account.Balance != nil {
		existing.Balance = *account.Balance
	}

	return &existing, repository.UpdateAccount(&existing)
}

// Мягкое удаление аккаунта
func DeleteAccount(accountID int, userID int) error {
	account, err := repository.GetAccountByID(accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrNotFound
		}
		return err
	}

	if account.Balance != 0 {
		return errs.ErrNoZeroBalance
	}

	if account.UserID != userID && userID != AdminID {
		return errs.ErrFraud
	}

	return repository.DeleteAccount(accountID)
}

// Получить аккаунт по ID
func GetAccountByID(accountID int) (*models.Account, error) {
	account, err := repository.GetAccountByID(accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return &account, nil
}

// Получить аккаунты пользователя
func GetAccountByUserID(userID int) (*models.Account, error) {
	_, err := repository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return repository.GetAccountsByUserID(userID)
}

// Получить неактивные аккаунты (только для админа)
func GetInactiveAccounts() ([]models.Account, error) {
	return repository.GetInactiveAccounts()
}

// Получить аккаунты по валюте
func GetAccountsByCurrency(currency string) ([]models.Account, error) {
	if !IsValidCurrency(currency) {
		return nil, errs.ErrInvalidCurrency
	}
	return repository.GetAccountsByCurrency(currency)
}

// Получить баланс аккаунта (с проверкой прав)
func GetAccountBalance(accountID int, requesterUserID int, isAdmin bool) (int64, error) {
	account, err := repository.GetAccountByID(accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errs.ErrNotFound
		}
		return 0, err
	}

	if account.UserID != requesterUserID && !isAdmin {
		return 0, errs.ErrFraud
	}
	return account.Balance, nil
}
