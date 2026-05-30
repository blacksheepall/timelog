package router

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
)

func assertGoldenJSON(t *testing.T, name string, value interface{}) {
	t.Helper()
	got, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal %s: %v", name, err)
	}
	got = append(got, '\n')
	path := filepath.Join("testdata", "contracts", name+".json")
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v", path, err)
	}
	if string(got) != string(want) {
		t.Fatalf("golden mismatch for %s\nwant:\n%s\ngot:\n%s", name, string(want), string(got))
	}
}

func TestGoldenTimelogCreateEnvelope(t *testing.T) {
	assertGoldenJSON(t, "timelog_create", SuccessResponse(&timelogv1.Timelog{
		Id:         1,
		StartTime:  "2026-05-30T12:00:00Z",
		CategoryId: 2,
		Remark:     stringPtr("implementation"),
		CreatedAt:  "2026-05-30T12:00:00Z",
		UpdatedAt:  "2026-05-30T12:30:00Z",
	}, "Time log created successfully"))
}

func TestGoldenCategoryTreeEnvelope(t *testing.T) {
	assertGoldenJSON(t, "category_tree", SuccessResponse([]*timelogv1.CategoryTreeNode{
		{
			Category: &timelogv1.Category{
				Id:          2,
				Name:        "Coding",
				Color:       "#3366ff",
				Description: "Deep work",
				Level:       0,
				SortOrder:   10,
				Path:        "/",
				CreatedAt:   "2026-05-30T12:00:00Z",
				UpdatedAt:   "2026-05-30T12:00:00Z",
			},
			Children: []*timelogv1.CategoryTreeNode{},
		},
	}, "Category tree retrieved successfully"))
}

func TestGoldenTaskStatsEnvelope(t *testing.T) {
	assertGoldenJSON(t, "task_stats", SuccessResponse(&timelogv1.TaskStats{
		TotalTasks:     4,
		CompletedTasks: 3,
		CompletionRate: 75,
	}, "Task stats retrieved successfully"))
}

func TestGoldenConstraintActiveEnvelope(t *testing.T) {
	assertGoldenJSON(t, "constraint_active", SuccessResponse(&timelogv1.Constraint{
		Id:              7,
		Description:     "No social media",
		PunishmentQuote: "Pay the price",
		StartDate:       "2026-05-01",
		IsActive:        true,
		CreatedAt:       "2026-05-01T00:00:00Z",
		UpdatedAt:       "2026-05-02T00:00:00Z",
	}, "Constraint retrieved successfully"))
}

func TestGoldenPasskeyLoginEnvelope(t *testing.T) {
	assertGoldenJSON(t, "passkey_login", SuccessResponse(&timelogv1.PasskeyLoginResponse{
		Token:     "abc",
		TokenType: "Bearer",
		ExpiresIn: 3600,
	}, "login success"))
}

func TestGoldenCommandNullEnvelope(t *testing.T) {
	assertGoldenJSON(t, "command_null", SuccessResponse(nil, "Task deleted successfully"))
}

func stringPtr(value string) *string {
	return &value
}
