package mapper

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
)

func TestTimelogToProtoFormatsOptionalEndTime(t *testing.T) {
	id := int32(10)
	taskID := int32(20)
	remark := "implementation"
	endTime := time.Date(2026, 5, 30, 12, 30, 0, 0, time.UTC)
	log := &gen.Timelog{
		ID:         &id,
		StartTime:  time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		EndTime:    &endTime,
		CategoryID: 3,
		TaskID:     &taskID,
		Remark:     &remark,
		CreatedAt:  time.Date(2026, 5, 30, 11, 0, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC),
	}
	got := TimelogToProto(log)
	if got.GetEndTime() != "2026-05-30T12:30:00Z" || got.GetTaskId() != 20 || got.GetRemark() != "implementation" {
		t.Fatalf("unexpected mapped timelog: %#v", got)
	}
}
