package service

import (
	"SB/internal/errs"
	"SB/internal/models"
	"SB/internal/repository"
	"SB/utils"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const AdminID = 1

// Хеширование пароля
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Проверка пароля
func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Создать пользователя
func CreateUser(user *models.User) error {
	if len(user.FullName) < 3 || len(user.FullName) > 50 {
		return errs.ErrInvalidFullName
	}

	user.FullName = strings.TrimSpace(user.FullName)
	if len(user.Password) < 6 {
		return errs.ErrPasswordTooShort
	}

	// Уникальность имени
	userFromDB, err := repository.GetUserByFullName(user.FullName)
	if userFromDB.ID > 0 {
		return errs.ErrUserAlreadyExists
	}

	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return errs.ErrDBUnavailable
	}

	hashed, err := hashPassword(user.Password)
	if err != nil {
		return errs.ErrCreateHash
	}

	user.Password = hashed

	user.ID, err = repository.CreateUser(user)
	return err
}

// Обновить пользователя
func UpdateUser(updateUser *models.UpdateUser) error {
	user, err := repository.GetUserByID(updateUser.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrNotFound
		}
		return err
	}

	if updateUser.FullName != nil {
		if len(user.FullName) < 3 || len(user.FullName) > 50 {
			return errs.ErrInvalidFullName
		}

		user.FullName = *updateUser.FullName
	}

	if updateUser.Password != nil {
		if len(user.Password) < 6 {
			return errs.ErrPasswordTooShort
		}

		hashed, err := hashPassword(user.Password)
		if err != nil {
			return errs.ErrCreateHash
		}

		user.Password = hashed
	}

	return repository.UpdateUser(&user)
}

// Удалить пользователя
func DeleteUser(userID int) error {

	account, err := repository.GetAccountsByUserID(userID)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return err
	}

	if account.ID > 0 {
		return errs.ErrAccountExists
	}

	credits, err := repository.GetCreditsByUserID(userID)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return err
	}

	if credits != nil {
		return errs.ErrCreditsExists
	}

	deposits, err := repository.GetDepositsByUserID(userID)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return err
	}

	if deposits != nil {
		return errs.ErrDepositsExists
	}

	return repository.DeleteUser(userID)
}

// Получить пользователя по ID
func GetUserByID(userID int) (*models.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Получить неактивных пользователей
func GetInactiveUsers() ([]models.User, error) {
	return repository.GetInactiveUsers()
}

// Аутентификация
func AuthenticateUser(fullName, password string) (string, *models.User, error) {
	user, err := repository.GetUserByFullName(fullName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, errs.ErrNotFound
		}
		return "", nil, err
	}

	if !checkPasswordHash(password, user.Password) {
		return "", nil, errs.ErrInvalidPassword
	}

	token, err := utils.GenerateToken(user.ID, user.FullName)
	if err != nil {
		return "", nil, errs.ErrGenerateToken
	}
	return token, user, nil
}

// Восстановить пользователя
func RestoreUser(fullName string) error {
	if len(fullName) < 3 || len(fullName) > 50 {
		return errs.ErrInvalidFullName
	}

	user, err := repository.GetUserByNameForRestore(fullName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrNotFound
		}
		return err
	}

	return repository.RestoreUser(user.ID)
}

// Поиск по имени
func FindUserByName(name string) ([]models.User, error) {
	if len(name) < 3 || len(name) > 50 {
		return nil, errs.ErrInvalidFullName
	}

	users, err := repository.GetUserByNameFilter(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return users, nil
}
