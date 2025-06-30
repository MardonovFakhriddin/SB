package repository

import (
	"SB/internal/db"
	"SB/internal/errs"
	"SB/internal/models"
)

// Создать кредит
func CreateCredit(credit *models.Credit) (int, error) {
	var id int
	err := db.GetDBConn().QueryRow(`
		INSERT INTO credits (user_id, amount, currency, duration_months, interest_rate, active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP) RETURNING id`,
		credit.UserID, credit.Amount, credit.Currency, credit.DurationMonths, credit.InterestRate, credit.Active).Scan(&id)
	return id, err
}

// Изменить кредит (только если active = true)
func UpdateCredit(credit *models.Credit) error {
	var active bool
	err := db.GetDBConn().Get(&active, `SELECT active FROM credits WHERE id = $1 AND approved_at IS NULL`, credit.ID)
	if err != nil {
		return err
	}
	if !active {
		return errs.ErrCreditNotActive
	}

	_, err = db.GetDBConn().Exec(`
		UPDATE credits
		SET amount = $1, currency = $2, duration_months = $3, interest_rate = $4, active = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6 AND active = TRUE`,
		credit.Amount, credit.Currency, credit.DurationMonths, credit.InterestRate, credit.Active, credit.ID)
	return err
}

// Взять кредит по ID (только если active = true)
func GetCreditByID(id int) (models.Credit, error) {
	var credit models.Credit
	err := db.GetDBConn().Get(&credit, `
		SELECT id, user_id, amount, currency, duration_months, interest_rate, created_at, approved_at, active
		FROM credits
		WHERE id = $1 AND active = TRUE`, id)
	return credit, err
}

// Взять кредиты по user_id (только если active = true)
func GetCreditsByUserID(userID int) ([]models.Credit, error) {
	var credits []models.Credit
	err := db.GetDBConn().Select(&credits, `
		SELECT id, user_id, amount, currency, duration_months, interest_rate, created_at, approved_at, active
		FROM credits
		WHERE user_id = $1 AND active = TRUE`, userID)
	return credits, err
}

// Взять все кредиты с active = true
func GetActiveCredits() ([]models.Credit, error) {
	var credits []models.Credit
	err := db.GetDBConn().Select(&credits, `
		SELECT id, user_id, amount, currency, duration_months, interest_rate, created_at, approved_at, active
		FROM credits
		WHERE active = TRUE`)
	return credits, err
}

// Взять все кредиты с active = false
func GetInactiveCredits() ([]models.Credit, error) {
	var credits []models.Credit
	err := db.GetDBConn().Select(&credits, `
		SELECT id, user_id, amount, currency, duration_months, interest_rate, created_at, approved_at, active
		FROM credits
		WHERE active = FALSE`)
	return credits, err
}

// Взять кредиты по валюте (currency) (только если active = true)
func GetCreditsByCurrency(currency string) ([]models.Credit, error) {
	var credits []models.Credit
	err := db.GetDBConn().Select(&credits, `
		SELECT id, user_id, amount, currency, duration_months, interest_rate, created_at, approved_at, active
		FROM credits
		WHERE currency = $1 AND active = TRUE`, currency)
	return credits, err
}
