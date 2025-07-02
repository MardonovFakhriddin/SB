package repository

import "SB/internal/db"

func WriteAuditLog(action, entity string, entityID, userID int) error {
	_, err := db.GetDBConn().Exec(`INSERT INTO audit_logs (action, entity, entity_id, user_id) VALUES ($1, $2, $3, $4)`,
		action, entity, entityID, userID)
	return err
}
