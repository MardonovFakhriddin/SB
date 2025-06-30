package repository

import (
	"SB/internal/db"
	"SB/internal/models"
)

// Создать депозит
func CreateDeposit(deposit *models.Deposit) (int, error) {
	var id int
	err := db.GetDBConn().QueryRow(`
		INSERT INTO deposits (user_id, amount, currency, interest_rate, duration_months, expires_at, active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
		RETURNING id`,
		deposit.UserID, deposit.Amount, deposit.Currency, deposit.InterestRate,
		deposit.DurationMonths, deposit.ExpiresAt, deposit.Active).Scan(&id)
	return id, err
}

// Изменить только поле active
func UpdateDepositActiveStatus(id int, active bool) error {
	_, err := db.GetDBConn().Exec(`
		UPDATE deposits
		SET active = $1
		WHERE id = $2`, active, id)
	return err
}

// Взять депозит по ID (только если active = true)
func GetDepositByID(id int) (models.Deposit, error) {
	var deposit models.Deposit
	err := db.GetDBConn().Get(&deposit, `
		SELECT id, user_id, amount, currency, interest_rate, duration_months, created_at, expires_at, active
		FROM deposits
		WHERE id = $1 AND active = TRUE`, id)
	return deposit, err
}

// Взять депозиты по user_id (только если active = true)
func GetDepositsByUserID(userID int) ([]models.Deposit, error) {
	var deposits []models.Deposit
	err := db.GetDBConn().Select(&deposits, `
		SELECT id, user_id, amount, currency, interest_rate, duration_months, created_at, expires_at, active
		FROM deposits
		WHERE user_id = $1 AND active = TRUE`, userID)
	return deposits, err
}

// Взять все активные депозиты
func GetActiveDeposits() ([]models.Deposit, error) {
	var deposits []models.Deposit
	err := db.GetDBConn().Select(&deposits, `
		SELECT id, user_id, amount, currency, interest_rate, duration_months, created_at, expires_at, active
		FROM deposits
		WHERE active = TRUE`)
	return deposits, err
}

// Взять все неактивные депозиты
func GetInactiveDeposits() ([]models.Deposit, error) {
	var deposits []models.Deposit
	err := db.GetDBConn().Select(&deposits, `
		SELECT id, user_id, amount, currency, interest_rate, duration_months, created_at, expires_at, active
		FROM deposits
		WHERE active = FALSE`)
	return deposits, err
}

// Взять депозиты по валюте (currency) (только если active = true)
func GetDepositsByCurrency(currency string) ([]models.Deposit, error) {
	var deposits []models.Deposit
	err := db.GetDBConn().Select(&deposits, `
		SELECT id, user_id, amount, currency, interest_rate, duration_months, created_at, expires_at, active
		FROM deposits
		WHERE currency = $1 AND active = TRUE`, currency)
	return deposits, err
}
