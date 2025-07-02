package logger

import (
	"log"
	"time"
)

func LogAction(action string, entity string, entityID int, performedBy int) {
	log.Printf("[AUDIT] %s - %s #%d by user #%d at %s",
		action, entity, entityID, performedBy, time.Now().Format(time.RFC3339))
}
