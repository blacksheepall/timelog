package mapper

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
)

func TestConstraintToProtoFormatsDates(t *testing.T) {
	id := int32(9)
	endReason := "done"
	active := false
	createdAt := time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 5, 2, 8, 0, 0, 0, time.UTC)
	c := &gen.Constraint{
		ID:              &id,
		Description:     "No social media",
		EndReason:       &endReason,
		PunishmentQuote: "Pay the price",
		StartDate:       time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		IsActive:        &active,
		CreatedAt:       &createdAt,
		UpdatedAt:       &updatedAt,
	}
	got := ConstraintToProto(c)
	if got.GetStartDate() != "2026-05-01" || got.GetEndReason() != "done" || got.GetIsActive() {
		t.Fatalf("unexpected mapped constraint: %#v", got)
	}
}
