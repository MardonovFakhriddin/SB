package service

import (
	"SB/internal/db"
	"SB/internal/errs"
	"SB/internal/models"
	"SB/internal/repository"
	"time"
)

func CreateDeposit(deposit *models.Deposit, fromAccountID int) (int, error) {
	user, err := repository.GetUserByID(deposit.UserID)
	if err != nil {
		return 0, errs.ErrNotFound
	}
	if !user.Active {
		return 0, errs.ErrUserNotActive
	}
	if deposit.Amount <= 0 {
		return 0, errs.ErrInvalidAmount
	}
	if deposit.InterestRate < 0 {
		return 0, errs.ErrInvalidInterestRate
	}
	if deposit.DurationMonths <= 0 {
		return 0, errs.ErrInvalidDuration
	}
	if !IsValidCurrency(deposit.Currency) {
		return 0, errs.ErrInvalidCurrency
	}

	acc, err := repository.GetAccountByID(fromAccountID)
	if err != nil {
		return 0, errs.ErrNotFound
	}
	if acc.UserID != deposit.UserID || !acc.Active {
		return 0, errs.ErrAccountNotActive
	}
	if acc.Currency != deposit.Currency {
		return 0, errs.ErrInvalidCurrency
	}
	if acc.Balance < deposit.Amount {
		return 0, errs.ErrInsufficientFunds
	}

	tx, err := db.GetDBConn().Beginx()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	deposit.Active = true
	deposit.CreatedAt = time.Now()
	deposit.ExpiresAt = deposit.CreatedAt.AddDate(0, deposit.DurationMonths, 0)

	newBalance := acc.Balance - deposit.Amount
	_, err = tx.Exec(`UPDATE accounts SET balance = $1, updated_at = now() WHERE id = $2`, newBalance, acc.ID)
	if err != nil {
		return 0, err
	}

	var depositID int
	err = tx.QueryRow(`
		INSERT INTO deposits (user_id, amount, currency, interest_rate, duration_months, expires_at, active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, true, now())
		RETURNING id`,
		deposit.UserID, deposit.Amount, deposit.Currency, deposit.InterestRate, deposit.DurationMonths, deposit.ExpiresAt,
	).Scan(&depositID)
	if err != nil {
		return 0, err
	}

	return depositID, tx.Commit()
}

func UpdateDepositStatus(depositID int, active bool) error {
	deposit, err := repository.GetDepositByID(depositID)
	if err != nil {
		return errs.ErrNotFound
	}
	if !deposit.Active {
		return errs.ErrDepositNotActive
	}
	return repository.UpdateDepositActiveStatus(depositID, active)
}

func GetDepositByID(id int) (*models.Deposit, error) {
	dep, err := repository.GetDepositByID(id)
	if err != nil {
		return nil, errs.ErrNotFound
	}
	if !dep.Active {
		return nil, errs.ErrDepositNotActive
	}
	return &dep, nil
}

func GetDepositsByUserID(userID int) ([]models.Deposit, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil || !user.Active {
		return nil, errs.ErrUserNotActive
	}
	return repository.GetDepositsByUserID(userID)
}

func GetActiveDeposits() ([]models.Deposit, error) {
	return repository.GetActiveDeposits()
}

func GetInactiveDeposits() ([]models.Deposit, error) {
	return repository.GetInactiveDeposits()
}

func GetDepositsByCurrency(currency string) ([]models.Deposit, error) {
	if !IsValidCurrency(currency) {
		return nil, errs.ErrInvalidCurrency
	}
	return repository.GetDepositsByCurrency(currency)
}

func CloseDeposit(depositID int, toAccountID int) error {
	deposit, err := repository.GetDepositByID(depositID)
	if err != nil {
		return errs.ErrNotFound
	}
	if !deposit.Active {
		return errs.ErrDepositNotActive
	}
	if time.Now().Before(deposit.ExpiresAt) {
		return errs.ErrEarlyCloseNotAllowed
	}

	acc, err := repository.GetAccountByID(toAccountID)
	if err != nil || !acc.Active {
		return errs.ErrAccountNotActive
	}
	if acc.Currency != deposit.Currency {
		return errs.ErrInvalidCurrency
	}

	interest := CalculateDepositInterest(deposit.Amount, deposit.InterestRate, deposit.DurationMonths)
	total := deposit.Amount + interest

	tx, err := db.GetDBConn().Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.Exec(`UPDATE accounts SET balance = balance + $1, updated_at = now() WHERE id = $2`, total, acc.ID)
	if err != nil {
		return err
	}

	err = repository.UpdateDepositActiveStatus(depositID, false)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func CalculateDepositInterest(amount int64, rate float64, months int) int64 {
	simpleInterest := float64(amount) * (rate / 100) * float64(months) / 12
	return int64(simpleInterest)
}
