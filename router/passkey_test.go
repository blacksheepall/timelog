package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/testutil"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

func setupPasskeyTestRouter(t *testing.T, enabled bool) (*gin.Engine, *gin.Engine, *config.Config, *service.Service) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	cfg := testutil.NewTestConfig()
	cfg.Passkey.Enabled = enabled
	cfg.Passkey.RPID = "localhost"
	cfg.Passkey.RPName = "test"
	cfg.Passkey.RPOrigins = []string{"http://localhost:8080"}
	cfg.Passkey.TokenTTL = 3600
	cfg.Passkey.TempPassword.TTL = 60

	svc := newTestServiceForRouter(t, cfg)
	if enabled {
		webAuthn, err := service.InitWebAuthnWithConfig(cfg)
		if err != nil {
			t.Fatalf("InitWebAuthnWithConfig: %v", err)
		}
		svc.SetWebAuthn(webAuthn)
	}
	deps := Dependencies{Service: svc}

	public := gin.New()
	api := public.Group("/api")
	setupAuthRoutes(api, cfg, deps)
	setupPasskeyRoutes(api, api, cfg, deps)

	protected := gin.New()
	papi := protected.Group("/api")
	papi.Use(func(c *gin.Context) {
		c.Set("user_id", int32(1))
		c.Next()
	})
	setupPasskeyRoutes(papi, papi, cfg, deps)

	return public, protected, cfg, svc
}

func TestPasskeyRegisterBeginBadJSON(t *testing.T) {
	r, _, _, _ := setupPasskeyTestRouter(t, false)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/passkey/register/begin", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPasskeyRegisterBeginWebAuthnNotInitialized(t *testing.T) {
	r, _, _, _ := setupPasskeyTestRouter(t, false)
	body, _ := json.Marshal(passkeyRegisterBeginRequest{TempPassword: "valid-temp"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/passkey/register/begin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPasskeyRegisterFinishBadJSON(t *testing.T) {
	r, _, _, _ := setupPasskeyTestRouter(t, false)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/passkey/register/finish", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPasskeyRegisterFinishWebAuthnNotInitialized(t *testing.T) {
	r, _, _, _ := setupPasskeyTestRouter(t, false)
	body, _ := json.Marshal(passkeyFinishRequest{SessionID: "s", Response: json.RawMessage(`{}`)})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/passkey/register/finish", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPasskeyLoginBeginWebAuthnNotInitialized(t *testing.T) {
	r, _, _, _ := setupPasskeyTestRouter(t, false)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/passkey/login/begin", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPasskeyLoginFinishBadJSON(t *testing.T) {
	r, _, _, _ := setupPasskeyTestRouter(t, false)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/passkey/login/finish", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPasskeyListCredentials(t *testing.T) {
	_, protected, _, _ := setupPasskeyTestRouter(t, false)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/passkey/credentials", nil)
	protected.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPasskeyDeleteCredentialInvalidID(t *testing.T) {
	_, protected, _, _ := setupPasskeyTestRouter(t, false)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/passkey/credentials/bad", nil)
	protected.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPasskeyDeleteCredentialNotFound(t *testing.T) {
	_, protected, _, _ := setupPasskeyTestRouter(t, false)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/passkey/credentials/9999", nil)
	protected.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPasskeyRegisterBeginSuccess(t *testing.T) {
	r, _, _, svc := setupPasskeyTestRouter(t, true)
	_, password, err := svc.CreateTempPassword(60)
	if err != nil {
		t.Fatalf("create temp password: %v", err)
	}

	body, _ := json.Marshal(passkeyRegisterBeginRequest{TempPassword: password})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/passkey/register/begin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPasskeyLoginBeginSuccess(t *testing.T) {
	r, _, _, _ := setupPasskeyTestRouter(t, true)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/passkey/login/begin", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPasskeyRegisterFinishInvalidSession(t *testing.T) {
	r, _, _, _ := setupPasskeyTestRouter(t, true)
	body, _ := json.Marshal(passkeyFinishRequest{SessionID: "bad", Response: json.RawMessage(`{}`)})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/passkey/register/finish", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPasskeyLoginFinishInvalidSession(t *testing.T) {
	r, _, _, _ := setupPasskeyTestRouter(t, true)
	body, _ := json.Marshal(passkeyFinishRequest{SessionID: "bad", Response: json.RawMessage(`{}`)})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/passkey/login/finish", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPasskeyRegisterFinishInvalidResponse(t *testing.T) {
	r, _, _, svc := setupPasskeyTestRouter(t, true)
	_, password, err := svc.CreateTempPassword(60)
	if err != nil {
		t.Fatalf("create temp password: %v", err)
	}

	body, _ := json.Marshal(passkeyRegisterBeginRequest{TempPassword: password})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/passkey/register/begin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("register begin expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			SessionID string `json:"session_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	finishBody, _ := json.Marshal(passkeyFinishRequest{SessionID: resp.Data.SessionID, Response: json.RawMessage(`{}`)})
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/passkey/register/finish", bytes.NewReader(finishBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPasskeyLoginFinishInvalidResponse(t *testing.T) {
	r, _, _, _ := setupPasskeyTestRouter(t, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/passkey/login/begin", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login begin expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			SessionID string `json:"session_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	finishBody, _ := json.Marshal(passkeyFinishRequest{SessionID: resp.Data.SessionID, Response: json.RawMessage(`{}`)})
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/passkey/login/finish", bytes.NewReader(finishBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 400/401, got %d: %s", w.Code, w.Body.String())
	}
}
