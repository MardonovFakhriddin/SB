package service

import (
	"SB/internal/db"
	"fmt"
	"time"
)

// WriteAuditLog сохраняет информацию об операции в таблицу audit_logs
func WriteAuditLog(action, entity string, entityID, userID int) error {
	_, err := db.GetDBConn().Exec(`
		INSERT INTO audit_logs (action, entity, entity_id, user_id, timestamp)
		VALUES ($1, $2, $3, $4, $5)`,
		action, entity, entityID, userID, time.Now())
	return err
}

// Log formats and отправляет лог в stdout (опционально)
func Log(action, entity string, entityID, userID int) {
	msg := fmt.Sprintf("[AUDIT] %s on %s (ID: %d) by user #%d at %s",
		action, entity, entityID, userID, time.Now().Format(time.RFC3339))
	fmt.Println(msg)
}
