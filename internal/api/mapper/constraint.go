package mapper

import (
	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
)

func ConstraintToProto(c *domain.Constraint) *timelogv1.Constraint {
	if c == nil {
		return nil
	}
	return &timelogv1.Constraint{
		Id:                c.ID,
		Description:       c.Description,
		EndReason:         &c.EndReason,
		PunishmentQuote:   c.PunishmentQuote,
		StartDate:         FormatDate(c.StartDate),
		EndDate:           optionalString(FormatDatePtr(c.EndDate)),
		IsActive:          c.IsActive,
		MetricId:          c.MetricID,
		MetricOperator:    c.MetricOperator,
		MetricTargetValue: c.MetricTargetValue,
		CreatedAt:         optionalString(FormatTimeUTC(c.CreatedAt)),
		UpdatedAt:         optionalString(FormatTimeUTC(c.UpdatedAt)),
	}
}

func ConstraintsToProto(constraints []domain.Constraint) []*timelogv1.Constraint {
	out := make([]*timelogv1.Constraint, 0, len(constraints))
	for i := range constraints {
		out = append(out, ConstraintToProto(&constraints[i]))
	}
	return out
}

func ConstraintFromCreateRequest(req *timelogv1.CreateConstraintRequest) (*domain.Constraint, error) {
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
	return &domain.Constraint{
		Description:       req.Description,
		PunishmentQuote:   req.PunishmentQuote,
		StartDate:         startDate,
		EndDate:           endDate,
		IsActive:          true,
		MetricID:          req.MetricId,
		MetricOperator:    req.MetricOperator,
		MetricTargetValue: req.MetricTargetValue,
	}, nil
}

func ApplyConstraintUpdate(c *domain.Constraint, req *timelogv1.UpdateConstraintRequest) error {
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
	if req.MetricId != nil {
		c.MetricID = req.MetricId
	}
	if req.MetricOperator != nil {
		c.MetricOperator = req.MetricOperator
	}
	if req.MetricTargetValue != nil {
		c.MetricTargetValue = req.MetricTargetValue
	}
	return nil
}
