package adapter

import (
	"reflect"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

func int32Ptr(v int32) *int32     { return &v }
func strPtr(v string) *string     { return &v }
func floatPtr(v float64) *float64 { return &v }
func boolPtr(v bool) *bool        { return &v }

func TestTimelogRoundTrip(t *testing.T) {
	endTime := time.Date(2026, 5, 30, 12, 30, 0, 0, time.UTC)
	taskID := int32(5)
	original := &domain.TimeLog{
		ID:         1,
		UserID:     2,
		StartTime:  time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		EndTime:    &endTime,
		CategoryID: 3,
		TaskID:     &taskID,
		Remark:     "note",
		CreatedAt:  time.Date(2026, 5, 30, 11, 0, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC),
	}
	g := toGenTimelog(original)
	back := toDomainTimelog(g)
	if !reflect.DeepEqual(original, back) {
		t.Fatalf("round trip mismatch\noriginal: %#v\nback: %#v", original, back)
	}
}

func TestTimelogNil(t *testing.T) {
	if toGenTimelog(nil) != nil || toDomainTimelog(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestTimelogZeroValues(t *testing.T) {
	original := &domain.TimeLog{
		StartTime:  time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		CategoryID: 1,
	}
	g := toGenTimelog(original)
	if g.ID != nil || g.UserID != nil || g.EndTime != nil || g.TaskID != nil || g.Remark != nil {
		t.Fatalf("expected zero values to stay nil: %#v", g)
	}
	back := toDomainTimelog(g)
	if back.ID != 0 || back.UserID != 0 || back.EndTime != nil || back.TaskID != nil || back.Remark != "" {
		t.Fatalf("expected zero values on round trip: %#v", back)
	}
}

func TestTimelogsToDomain(t *testing.T) {
	list := []gen.Timelog{
		{ID: int32Ptr(1), StartTime: time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC), CategoryID: 1},
		{ID: int32Ptr(2), StartTime: time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC), CategoryID: 1},
	}
	got := toDomainTimelogs(list)
	if len(got) != 2 || got[0].ID != 1 || got[1].ID != 2 {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestCategoryRoundTrip(t *testing.T) {
	parentID := int32(3)
	level := int32(2)
	sortOrder := int32(1)
	path := "/3"
	original := &domain.Category{
		ID:          7,
		Name:        "Coding",
		Color:       "#3366ff",
		Description: "Deep work",
		ParentID:    &parentID,
		Level:       level,
		SortOrder:   sortOrder,
		Path:        path,
		CreatedAt:   time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC),
	}
	g := toGenCategory(original)
	back := toDomainCategory(g)
	if !reflect.DeepEqual(original, back) {
		t.Fatalf("round trip mismatch\noriginal: %#v\nback: %#v", original, back)
	}
}

func TestCategoryNil(t *testing.T) {
	if toGenCategory(nil) != nil || toDomainCategory(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestCategoryZeroValues(t *testing.T) {
	original := &domain.Category{Name: "Root"}
	g := toGenCategory(original)
	if g.ID != nil || g.Color != nil || g.Description != nil || g.ParentID != nil || g.Level != nil || g.SortOrder != nil || g.Path != nil {
		t.Fatalf("expected zero values to stay nil: %#v", g)
	}
}

func TestCategoriesToDomain(t *testing.T) {
	list := []gen.Category{
		{ID: int32Ptr(1), Name: "A"},
		{ID: int32Ptr(2), Name: "B"},
	}
	got := toDomainCategories(list)
	if len(got) != 2 || got[0].Name != "A" || got[1].Name != "B" {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestCategoryNodeToDomain(t *testing.T) {
	node := &model.CategoryNode{
		Category: gen.Category{ID: int32Ptr(1), Name: "Root"},
		Children: []*model.CategoryNode{
			{Category: gen.Category{ID: int32Ptr(2), Name: "Child"}},
		},
	}
	got := toDomainCategoryNode(node)
	if got.Category.Name != "Root" || len(got.Children) != 1 || got.Children[0].Category.Name != "Child" {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestCategoryNodesToDomain(t *testing.T) {
	list := []*model.CategoryNode{
		{Category: gen.Category{ID: int32Ptr(1), Name: "A"}},
	}
	got := toDomainCategoryNodes(list)
	if len(got) != 1 || got[0].Category.Name != "A" {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestTaskRoundTrip(t *testing.T) {
	completedAt := time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)
	original := &domain.Task{
		ID:               5,
		Title:            "IDL",
		Description:      "Write tests",
		CategoryID:       3,
		DueDate:          time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC),
		EstimatedMinutes: 90,
		IsCompleted:      true,
		CompletedAt:      &completedAt,
		IsSuspended:      true,
		CreatedAt:        time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC),
	}
	g := toGenTask(original)
	back := toDomainTask(g)
	if !reflect.DeepEqual(original, back) {
		t.Fatalf("round trip mismatch\noriginal: %#v\nback: %#v", original, back)
	}
}

func TestTaskNil(t *testing.T) {
	if toGenTask(nil) != nil || toDomainTask(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestTaskZeroValues(t *testing.T) {
	original := &domain.Task{
		Title:            "Test",
		CategoryID:       1,
		DueDate:          time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC),
		EstimatedMinutes: 30,
	}
	g := toGenTask(original)
	if g.ID != nil || g.Description != nil || g.CompletedAt != nil {
		t.Fatalf("expected zero values to stay nil: %#v", g)
	}
}

func TestTasksToDomain(t *testing.T) {
	list := []gen.Task{
		{ID: int32Ptr(1), Title: "A", CategoryID: 1, DueDate: time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC), EstimatedMinutes: 30},
	}
	got := toDomainTasks(list)
	if len(got) != 1 || got[0].Title != "A" {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestConstraintRoundTrip(t *testing.T) {
	endDate := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	metricID := int32(5)
	op := "gt"
	target := 100.0
	original := &domain.Constraint{
		ID:                9,
		Description:       "No social media",
		EndReason:         "done",
		PunishmentQuote:   "Pay the price",
		StartDate:         time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		EndDate:           &endDate,
		IsActive:          false,
		MetricID:          &metricID,
		MetricOperator:    &op,
		MetricTargetValue: &target,
		CreatedAt:         time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2026, 5, 2, 8, 0, 0, 0, time.UTC),
	}
	g := toGenConstraint(original)
	back := toDomainConstraint(g)
	if !reflect.DeepEqual(original, back) {
		t.Fatalf("round trip mismatch\noriginal: %#v\nback: %#v", original, back)
	}
}

func TestConstraintNil(t *testing.T) {
	if toGenConstraint(nil) != nil || toDomainConstraint(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestConstraintZeroValues(t *testing.T) {
	original := &domain.Constraint{
		Description: "Test",
		StartDate:   time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
	}
	g := toGenConstraint(original)
	if g.ID != nil || g.EndReason != nil || g.EndDate != nil || g.MetricID != nil || g.MetricOperator != nil || g.MetricTargetValue != nil {
		t.Fatalf("expected zero values to stay nil: %#v", g)
	}
}

func TestConstraintsToDomain(t *testing.T) {
	list := []gen.Constraint{
		{ID: int32Ptr(1), Description: "A", StartDate: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)},
	}
	got := toDomainConstraints(list)
	if len(got) != 1 || got[0].Description != "A" {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestMetricRoundTrip(t *testing.T) {
	current := 42.0
	lastRecorded := time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)
	original := &domain.Metric{
		ID:             1,
		Name:           "Pages",
		Description:    "Pages read",
		MetricType:     "counter",
		Unit:           "page",
		CurrentValue:   &current,
		LastRecordedAt: &lastRecorded,
		Extra:          "extra",
		CreatedAt:      time.Date(2026, 5, 30, 10, 0, 0, 0, time.UTC),
		UpdatedAt:      time.Date(2026, 5, 30, 11, 0, 0, 0, time.UTC),
	}
	g := toGenMetric(original)
	back := toDomainMetric(g)
	if !reflect.DeepEqual(original, back) {
		t.Fatalf("round trip mismatch\noriginal: %#v\nback: %#v", original, back)
	}
}

func TestMetricNil(t *testing.T) {
	if toGenMetric(nil) != nil || toDomainMetric(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestMetricZeroValues(t *testing.T) {
	original := &domain.Metric{Name: "Test", MetricType: "counter", Unit: "count"}
	g := toGenMetric(original)
	if g.ID != nil || g.Description != nil || g.CurrentValue != nil || g.LastRecordedAt != nil || g.Extra != nil {
		t.Fatalf("expected zero values to stay nil: %#v", g)
	}
}

func TestMetricsToDomain(t *testing.T) {
	list := []gen.Metric{
		{ID: int32Ptr(1), Name: "A", MetricType: "counter", Unit: "count"},
	}
	got := toDomainMetrics(list)
	if len(got) != 1 || got[0].Name != "A" {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestMetricRecordRoundTrip(t *testing.T) {
	recordedAt := time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)
	original := &domain.MetricRecord{
		ID:         1,
		MetricID:   2,
		Value:      10,
		Source:     "manual",
		RecordedAt: &recordedAt,
		CreatedAt:  time.Date(2026, 5, 30, 11, 0, 0, 0, time.UTC),
	}
	g := toGenMetricRecord(original)
	back := toDomainMetricRecord(g)
	if !reflect.DeepEqual(original, back) {
		t.Fatalf("round trip mismatch\noriginal: %#v\nback: %#v", original, back)
	}
}

func TestMetricRecordNil(t *testing.T) {
	if toGenMetricRecord(nil) != nil || toDomainMetricRecord(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestMetricRecordZeroValues(t *testing.T) {
	original := &domain.MetricRecord{MetricID: 1, Value: 5}
	g := toGenMetricRecord(original)
	if g.ID != nil || g.Source != nil || g.RecordedAt != nil {
		t.Fatalf("expected zero values to stay nil: %#v", g)
	}
}

func TestMetricRecordsToDomain(t *testing.T) {
	list := []gen.MetricRecord{
		{ID: int32Ptr(1), MetricID: 1, Value: 10},
	}
	got := toDomainMetricRecords(list)
	if len(got) != 1 || got[0].Value != 10 {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestPasskeyCredentialRoundTrip(t *testing.T) {
	original := &domain.PasskeyCredential{
		ID:                            1,
		CredentialID:                  []byte("cred"),
		PublicKey:                     []byte("pub"),
		AttestationType:               "none",
		Transport:                     "usb",
		DeviceName:                    "Laptop",
		UserPresent:                   true,
		UserVerified:                  true,
		BackupEligible:                true,
		BackupState:                   false,
		AuthenticatorAaguid:           []byte("aaguid"),
		AuthenticatorSignCount:        5,
		AuthenticatorCloneWarning:     false,
		AuthenticatorAttachment:       "platform",
		AttestationClientDataJSON:     []byte("client"),
		AttestationClientDataHash:     []byte("hash"),
		AttestationAuthenticatorData:  []byte("auth"),
		AttestationPublicKeyAlgorithm: -7,
		AttestationObject:             []byte("obj"),
		CreatedAt:                     time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		UpdatedAt:                     time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC),
	}
	m := toModelPasskeyCredential(original)
	back := toDomainPasskeyCredential(m)
	if !reflect.DeepEqual(original, back) {
		t.Fatalf("round trip mismatch\noriginal: %#v\nback: %#v", original, back)
	}
}

func TestPasskeyCredentialNil(t *testing.T) {
	if toModelPasskeyCredential(nil) != nil || toDomainPasskeyCredential(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestPasskeyCredentialsToDomain(t *testing.T) {
	list := []model.WebAuthnCredential{
		{ID: 1, DeviceName: "A"},
	}
	got := toDomainPasskeyCredentials(list)
	if len(got) != 1 || got[0].DeviceName != "A" {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestTempPasswordRoundTrip(t *testing.T) {
	original := &domain.TempPassword{
		ID:           1,
		PasswordHash: "hash",
		ExpiresAt:    time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		CreatedAt:    time.Date(2026, 5, 30, 11, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2026, 5, 30, 11, 0, 0, 0, time.UTC),
	}
	m := toModelTempPassword(original)
	back := toDomainTempPassword(m)
	if !reflect.DeepEqual(original, back) {
		t.Fatalf("round trip mismatch\noriginal: %#v\nback: %#v", original, back)
	}
}

func TestTempPasswordNil(t *testing.T) {
	if toModelTempPassword(nil) != nil || toDomainTempPassword(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestTempPasswordsToDomain(t *testing.T) {
	list := []model.TempPassword{
		{ID: 1, PasswordHash: "h"},
	}
	got := toDomainTempPasswords(list)
	if len(got) != 1 || got[0].PasswordHash != "h" {
		t.Fatalf("unexpected: %#v", got)
	}
}
