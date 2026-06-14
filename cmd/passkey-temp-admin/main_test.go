package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
)

type fakeLogger struct{}

func (fakeLogger) Debug(fields ...interface{})                     {}
func (fakeLogger) Debugw(msg string, keysAndValues ...interface{}) {}
func (fakeLogger) Info(fields ...interface{})                      {}
func (fakeLogger) Infow(msg string, keysAndValues ...interface{})  {}
func (fakeLogger) Warn(fields ...interface{})                      {}
func (fakeLogger) Warnw(msg string, keysAndValues ...interface{})  {}
func (fakeLogger) Error(fields ...interface{})                     {}
func (fakeLogger) Errorw(msg string, keysAndValues ...interface{}) {}
func (fakeLogger) Fatal(fields ...interface{})                     {}
func (fakeLogger) Fatalw(msg string, keysAndValues ...interface{}) {}

var testCfg = &config.Config{}

func TestResolveTTL(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		defaultTTL int
		wantTTL    int
		wantErr    string
	}{
		{
			name:       "uses config default when ttl omitted",
			args:       nil,
			defaultTTL: 900,
			wantTTL:    900,
		},
		{
			name:       "uses positional ttl override",
			args:       []string{"1200"},
			defaultTTL: 900,
			wantTTL:    1200,
		},
		{
			name:       "rejects negative ttl",
			args:       []string{"-1"},
			defaultTTL: 900,
			wantErr:    "ttl must be >= 0",
		},
		{
			name:       "rejects extra arguments",
			args:       []string{"900", "extra"},
			defaultTTL: 900,
			wantErr:    "too many positional arguments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTTL, err := resolveTTL(tt.args, tt.defaultTTL)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveTTL returned error: %v", err)
			}
			if gotTTL != tt.wantTTL {
				t.Fatalf("expected ttl %d, got %d", tt.wantTTL, gotTTL)
			}
		})
	}
}

func setupCommandTestEnv(t *testing.T) (*service.Service, *model.Dao) {
	t.Helper()

	testCfg.Database.Host = ":memory:"
	testCfg.Log.ORMLogLevel = 1
	testCfg.Passkey.TempPassword.TTL = 900

	dao, err := model.NewDao(testCfg, fakeLogger{})
	if err != nil {
		t.Fatalf("NewDao: %v", err)
	}
	repos := adapter.NewRepositories(dao)
	svc := service.NewService(repos, repos, repos, repos, repos, repos, repos, repos, repos, fakeLogger{}, testCfg, nil)

	if err := dao.Db().AutoMigrate(&model.TempPassword{}); err != nil {
		t.Fatalf("auto migrate temp_passwords: %v", err)
	}
	if err := dao.Db().Exec("DELETE FROM temp_passwords").Error; err != nil {
		t.Fatalf("clear temp_passwords: %v", err)
	}
	return svc, dao
}

func TestRunCreateUsesDefaultTTLAndPersistsRecord(t *testing.T) {
	svc, dao := setupCommandTestEnv(t)

	var stdout bytes.Buffer
	err := runCreate(nil, svc, testCfg, &stdout)
	if err != nil {
		t.Fatalf("runCreate returned error: %v", err)
	}

	var count int64
	if err := dao.Db().Model(&model.TempPassword{}).Count(&count).Error; err != nil {
		t.Fatalf("count temp_passwords: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 temp password record, got %d", count)
	}
	if !strings.Contains(stdout.String(), "temp password:") {
		t.Fatalf("expected plaintext password in output, got %q", stdout.String())
	}
}

func TestRunCreateRejectsExtraArgs(t *testing.T) {
	svc, _ := setupCommandTestEnv(t)

	var stdout bytes.Buffer
	err := runCreate([]string{"900", "extra"}, svc, testCfg, &stdout)
	if err == nil || err.Error() != "too many positional arguments" {
		t.Fatalf("expected extra-argument error, got %v", err)
	}
}
