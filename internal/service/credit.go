package service

import (
	"SB/internal/db"
	"SB/internal/errs"
	"SB/internal/models"
	"SB/internal/repository"
	"errors"
	"fmt"
	"time"
)

// Валидация ISO валюты (простейшая)
//func isValidCurrency(currency string) bool {
//	validCurrencies := map[string]bool{"USD": true, "EUR": true, "RUB": true}
//	return validCurrencies[currency]
//}

// Максимальные ограничения для кредитов (пример)
const (
	MaxCreditAmount   = 1_000_000
	MaxCreditDuration = 60 // месяцев
	MinInterestRate   = 0.0
	MaxInterestRate   = 100.0
)

// Создать заявку на кредит
func CreateCredit(credit *models.Credit) (int, error) {
	user, err := repository.GetUserByID(credit.UserID)
	if err != nil {
		return 0, errs.ErrNotFound
	}
	if !user.Active {
		return 0, errs.ErrUserNotActive
	}

	if credit.Amount <= 0 || credit.Amount > MaxCreditAmount {
		return 0, fmt.Errorf("amount must be > 0 and <= %d", MaxCreditAmount)
	}
	if !IsValidCurrency(credit.Currency) {
		return 0, fmt.Errorf("invalid currency")
	}
	if credit.DurationMonths <= 0 || credit.DurationMonths > MaxCreditDuration {
		return 0, fmt.Errorf("duration_months must be > 0 and <= %d", MaxCreditDuration)
	}
	if credit.InterestRate < MinInterestRate || credit.InterestRate > MaxInterestRate {
		return 0, fmt.Errorf("interest_rate must be between %.2f and %.2f", MinInterestRate, MaxInterestRate)
	}

	credit.Active = true
	credit.CreatedAt = time.Now()

	return repository.CreateCredit(credit)
}

// Обновить кредит (только если активен и не одобрен)
func UpdateCredit(credit *models.Credit) error {
	existing, err := repository.GetCreditByID(credit.ID)
	if err != nil {
		return errs.ErrNotFound
	}
	if !existing.Active {
		return errs.ErrCreditNotActive
	}
	if existing.ApprovedAt != nil {
		return errors.New("approved credits cannot be updated")
	}

	if credit.Amount <= 0 || credit.Amount > MaxCreditAmount {
		return fmt.Errorf("amount must be > 0 and <= %d", MaxCreditAmount)
	}
	if !IsValidCurrency(credit.Currency) {
		return fmt.Errorf("invalid currency")
	}
	if credit.DurationMonths <= 0 || credit.DurationMonths > MaxCreditDuration {
		return fmt.Errorf("duration_months must be > 0 and <= %d", MaxCreditDuration)
	}
	if credit.InterestRate < MinInterestRate || credit.InterestRate > MaxInterestRate {
		return fmt.Errorf("interest_rate must be between %.2f and %.2f", MinInterestRate, MaxInterestRate)
	}

	return repository.UpdateCredit(credit)
}

// Одобрить кредит (админ)
//func ApproveCredit(creditID int) error {
//	credit, err := repository.GetCreditByID(creditID)
//	if err != nil {
//		return errs.ErrNotFound
//	}
//	if !credit.Active {
//		return errs.ErrCreditNotActive
//	}
//	if credit.ApprovedAt != nil {
//		return errors.New("credit already approved")
//	}
//
//	accounts, err := repository.GetAccountsByUserID(credit.UserID)
//	if err != nil {
//		return err
//	}
//
//	var targetAccount *models.Account
//	for i := range accounts {
//		if accounts[i].Currency == credit.Currency && accounts[i].Active {
//			targetAccount = &accounts[i]
//			break
//		}
//	}
//	if targetAccount == nil {
//		return errors.New("no active account found for user with matching currency")
//	}
//
//	tx, err := db.GetDBConn().Begin()
//	if err != nil {
//		return err
//	}
//	defer func() {
//		if err != nil {
//			_ = tx.Rollback()
//		}
//	}()
//
//	now := time.Now()
//	_, err = tx.Exec(`
//		UPDATE credits SET approved_at = $1 WHERE id = $2`,
//		now, credit.ID)
//	if err != nil {
//		return err
//	}
//
//	newBalance := targetAccount.Balance + credit.Amount
//	_, err = tx.Exec(`
//		UPDATE accounts SET balance = $1, updated_at = $2 WHERE id = $3`,
//		newBalance, now, targetAccount.ID)
//	if err != nil {
//		return err
//	}
//
//	err = tx.Commit()
//	return err
//}

