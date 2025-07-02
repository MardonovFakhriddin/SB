package utils

import "SB/internal/models"

func IsUserActive(user *models.User) bool {
	return user.Active && user.DeletedAt == nil
}

func IsAccountActive(account *models.Account) bool {
	return account.Active && account.DeletedAt == nil
}

func IsCreditActive(credit *models.Credit) bool {
	return credit.Active
}

func IsDepositActive(deposit *models.Deposit) bool {
	return deposit.Active
}
