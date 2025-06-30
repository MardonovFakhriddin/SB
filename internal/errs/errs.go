package errs

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrAccountNotActive = errors.New("account is not active")
	ErrUserNotActive    = errors.New("user is not active")
	ErrCreditNotActive  = errors.New("credit is not active")
	ErrDepositNotActive = errors.New("deposit is not active")
)