// Получить кредит по ID (только если активен)
func GetCreditByID(creditID int) (*models.Credit, error) {
	credit, err := repository.GetCreditByID(creditID)
	if err != nil {
		return nil, errs.ErrNotFound
	}
	if !credit.Active {
		return nil, errs.ErrCreditNotActive
	}
	return &credit, nil
}

// Получить кредиты пользователя (только активные)
func GetCreditsByUserID(userID int) ([]models.Credit, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, errs.ErrNotFound
	}
	if !user.Active {
		return nil, errs.ErrUserNotActive
	}
	return repository.GetCreditsByUserID(userID)
}

// Получить активные кредиты (админ)
func GetActiveCredits() ([]models.Credit, error) {
	return repository.GetActiveCredits()
}

// Получить неактивные кредиты (админ)
func GetInactiveCredits() ([]models.Credit, error) {
	return repository.GetInactiveCredits()
}

// Получить кредиты по валюте
func GetCreditsByCurrency(currency string) ([]models.Credit, error) {
	if !IsValidCurrency(currency) {
		return nil, fmt.Errorf("invalid currency")
	}
	return repository.GetCreditsByCurrency(currency)
}

// Погашение кредита
func RepayCredit(creditID, accountID int, amountToPay int64) error {
	credit, err := repository.GetCreditByID(creditID)
	if err != nil {
		return errs.ErrNotFound
	}
	if !credit.Active {
		return errs.ErrCreditNotActive
	}
	if credit.Amount <= 0 {
		return errors.New("credit amount invalid")
	}

	account, err := repository.GetAccountByID(accountID)
	if err != nil {
		return errs.ErrNotFound
	}
	if !account.Active {
		return errs.ErrAccountNotActive
	}
	if account.Balance < amountToPay {
		return errors.New("insufficient funds on account")
	}

	remaining := credit.Amount - amountToPay
	if remaining < 0 {
		return errors.New("overpayment not allowed")
	}

	tx, err := db.GetDBConn().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	newBalance := account.Balance - amountToPay
	_, err = tx.Exec(`UPDATE accounts SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, newBalance, account.ID)
	if err != nil {
		return err
	}

	newActive := remaining > 0
	_, err = tx.Exec(`UPDATE credits SET amount = $1, active = $2 WHERE id = $3`, remaining, newActive, credit.ID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// График платежей - аннуитетный метод
func CalculatePaymentSchedule(amount int64, months int, rate float64) ([]PaymentScheduleEntry, error) {
	if amount <= 0 || months <= 0 || rate < 0 {
		return nil, errors.New("invalid parameters for schedule calculation")
	}

	rateMonthly := rate / 100 / 12
	var annuity float64

	if rateMonthly == 0 {
		annuity = float64(amount) / float64(months)
	} else {
		annuity = float64(amount) * rateMonthly / (1 - (1 / pow(1+rateMonthly, months)))
	}

	var schedule []PaymentScheduleEntry
	for i := 1; i <= months; i++ {
		date := time.Now().AddDate(0, i, 0)
		interest := float64(amount) * rateMonthly
		principal := annuity - interest
		schedule = append(schedule, PaymentScheduleEntry{
			Month:     i,
			DueDate:   date,
			Payment:   annuity,
			Principal: principal,
			Interest:  interest,
		})
	}

	return schedule, nil
}

type PaymentScheduleEntry struct {
	Month     int
	DueDate   time.Time
	Payment   float64
	Principal float64
	Interest  float64
}

func pow(x float64, y int) float64 {
	result := 1.0
	for i := 0; i < y; i++ {
		result *= x
	}
	return result
}
