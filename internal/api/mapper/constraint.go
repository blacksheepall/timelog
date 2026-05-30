package mapper

import (
	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/model/gen"
)

func ConstraintToProto(c *gen.Constraint) *timelogv1.Constraint {
	if c == nil {
		return nil
	}
	return &timelogv1.Constraint{
		Id:              Int32Value(c.ID),
		Description:     c.Description,
		EndReason:       c.EndReason,
		PunishmentQuote: c.PunishmentQuote,
		StartDate:       FormatDate(c.StartDate),
		EndDate:         optionalString(FormatDatePtr(c.EndDate)),
		IsActive:        BoolValue(c.IsActive),
		CreatedAt:       FormatTimeUTCPtr(c.CreatedAt),
		UpdatedAt:       FormatTimeUTCPtr(c.UpdatedAt),
	}
}

func ConstraintsToProto(constraints []gen.Constraint) []*timelogv1.Constraint {
	out := make([]*timelogv1.Constraint, 0, len(constraints))
	for i := range constraints {
		out = append(out, ConstraintToProto(&constraints[i]))
	}
	return out
}

func ConstraintFromCreateRequest(req *timelogv1.CreateConstraintRequest) (*gen.Constraint, error) {
	if req == nil {
		return nil, nil
	}
	startDate, err := ParseDate(req.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := ParseOptionalDate(req.EndDate)
	if err != nil {
		return nil, err
	}
	active := true
	return &gen.Constraint{
		Description:     req.Description,
		PunishmentQuote: req.PunishmentQuote,
		StartDate:       startDate,
		EndDate:         endDate,
		IsActive:        &active,
	}, nil
}

func ApplyConstraintUpdate(c *gen.Constraint, req *timelogv1.UpdateConstraintRequest) error {
	if c == nil || req == nil {
		return nil
	}
	if req.Description != nil {
		c.Description = req.GetDescription()
	}
	if req.PunishmentQuote != nil {
		c.PunishmentQuote = req.GetPunishmentQuote()
	}
	if req.StartDate != nil {
		startDate, err := ParseDate(req.GetStartDate())
		if err != nil {
			return err
		}
		c.StartDate = startDate
	}
	if req.EndDate != nil {
		endDate, err := ParseOptionalDate(req.EndDate)
		if err != nil {
			return err
		}
		c.EndDate = endDate
	}
	return nil
}
