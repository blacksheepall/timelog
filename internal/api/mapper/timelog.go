package mapper

import (
	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/model/gen"
)

func TimelogToProto(t *gen.Timelog) *timelogv1.Timelog {
	if t == nil {
		return nil
	}
	return &timelogv1.Timelog{
		Id:         Int32Value(t.ID),
		StartTime:  FormatTimeUTC(t.StartTime),
		EndTime:    optionalString(FormatTimeUTCPtr(t.EndTime)),
		CategoryId: t.CategoryID,
		TaskId:     t.TaskID,
		Remark:     t.Remark,
		CreatedAt:  FormatTimeUTC(t.CreatedAt),
		UpdatedAt:  FormatTimeUTC(t.UpdatedAt),
	}
}

func TimelogsToProto(logs []gen.Timelog) []*timelogv1.Timelog {
	out := make([]*timelogv1.Timelog, 0, len(logs))
	for i := range logs {
		out = append(out, TimelogToProto(&logs[i]))
	}
	return out
}

func TimelogFromCreateRequest(req *timelogv1.CreateTimelogRequest) (*gen.Timelog, error) {
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
	return &gen.Timelog{
		StartTime:  startTime,
		EndTime:    endTime,
		CategoryID: req.CategoryId,
		TaskID:     req.TaskId,
		Remark:     &req.Remark,
	}, nil
}

func ApplyTimelogUpdate(t *gen.Timelog, req *timelogv1.UpdateTimelogRequest) error {
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
		remark := req.GetRemark()
		t.Remark = &remark
	}
	return nil
}
