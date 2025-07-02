package service

import (
	"SB/internal/errs"
	"SB/internal/models"
	"SB/internal/repository"
	"database/sql"
	"errors"
)

// Создать транзакцию (перевод денег)
func CreateTransfer(tx *models.Transfer, userID int) (*models.TransferTxResult, error) {
	fromAccount, err := repository.GetAccountByID(tx.FromAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	if fromAccount.Currency != tx.Currency {
		return nil, errs.ErrInvalidCurrency
	}

	toAccount, err := repository.GetAccountByID(tx.ToAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	if toAccount.Currency != tx.Currency {
		return nil, errs.ErrInvalidCurrency
	}

	if fromAccount.UserID != userID {
		return nil, errs.ErrFraud
	}

	if fromAccount.Balance < int64(tx.Amount) {
		return nil, errs.ErrInsufficientBalance
	}

	return repository.CreateTransaction(tx)
}
