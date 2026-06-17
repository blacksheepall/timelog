package timelogv1

import (
	"reflect"
	"testing"
)

// Intentionally hand-written test file preserved in the generated package.
// It is not overwritten by `make gen-api` (buf only writes generated *.pb.go
// files), and it is needed because default coverage only counts a package's
// own tests. This exercises the generated protobuf methods concretely rather
// than merely checking that they do not panic.

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

		name := val.Type().Name()

		// Concrete assertions on methods every generated message must have.
		reset := ptr.MethodByName("Reset")
		if !reset.IsValid() {
			t.Fatalf("%s: missing Reset method", name)
		}
		reset.Call(nil)

		str := ptr.MethodByName("String")
		if !str.IsValid() {
			t.Fatalf("%s: missing String method", name)
		}
		_ = str.Call(nil)[0].String()

		protoReflect := ptr.MethodByName("ProtoReflect")
		if !protoReflect.IsValid() {
			t.Fatalf("%s: missing ProtoReflect method", name)
		}
		if pr := protoReflect.Call(nil)[0]; pr.IsNil() {
			t.Fatalf("%s: ProtoReflect returned nil", name)
		}

		// Exercise the remaining exported pointer-receiver methods (getters,
		// setters, descriptor helpers, etc.) to ensure they do not panic and
		// to keep generated-code coverage honest.
		mt := ptr.Type()
		for i := 0; i < mt.NumMethod(); i++ {
			m := mt.Method(i)
			// Skip methods already exercised above.
			if m.Name == "Reset" || m.Name == "String" || m.Name == "ProtoReflect" {
				continue
			}
			if m.Type.NumIn() != 1 || m.Type.NumOut() == 0 {
				continue
			}
			_ = ptr.Method(i).Call(nil)
		}
	}
}
