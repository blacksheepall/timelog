package timelogv1

import (
	"reflect"
	"testing"
)

func TestGeneratedMessageMethods(t *testing.T) {
	types := []interface{}{
		Category{},
		CategoryTreeNode{},
		CompleteConstraintRequest{},
		Constraint{},
		ConstraintEvaluation{},
		CreateCategoryRequest{},
		CreateConstraintRequest{},
		CreateMetricRequest{},
		CreateTaskRequest{},
		CreateTimelogRequest{},
		DeleteMetricRequest{},
		EvaluateConstraintRequest{},
		GetMetricRequest{},
		ListCategoriesQuery{},
		ListConstraintsQuery{},
		ListMetricRecordsRequest{},
		ListMetricsRequest{},
		ListTasksQuery{},
		ListTimelogsQuery{},
		Metric{},
		MetricRecord{},
		MoveCategoryRequest{},
		PasskeyBeginPayload{},
		PasskeyCredential{},
		PasskeyFinishRequest{},
		PasskeyLoginResponse{},
		PasskeyRegisterBeginRequest{},
		Task{},
		TaskStats{},
		Timelog{},
		UpdateCategoryRequest{},
		UpdateConstraintRequest{},
		UpdateMetricRequest{},
		UpdateTaskRequest{},
		UpdateTimelogRequest{},
	}

	for _, typ := range types {
		val := reflect.ValueOf(typ)
		ptr := reflect.New(val.Type())
		ptr.Elem().Set(val)

		// Call all exported pointer-receiver methods.
		mt := ptr.Type()
		for i := 0; i < mt.NumMethod(); i++ {
			m := mt.Method(i)
			if m.Type.NumIn() != 1 || m.Type.NumOut() == 0 {
				continue
			}
			_ = ptr.Method(i).Call(nil)
		}
	}
}
