package repository

import (
	"SB/internal/db"
	"SB/internal/errs"
	"SB/internal/models"
)

// Создать пользователя
func CreateUser(user *models.User) (int, error) {
	var id int
	err := db.GetDBConn().QueryRow(`
		INSERT INTO users (full_name, password, active, created_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP) RETURNING id`,
		user.FullName, user.Password, user.Active).Scan(&id)
	return id, err
}

// Изменить пользователя (только если active = true)
func UpdateUser(user *models.User) error {
	var active bool
	err := db.GetDBConn().Get(&active, `SELECT active FROM users WHERE id = $1 AND deleted_at IS NULL`, user.ID)
	if err != nil {
		return err
	}
	if !active {
		return errs.ErrUserNotActive
	}

	_, err = db.GetDBConn().Exec(`
		UPDATE users
		SET full_name = $1, password = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND active = TRUE AND deleted_at IS NULL`,
		user.FullName, user.Password, user.ID)
	return err
}

// Мягкое удаление пользователя (deleted_at = now)
func DeleteUser(id int) error {
	_, err := db.GetDBConn().Exec(`
		UPDATE users
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

// Взять пользователя по ID (только если active = true)
func GetUserByID(id int) (models.User, error) {
	var user models.User
	err := db.GetDBConn().Get(&user, `
		SELECT id, full_name, password, created_at, updated_at, active, deleted_at
		FROM users
		WHERE id = $1 AND active = TRUE AND deleted_at IS NULL`, id)
	return user, err
}

// Взять пользователя по имени (только если active = true)
func GetUserByFullName(fullName string) (models.User, error) {
	var user models.User
	err := db.GetDBConn().Get(&user, `
		SELECT id, full_name, password, created_at, updated_at, active, deleted_at
		FROM users
		WHERE full_name = $1 AND active = TRUE AND deleted_at IS NULL`, fullName)
	return user, err
}

// Взять всех неактивных пользователей
func GetInactiveUsers() ([]models.User, error) {
	var users []models.User
	err := db.GetDBConn().Select(&users, `
		SELECT id, full_name, password, created_at, updated_at, active, deleted_at
		FROM users
		WHERE active = FALSE AND deleted_at IS NULL`)
	return users, err
}

// Взять всех премиум пользователей
func GetPremiumUsers() ([]models.User, error) {
	var users []models.User
	err := db.GetDBConn().Select(&users, `
		SELECT id, full_name, password, created_at, updated_at, active, deleted_at
		FROM users
		WHERE active = TRUE AND deleted_at IS NULL AND is_premium = TRUE`)
	return users, err
}
