package service

import (
	"testing"
)

func TestNewServiceCreatesNonNilService(t *testing.T) {
	svc, _ := newTestService(t)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}
