package mapper

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestConstraintToProtoFormatsDates(t *testing.T) {
	c := &domain.Constraint{
		ID:              9,
		Description:     "No social media",
		EndReason:       "done",
		PunishmentQuote: "Pay the price",
		StartDate:       time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		IsActive:        false,
		CreatedAt:       time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2026, 5, 2, 8, 0, 0, 0, time.UTC),
	}
	got := ConstraintToProto(c)
	if got.GetStartDate() != "2026-05-01" || got.GetEndReason() != "done" || got.GetIsActive() {
		t.Fatalf("unexpected mapped constraint: %#v", got)
	}
}
