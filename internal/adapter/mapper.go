package adapter

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/go-webauthn/webauthn/protocol"
)

// --- TimeLog ---

func toDomainTimelog(g *gen.Timelog) *domain.TimeLog {
	if g == nil {
		return nil
	}
	d := &domain.TimeLog{
		StartTime:  g.StartTime,
		CategoryID: g.CategoryID,
		CreatedAt:  g.CreatedAt,
		UpdatedAt:  g.UpdatedAt,
	}
	if g.ID != nil {
		d.ID = *g.ID
	}
	if g.UserID != nil {
		d.UserID = *g.UserID
	}
	if g.EndTime != nil {
		d.EndTime = g.EndTime
	}
	if g.TaskID != nil {
		d.TaskID = g.TaskID
	}
	if g.Remark != nil {
		d.Remark = *g.Remark
	}
	return d
}

func toGenTimelog(d *domain.TimeLog) *gen.Timelog {
	if d == nil {
		return nil
	}
	g := &gen.Timelog{
		StartTime:  d.StartTime,
		CategoryID: d.CategoryID,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
	if d.ID != 0 {
		id := d.ID
		g.ID = &id
	}
	if d.UserID != 0 {
		userID := d.UserID
		g.UserID = &userID
	}
	if d.EndTime != nil {
		g.EndTime = d.EndTime
	}
	if d.TaskID != nil {
		g.TaskID = d.TaskID
	}
	if d.Remark != "" {
		g.Remark = &d.Remark
	}
	return g
}

func toDomainTimelogs(list []gen.Timelog) []domain.TimeLog {
	out := make([]domain.TimeLog, len(list))
	for i := range list {
		if tl := toDomainTimelog(&list[i]); tl != nil {
			out[i] = *tl
		}
	}
	return out
}

// --- Category ---

func toDomainCategory(g *gen.Category) *domain.Category {
	if g == nil {
		return nil
	}
	d := &domain.Category{
		Name:      g.Name,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
	if g.ID != nil {
		d.ID = *g.ID
	}
	if g.Color != nil {
		d.Color = *g.Color
	}
	if g.Description != nil {
		d.Description = *g.Description
	}
	if g.ParentID != nil {
		d.ParentID = g.ParentID
	}
	if g.Level != nil {
		d.Level = *g.Level
	}
	if g.SortOrder != nil {
		d.SortOrder = *g.SortOrder
	}
	if g.Path != nil {
		d.Path = *g.Path
	}
	return d
}

func toGenCategory(d *domain.Category) *gen.Category {
	if d == nil {
		return nil
	}
	g := &gen.Category{
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
	if d.ID != 0 {
		id := d.ID
		g.ID = &id
	}
	if d.Color != "" {
		g.Color = &d.Color
	}
	if d.Description != "" {
		g.Description = &d.Description
	}
	if d.ParentID != nil {
		g.ParentID = d.ParentID
	}
	if d.Level != 0 {
		level := d.Level
		g.Level = &level
	}
	if d.SortOrder != 0 {
		so := d.SortOrder
		g.SortOrder = &so
	}
	if d.Path != "" {
		g.Path = &d.Path
	}
	return g
}

func toDomainCategories(list []gen.Category) []domain.Category {
	out := make([]domain.Category, len(list))
	for i := range list {
		if c := toDomainCategory(&list[i]); c != nil {
			out[i] = *c
		}
	}
	return out
}

func toDomainCategoryNode(n *model.CategoryNode) *domain.CategoryNode {
	if n == nil {
		return nil
	}
	node := &domain.CategoryNode{
		Category: *toDomainCategory(&n.Category),
	}
	if len(n.Children) > 0 {
		node.Children = make([]*domain.CategoryNode, len(n.Children))
		for i, child := range n.Children {
			node.Children[i] = toDomainCategoryNode(child)
		}
	}
	return node
}

func toDomainCategoryNodes(list []*model.CategoryNode) []*domain.CategoryNode {
	out := make([]*domain.CategoryNode, len(list))
	for i := range list {
		out[i] = toDomainCategoryNode(list[i])
	}
	return out
}

// --- Task ---

func toDomainTask(g *gen.Task) *domain.Task {
	if g == nil {
		return nil
	}
	d := &domain.Task{
		Title:            g.Title,
		Description:      "",
		CategoryID:       g.CategoryID,
		DueDate:          g.DueDate,
		EstimatedMinutes: g.EstimatedMinutes,
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
	}
	if g.ID != nil {
		d.ID = *g.ID
	}
	if g.Description != nil {
		d.Description = *g.Description
	}
	if g.IsCompleted != nil {
		d.IsCompleted = *g.IsCompleted
	}
	if g.CompletedAt != nil {
		d.CompletedAt = g.CompletedAt
	}
	if g.IsSuspended != nil {
		d.IsSuspended = *g.IsSuspended
	}
	if g.CreatedAt != nil {
		d.CreatedAt = *g.CreatedAt
	}
	if g.UpdatedAt != nil {
		d.UpdatedAt = *g.UpdatedAt
	}
	return d
}

func toGenTask(d *domain.Task) *gen.Task {
	if d == nil {
		return nil
	}
	g := &gen.Task{
		Title:            d.Title,
		CategoryID:       d.CategoryID,
		DueDate:          d.DueDate,
		EstimatedMinutes: d.EstimatedMinutes,
	}
	if d.ID != 0 {
		id := d.ID
		g.ID = &id
	}
	if d.Description != "" {
		g.Description = &d.Description
	}
	isComp := d.IsCompleted
	g.IsCompleted = &isComp
	if d.CompletedAt != nil {
		g.CompletedAt = d.CompletedAt
	}
	isSus := d.IsSuspended
	g.IsSuspended = &isSus
	g.CreatedAt = &d.CreatedAt
	g.UpdatedAt = &d.UpdatedAt
	return g
}

func toDomainTasks(list []gen.Task) []domain.Task {
	out := make([]domain.Task, len(list))
	for i := range list {
		if t := toDomainTask(&list[i]); t != nil {
			out[i] = *t
		}
	}
	return out
}

// --- Constraint ---

func toDomainConstraint(g *gen.Constraint) *domain.Constraint {
	if g == nil {
		return nil
	}
	d := &domain.Constraint{
		Description:     g.Description,
		PunishmentQuote: g.PunishmentQuote,
		StartDate:       g.StartDate,
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}
	if g.ID != nil {
		d.ID = *g.ID
	}
	if g.EndReason != nil {
		d.EndReason = *g.EndReason
	}
	if g.EndDate != nil {
		d.EndDate = g.EndDate
	}
	if g.IsActive != nil {
		d.IsActive = *g.IsActive
	}
	if g.MetricID != nil {
		d.MetricID = g.MetricID
	}
	if g.MetricOperator != nil {
		d.MetricOperator = g.MetricOperator
	}
	if g.MetricTargetValue != nil {
		d.MetricTargetValue = g.MetricTargetValue
	}
	if g.CreatedAt != nil {
		d.CreatedAt = *g.CreatedAt
	}
	if g.UpdatedAt != nil {
		d.UpdatedAt = *g.UpdatedAt
	}
	return d
}

func toGenConstraint(d *domain.Constraint) *gen.Constraint {
	if d == nil {
		return nil
	}
	g := &gen.Constraint{
		Description:     d.Description,
		PunishmentQuote: d.PunishmentQuote,
		StartDate:       d.StartDate,
	}
	if d.ID != 0 {
		id := d.ID
		g.ID = &id
	}
	if d.EndReason != "" {
		g.EndReason = &d.EndReason
	}
	if d.EndDate != nil {
		g.EndDate = d.EndDate
	}
	isActive := d.IsActive
	g.IsActive = &isActive
	g.MetricID = d.MetricID
	g.MetricOperator = d.MetricOperator
	g.MetricTargetValue = d.MetricTargetValue
	g.CreatedAt = &d.CreatedAt
	g.UpdatedAt = &d.UpdatedAt
	return g
}

func toDomainConstraints(list []gen.Constraint) []domain.Constraint {
	out := make([]domain.Constraint, len(list))
	for i := range list {
		if c := toDomainConstraint(&list[i]); c != nil {
			out[i] = *c
		}
	}
	return out
}

// --- Metric ---

func toDomainMetric(g *gen.Metric) *domain.Metric {
	if g == nil {
		return nil
	}
	d := &domain.Metric{
		Name:        g.Name,
		Description: "",
		MetricType:  g.MetricType,
		Unit:        g.Unit,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
	if g.ID != nil {
		d.ID = *g.ID
	}
	if g.Description != nil {
		d.Description = *g.Description
	}
	if g.CurrentValue != nil {
		d.CurrentValue = g.CurrentValue
	}
	if g.LastRecordedAt != nil {
		d.LastRecordedAt = g.LastRecordedAt
	}
	if g.Extra != nil {
		d.Extra = *g.Extra
	}
	if g.CreatedAt != nil {
		d.CreatedAt = *g.CreatedAt
	}
	if g.UpdatedAt != nil {
		d.UpdatedAt = *g.UpdatedAt
	}
	return d
}

func toGenMetric(d *domain.Metric) *gen.Metric {
	if d == nil {
		return nil
	}
	g := &gen.Metric{
		Name:       d.Name,
		MetricType: d.MetricType,
		Unit:       d.Unit,
	}
	if d.ID != 0 {
		id := d.ID
		g.ID = &id
	}
	if d.Description != "" {
		g.Description = &d.Description
	}
	if d.CurrentValue != nil {
		g.CurrentValue = d.CurrentValue
	}
	if d.LastRecordedAt != nil {
		g.LastRecordedAt = d.LastRecordedAt
	}
	if d.Extra != "" {
		g.Extra = &d.Extra
	}
	g.CreatedAt = &d.CreatedAt
	g.UpdatedAt = &d.UpdatedAt
	return g
}

func toDomainMetrics(list []gen.Metric) []domain.Metric {
	out := make([]domain.Metric, len(list))
	for i := range list {
		if m := toDomainMetric(&list[i]); m != nil {
			out[i] = *m
		}
	}
	return out
}

func toDomainMetricRecord(g *gen.MetricRecord) *domain.MetricRecord {
	if g == nil {
		return nil
	}
	d := &domain.MetricRecord{
		MetricID:  g.MetricID,
		Value:     g.Value,
		Source:    "",
		CreatedAt: time.Time{},
	}
	if g.ID != nil {
		d.ID = *g.ID
	}
	if g.Source != nil {
		d.Source = *g.Source
	}
	if g.RecordedAt != nil {
		d.RecordedAt = g.RecordedAt
	}
	if g.CreatedAt != nil {
		d.CreatedAt = *g.CreatedAt
	}
	return d
}

func toGenMetricRecord(d *domain.MetricRecord) *gen.MetricRecord {
	if d == nil {
		return nil
	}
	g := &gen.MetricRecord{
		MetricID: d.MetricID,
		Value:    d.Value,
	}
	if d.ID != 0 {
		id := d.ID
		g.ID = &id
	}
	if d.Source != "" {
		g.Source = &d.Source
	}
	if d.RecordedAt != nil {
		g.RecordedAt = d.RecordedAt
	}
	g.CreatedAt = &d.CreatedAt
	return g
}

func toDomainMetricRecords(list []gen.MetricRecord) []domain.MetricRecord {
	out := make([]domain.MetricRecord, len(list))
	for i := range list {
		if r := toDomainMetricRecord(&list[i]); r != nil {
			out[i] = *r
		}
	}
	return out
}

// --- PasskeyCredential ---

func toDomainPasskeyCredential(m *model.WebAuthnCredential) *domain.PasskeyCredential {
	if m == nil {
		return nil
	}
	return &domain.PasskeyCredential{
		ID:                            int32(m.ID),
		CredentialID:                  m.CredentialID,
		PublicKey:                     m.PublicKey,
		AttestationType:               m.AttestationType,
		Transport:                     m.Transport,
		DeviceName:                    m.DeviceName,
		UserPresent:                   m.UserPresent,
		UserVerified:                  m.UserVerified,
		BackupEligible:                m.BackupEligible,
		BackupState:                   m.BackupState,
		AuthenticatorAaguid:           m.AuthenticatorAAGUID,
		AuthenticatorSignCount:        int32(m.AuthenticatorSignCount),
		AuthenticatorCloneWarning:     m.AuthenticatorCloneWarning,
		AuthenticatorAttachment:       m.AuthenticatorAttachment,
		AttestationClientDataJSON:     m.AttestationClientDataJSON,
		AttestationClientDataHash:     m.AttestationClientDataHash,
		AttestationAuthenticatorData:  m.AttestationAuthenticatorData,
		AttestationPublicKeyAlgorithm: int32(m.AttestationPublicKeyAlgorithm),
		AttestationObject:             m.AttestationObject,
		CreatedAt:                     m.CreatedAt,
		UpdatedAt:                     m.UpdatedAt,
	}
}

func toModelPasskeyCredential(d *domain.PasskeyCredential) *model.WebAuthnCredential {
	if d == nil {
		return nil
	}
	return &model.WebAuthnCredential{
		ID:                            uint(d.ID),
		CredentialID:                  d.CredentialID,
		PublicKey:                     d.PublicKey,
		AttestationType:               d.AttestationType,
		Transport:                     d.Transport,
		DeviceName:                    d.DeviceName,
		UserPresent:                   d.UserPresent,
		UserVerified:                  d.UserVerified,
		BackupEligible:                d.BackupEligible,
		BackupState:                   d.BackupState,
		AuthenticatorAAGUID:           d.AuthenticatorAaguid,
		AuthenticatorSignCount:        uint32(d.AuthenticatorSignCount),
		AuthenticatorCloneWarning:     d.AuthenticatorCloneWarning,
		AuthenticatorAttachment:       d.AuthenticatorAttachment,
		AttestationClientDataJSON:     d.AttestationClientDataJSON,
		AttestationClientDataHash:     d.AttestationClientDataHash,
		AttestationAuthenticatorData:  d.AttestationAuthenticatorData,
		AttestationPublicKeyAlgorithm: int64(d.AttestationPublicKeyAlgorithm),
		AttestationObject:             d.AttestationObject,
		CreatedAt:                     d.CreatedAt,
		UpdatedAt:                     d.UpdatedAt,
	}
}

func toDomainPasskeyCredentials(list []model.WebAuthnCredential) []domain.PasskeyCredential {
	out := make([]domain.PasskeyCredential, len(list))
	for i := range list {
		if c := toDomainPasskeyCredential(&list[i]); c != nil {
			out[i] = *c
		}
	}
	return out
}

// --- TempPassword ---

func toDomainTempPassword(m *model.TempPassword) *domain.TempPassword {
	if m == nil {
		return nil
	}
	return &domain.TempPassword{
		ID:           int32(m.ID),
		PasswordHash: m.PasswordHash,
		ExpiresAt:    m.ExpiresAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func toModelTempPassword(d *domain.TempPassword) *model.TempPassword {
	if d == nil {
		return nil
	}
	return &model.TempPassword{
		ID:           uint(d.ID),
		PasswordHash: d.PasswordHash,
		ExpiresAt:    d.ExpiresAt,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

func toDomainTempPasswords(list []model.TempPassword) []domain.TempPassword {
	out := make([]domain.TempPassword, len(list))
	for i := range list {
		if tp := toDomainTempPassword(&list[i]); tp != nil {
			out[i] = *tp
		}
	}
	return out
}

// --- helpers for model/webauthn_credential.go compatibility ---

func parseCredentialTransport(s string) []protocol.AuthenticatorTransport {
	if s == "" {
		return nil
	}
	return []protocol.AuthenticatorTransport{protocol.AuthenticatorTransport(s)}
}

func serializeCredentialTransport(t []protocol.AuthenticatorTransport) string {
	if len(t) == 0 {
		return ""
	}
	return string(t[0])
}

func parseAuthenticatorAttachment(s string) protocol.AuthenticatorAttachment {
	return protocol.AuthenticatorAttachment(s)
}
