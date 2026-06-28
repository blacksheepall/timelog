package adapter

import (
	"context"

	"github.com/blacksheepaul/timelog/core/audit"
	"github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/model"
)

// auditRepo implements audit.Logger. It persists audit records to the database
// and also emits a structured log line so operators can grep the application log.
type auditRepo struct {
	db model.DBProvider
	l  logger.Logger
}

var _ audit.Logger = (*auditRepo)(nil)

func newAuditRepo(db model.DBProvider, l logger.Logger) *auditRepo {
	return &auditRepo{db: db, l: l}
}

func (r *auditRepo) Log(ctx context.Context, action, entityType string, entityID int32, payload map[string]any) error {
	source := audit.SourceFromContext(ctx)
	if source == "" {
		source = "unknown"
	}
	reqID := audit.RequestIDFromContext(ctx)

	// Persist to the database.
	if err := model.CreateAuditLog(r.db.Db(), action, entityType, entityID, source, reqID, payload); err != nil {
		// Log the persistence failure but do not fail the original operation.
		r.l.Errorw("failed to persist audit log",
			"error", err,
			"action", action,
			"entity_type", entityType,
			"entity_id", entityID,
			"source", source,
			"request_id", reqID,
		)
	}

	// Emit a structured audit log line for easy searching.
	r.l.WithContext(ctx).Infow("audit",
		"audit", true,
		"action", action,
		"entity_type", entityType,
		"entity_id", entityID,
		"source", source,
		"request_id", reqID,
		"payload", payload,
	)

	return nil
}
