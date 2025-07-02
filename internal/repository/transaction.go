package repository

import (
	"SB/internal/db"
	"SB/internal/models"
	"SB/logger"
	"context"
)

var (
	createTransferQuery = `insert into 
    transactions (from_account_id, to_account_id, amount, currency) 
	values ($1, $2, $3, $4) 
	returning id, from_account_id, to_account_id, amount, created_at`

	createEntryQuery = `insert into 
    entries (account_id, amount) 
	values ($1, $2) 
	returning id, account_id, amount, created_at`

	updateAccountQuery = `update accounts 
	set balance = balance + $2 
	where id = $1
	and active = true 
	returning id, user_id, balance, currency, created_at`
)

// Создать транзакцию
func CreateTransaction(trnx *models.Transfer) (*models.TransferTxResult, error) {
	const op = "CreateTransaction"

	tx, err := db.GetDBConn().BeginTxx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.Error.Printf("%s: failed to rollback transaction: %v", op, rollbackErr)
			}
		}
	}()

	var transfer models.Transfer
	err = tx.QueryRow(
		createTransferQuery,
		trnx.FromAccountID,
		trnx.ToAccountID,
		trnx.Amount,
		trnx.Currency,
	).Scan(
		&transfer.ID,
		&transfer.FromAccountID,
		&transfer.ToAccountID,
		&transfer.Amount,
		&transfer.CreatedAt)
	if err != nil {
		return nil, err
	}

	transfer.Currency = trnx.Currency

	var entry1 models.Entry
	err = tx.QueryRow(
		createEntryQuery,
		trnx.FromAccountID,
		-trnx.Amount,
	).Scan(
		&entry1.ID,
		&entry1.AccountID,
		&entry1.Amount,
		&entry1.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	var entry2 models.Entry
	err = tx.QueryRow(
		createEntryQuery,
		trnx.ToAccountID,
		trnx.Amount,
	).Scan(
		&entry2.ID,
		&entry2.AccountID,
		&entry2.Amount,
		&entry2.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	var (
		fromAccount models.Account
		toAccount   models.Account
	)

	if trnx.FromAccountID < trnx.ToAccountID {
		err = tx.QueryRow(
			updateAccountQuery,
			trnx.FromAccountID,
			-trnx.Amount,
		).Scan(
			&fromAccount.ID,
			&fromAccount.UserID,
			&fromAccount.Balance,
			&fromAccount.Currency,
			&fromAccount.CreatedAt,
		)

		err = tx.QueryRow(
			updateAccountQuery,
			trnx.ToAccountID,
			trnx.Amount,
		).Scan(
			&toAccount.ID,
			&toAccount.UserID,
			&toAccount.Balance,
			&toAccount.Currency,
			&toAccount.CreatedAt,
		)
	} else {
		err = tx.QueryRow(
			updateAccountQuery,
			trnx.ToAccountID,
			trnx.Amount,
		).Scan(
			&toAccount.ID,
			&toAccount.UserID,
			&toAccount.Balance,
			&toAccount.Currency,
			&toAccount.CreatedAt,
		)

		err = tx.QueryRow(
			updateAccountQuery,
			trnx.FromAccountID,
			-trnx.Amount,
		).Scan(
			&fromAccount.ID,
			&fromAccount.UserID,
			&fromAccount.Balance,
			&fromAccount.Currency,
			&fromAccount.CreatedAt,
		)
	}

	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	res := &models.TransferTxResult{
		Transfer:    transfer,
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		FromEntry:   entry1,
		ToEntry:     entry2,
	}
	return res, err
}

// Взять транзакцию по ID
func GetTransactionByID(id int) (models.Transfer, error) {
	var tx models.Transfer
	err := db.GetDBConn().Get(&tx, `
		SELECT id, from_account_id, to_account_id, amount, commission, created_at
		FROM transactions
		WHERE id = $1`, id)
	return tx, err
}

// Взять транзакции, куда был отправлен счёт (to_account_id)
func GetTransactionsByToAccountID(toAccountID int) ([]models.Transfer, error) {
	var txs []models.Transfer
	err := db.GetDBConn().Select(&txs, `
		SELECT id, from_account_id, to_account_id, amount, commission, created_at
		FROM transactions
		WHERE to_account_id = $1`, toAccountID)
	return txs, err
}

// Взять транзакции, которые отправил счёт (from_account_id)
func GetTransactionsByFromAccountID(fromAccountID int) ([]models.Transfer, error) {
	var txs []models.Transfer
	err := db.GetDBConn().Select(&txs, `
		SELECT id, from_account_id, to_account_id, amount, commission, created_at
		FROM transactions
		WHERE from_account_id = $1`, fromAccountID)
	return txs, err
}
