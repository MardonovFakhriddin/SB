package repository

import (
	"SB/internal/db"
	"SB/internal/models"
)

// Создать транзакцию
func CreateTransaction(tx *models.Transaction) (int, error) {
	var id int
	err := db.GetDBConn().QueryRow(`
		INSERT INTO transactions (from_account_id, to_account_id, amount, commission, created_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) RETURNING id`,
		tx.FromAccountID, tx.ToAccountID, tx.Amount, tx.Commission).Scan(&id)
	return id, err
}

// Взять транзакцию по ID
func GetTransactionByID(id int) (models.Transaction, error) {
	var tx models.Transaction
	err := db.GetDBConn().Get(&tx, `
		SELECT id, from_account_id, to_account_id, amount, commission, created_at
		FROM transactions
		WHERE id = $1`, id)
	return tx, err
}

// Взять транзакции, куда был отправлен счёт (to_account_id)
func GetTransactionsByToAccountID(toAccountID int) ([]models.Transaction, error) {
	var txs []models.Transaction
	err := db.GetDBConn().Select(&txs, `
		SELECT id, from_account_id, to_account_id, amount, commission, created_at
		FROM transactions
		WHERE to_account_id = $1`, toAccountID)
	return txs, err
}

// Взять транзакции, которые отправил счёт (from_account_id)
func GetTransactionsByFromAccountID(fromAccountID int) ([]models.Transaction, error) {
	var txs []models.Transaction
	err := db.GetDBConn().Select(&txs, `
		SELECT id, from_account_id, to_account_id, amount, commission, created_at
		FROM transactions
		WHERE from_account_id = $1`, fromAccountID)
	return txs, err
}
