package mapper

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestTimelogToProtoFormatsOptionalEndTime(t *testing.T) {
	endTime := time.Date(2026, 5, 30, 12, 30, 0, 0, time.UTC)
	log := &domain.TimeLog{
		ID:         10,
		StartTime:  time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		EndTime:    &endTime,
		CategoryID: 3,
		TaskID:     &[]int32{20}[0],
		Remark:     "implementation",
		CreatedAt:  time.Date(2026, 5, 30, 11, 0, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC),
	}
	got := TimelogToProto(log)
	if got.GetEndTime() != "2026-05-30T12:30:00Z" || got.GetTaskId() != 20 || got.GetRemark() != "implementation" {
		t.Fatalf("unexpected mapped timelog: %#v", got)
	}
}
