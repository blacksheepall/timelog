package router

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/api/mapper"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type passkeyFinishRequest struct {
	SessionID  string          `json:"session_id" binding:"required"`
	Response   json.RawMessage `json:"response" binding:"required"`
	DeviceName string          `json:"device_name"`
}

type passkeyRegisterBeginRequest struct {
	TempPassword string `json:"temp_password" binding:"required"`
	DeviceName   string `json:"device_name"`
}

func setupPasskeyRoutes(public *gin.RouterGroup, protected *gin.RouterGroup, cfg *config.Config, deps Dependencies) {
	public.POST("/passkey/register/begin", passkeyRegisterBeginHandler(cfg, deps))
	public.POST("/passkey/register/finish", passkeyRegisterFinishHandler(deps))
	public.POST("/passkey/login/begin", passkeyLoginBeginHandler(cfg, deps))
	public.POST("/passkey/login/finish", passkeyLoginFinishHandler(cfg, deps))

	protected.GET("/passkey/credentials", passkeyListCredentialsHandler(deps))
	protected.DELETE("/passkey/credentials/:id", passkeyDeleteCredentialHandler(deps))
}

func passkeyRegisterBeginHandler(cfg *config.Config, deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request passkeyRegisterBeginRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		if cfg == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, "config not initialized"))
			return
		}

		record, err := deps.Service.ValidateTempPassword(strings.TrimSpace(request.TempPassword))
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse(http.StatusUnauthorized, "invalid or expired temp password"))
			return
		}
		if err := deps.Service.CleanupExpiredTempPasswords(); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
		if record != nil {
			_ = deps.Service.DeleteTempPassword(record.ID)
		}

		user, err := deps.Service.LoadPasskeyUser()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		webAuthn := deps.Service.GetWebAuthn()
		if webAuthn == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, "webauthn not initialized"))
			return
		}

		options := []webauthn.RegistrationOption{
			webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementRequired),
			webauthn.WithExclusions(webauthn.Credentials(user.WebAuthnCredentials()).CredentialDescriptors()),
		}

		creation, session, err := webAuthn.BeginRegistration(user, options...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		sessionID, err := deps.Service.GenerateSessionToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		if err := deps.Service.StorePasskeySession(sessionID, session, int64(cfg.Passkey.TempPassword.TTL)); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(passkeyRegisterCreationResponse{SessionID: sessionID, Data: creation}, "passkey register begin"))
	}
}

func passkeyRegisterFinishHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request passkeyFinishRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		webAuthn := deps.Service.GetWebAuthn()
		if webAuthn == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, "webauthn not initialized"))
			return
		}

		session, err := deps.Service.LoadPasskeySession(request.SessionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		parsed, err := protocol.ParseCredentialCreationResponseBytes(request.Response)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		user, err := deps.Service.LoadPasskeyUser()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		credential, err := webAuthn.CreateCredential(user, *session, parsed)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		record, err := deps.Service.CreatePasskeyCredential(credential, strings.TrimSpace(request.DeviceName))
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(mapper.PasskeyCredentialToProto(record), "passkey registered"))
	}
}

func passkeyLoginBeginHandler(cfg *config.Config, deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		webAuthn := deps.Service.GetWebAuthn()
		if webAuthn == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, "webauthn not initialized"))
			return
		}

		assertion, session, err := webAuthn.BeginDiscoverableLogin(
			webauthn.WithUserVerification(protocol.VerificationPreferred),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		sessionID, err := deps.Service.GenerateSessionToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		if cfg == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, "config not initialized"))
			return
		}

		if err := deps.Service.StorePasskeySession(sessionID, session, int64(cfg.Passkey.TempPassword.TTL)); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(passkeyLoginAssertionResponse{SessionID: sessionID, Data: assertion}, "passkey login begin"))
	}
}

func passkeyLoginFinishHandler(cfg *config.Config, deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request passkeyFinishRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		webAuthn := deps.Service.GetWebAuthn()
		if webAuthn == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, "webauthn not initialized"))
			return
		}

		session, err := deps.Service.LoadPasskeySession(request.SessionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		parsed, err := protocol.ParseCredentialRequestResponseBytes(request.Response)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		_, credential, err := webAuthn.ValidatePasskeyLogin(deps.Service.LoadPasskeyUserByHandle, *session, parsed)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse(http.StatusUnauthorized, err.Error()))
			return
		}

		_ = deps.Service.UpdatePasskeyCredentialAuth(credential)

		token, err := deps.Service.GenerateSessionToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		if cfg == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, "config not initialized"))
			return
		}

		if err := deps.Service.StoreSessionToken(token, int64(cfg.Passkey.TokenTTL)); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(mapper.LoginResponse(token, "Bearer", int64(cfg.Passkey.TokenTTL)), "login success"))
	}
}

func passkeyListCredentialsHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		credentials, err := deps.Service.ListPasskeyCredentials()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(mapper.PasskeyCredentialsToProto(credentials), "passkey credentials"))
	}
}

func passkeyDeleteCredentialHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id uint
		if err := parseUintParam(c, "id", &id); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		if err := deps.Service.DeletePasskeyCredential(id); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(nil, "passkey credential deleted"))
	}
}
