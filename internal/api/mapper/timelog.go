package mapper

import (
	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
)

func TimelogToProto(t *domain.TimeLog) *timelogv1.Timelog {
	if t == nil {
		return nil
	}
	return &timelogv1.Timelog{
		Id:         t.ID,
		StartTime:  FormatTimeUTC(t.StartTime),
		EndTime:    optionalString(FormatTimeUTCPtr(t.EndTime)),
		CategoryId: t.CategoryID,
		TaskId:     t.TaskID,
		Remark:     &t.Remark,
		CreatedAt:  FormatTimeUTC(t.CreatedAt),
		UpdatedAt:  FormatTimeUTC(t.UpdatedAt),
	}
}

func TimelogsToProto(logs []domain.TimeLog) []*timelogv1.Timelog {
	out := make([]*timelogv1.Timelog, 0, len(logs))
	for i := range logs {
		out = append(out, TimelogToProto(&logs[i]))
	}
	return out
}

func TimelogFromCreateRequest(req *timelogv1.CreateTimelogRequest) (*domain.TimeLog, error) {
	if req == nil {
		return nil, nil
	}
	startTime, err := ParseTimeUTC(req.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := ParseOptionalTimeUTC(req.EndTime)
	if err != nil {
		return nil, err
	}
	return &domain.TimeLog{
		StartTime:  startTime,
		EndTime:    endTime,
		CategoryID: req.CategoryId,
		TaskID:     req.TaskId,
		Remark:     req.GetRemark(),
	}, nil
}

func ApplyTimelogUpdate(t *domain.TimeLog, req *timelogv1.UpdateTimelogRequest) error {
	if t == nil || req == nil {
		return nil
	}
	if req.StartTime != nil {
		startTime, err := ParseTimeUTC(req.GetStartTime())
		if err != nil {
			return err
		}
		t.StartTime = startTime
	}
	if req.EndTime != nil {
		endTime, err := ParseOptionalTimeUTC(req.EndTime)
		if err != nil {
			return err
		}
		t.EndTime = endTime
	}
	if req.CategoryId != nil {
		t.CategoryID = req.GetCategoryId()
	}
	if req.TaskId != nil {
		t.TaskID = req.TaskId
	}
	if req.Remark != nil {
		t.Remark = req.GetRemark()
	}
	return nil
}
