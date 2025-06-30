package repository

import (
	"SB/internal/db"
	"SB/internal/errs"
	"SB/internal/models"
)

// Создать аккаунт
func CreateAccount(account *models.Account) (int, error) {
	var id int
	err := db.GetDBConn().QueryRow(`
		INSERT INTO accounts (user_id, balance, currency, active, created_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) RETURNING id`,
		account.UserID, account.Balance, account.Currency, account.Active).Scan(&id)
	return id, err
}

// Изменить аккаунт (только если active = true)
func UpdateAccount(account *models.Account) error {
	var active bool
	err := db.GetDBConn().Get(&active, `SELECT active FROM accounts WHERE id = $1 AND deleted_at IS NULL`, account.ID)
	if err != nil {
		return err
	}
	if !active {
		return errs.ErrAccountNotActive
	}

	_, err = db.GetDBConn().Exec(`
		UPDATE accounts
		SET balance = $1, currency = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND active = TRUE AND deleted_at IS NULL`,
		account.Balance, account.Currency, account.ID)
	return err
}

// Мягкое удаление аккаунта (deleted_at = now)
func DeleteAccount(id int) error {
	_, err := db.GetDBConn().Exec(`
		UPDATE accounts
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

// Взять аккаунт по ID (только если active = true)
func GetAccountByID(id int) (models.Account, error) {
	var account models.Account
	err := db.GetDBConn().Get(&account, `
		SELECT id, user_id, balance, currency, active, created_at, updated_at, deleted_at
		FROM accounts
		WHERE id = $1 AND active = TRUE AND deleted_at IS NULL`, id)
	return account, err
}

// Взять все аккаунты пользователя по user_id (только если active = true)
func GetAccountsByUserID(userID int) ([]models.Account, error) {
	var accounts []models.Account
	err := db.GetDBConn().Select(&accounts, `
		SELECT id, user_id, balance, currency, active, created_at, updated_at, deleted_at
		FROM accounts
		WHERE user_id = $1 AND active = TRUE AND deleted_at IS NULL`, userID)
	return accounts, err
}

// Взять все неактивные аккаунты
func GetInactiveAccounts() ([]models.Account, error) {
	var accounts []models.Account
	err := db.GetDBConn().Select(&accounts, `
		SELECT id, user_id, balance, currency, active, created_at, updated_at, deleted_at
		FROM accounts
		WHERE active = FALSE AND deleted_at IS NULL`)
	return accounts, err
}

// Взять аккаунты по валюте (только если active = true)
func GetAccountsByCurrency(currency string) ([]models.Account, error) {
	var accounts []models.Account
	err := db.GetDBConn().Select(&accounts, `
		SELECT id, user_id, balance, currency, active, created_at, updated_at, deleted_at
		FROM accounts
		WHERE currency = $1 AND active = TRUE AND deleted_at IS NULL`, currency)
	return accounts, err
}
