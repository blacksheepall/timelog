package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// AuditLog represents an auditable change to an entity in the system.
type AuditLog struct {
	ID         int32          `gorm:"primaryKey" json:"id"`
	Action     string         `gorm:"not null" json:"action"`
	EntityType string         `gorm:"not null" json:"entity_type"`
	EntityID   int32          `json:"entity_id"`
	Source     string         `gorm:"not null" json:"source"`
	RequestID  string         `json:"request_id"`
	Payload    string         `json:"payload"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default table name for AuditLog.
func (AuditLog) TableName() string {
	return "audit_logs"
}

// CreateAuditLog inserts a new audit log record.
func CreateAuditLog(db *gorm.DB, action, entityType string, entityID int32, source, requestID string, payload map[string]any) error {
	var payloadJSON string
	if len(payload) > 0 {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		payloadJSON = string(b)
	}

	return db.Create(&AuditLog{
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Source:     source,
		RequestID:  requestID,
		Payload:    payloadJSON,
		CreatedAt:  time.Now().UTC(),
	}).Error
}
