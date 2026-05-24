package router

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/service"
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

// passkeyCredentialDTO is a sanitized credential response that excludes cryptographic material
type passkeyCredentialDTO struct {
	ID         uint   `json:"id"`
	DeviceName string `json:"device_name"`
	CreatedAt  string `json:"created_at"`
}

func setupPasskeyRoutes(public *gin.RouterGroup, protected *gin.RouterGroup, cfg *config.Config) {
	public.POST("/passkey/register/begin", passkeyRegisterBeginHandler(cfg))
	public.POST("/passkey/register/finish", passkeyRegisterFinishHandler())
	public.POST("/passkey/login/begin", passkeyLoginBeginHandler(cfg))
	public.POST("/passkey/login/finish", passkeyLoginFinishHandler(cfg))

	protected.GET("/passkey/credentials", passkeyListCredentialsHandler)
	protected.DELETE("/passkey/credentials/:id", passkeyDeleteCredentialHandler)
}

func passkeyRegisterBeginHandler(cfg *config.Config) gin.HandlerFunc {
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

		record, err := service.ValidateTempPassword(strings.TrimSpace(request.TempPassword))
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse(http.StatusUnauthorized, "invalid or expired temp password"))
			return
		}
		if err := service.CleanupExpiredTempPasswords(); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
		if record != nil {
			_ = service.DeleteTempPassword(record.ID)
		}

		user, err := service.LoadPasskeyUser()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		webAuthn := service.GetWebAuthn()
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

		sessionID, err := service.GenerateSessionToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		if err := service.StorePasskeySession(sessionID, session, int64(cfg.Passkey.TempPassword.TTL)); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(passkeyRegisterCreationResponse{SessionID: sessionID, Data: creation}, "passkey register begin"))
	}
}

func passkeyRegisterFinishHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request passkeyFinishRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		webAuthn := service.GetWebAuthn()
		if webAuthn == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, "webauthn not initialized"))
			return
		}

		session, err := service.LoadPasskeySession(request.SessionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		parsed, err := protocol.ParseCredentialCreationResponseBytes(request.Response)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		user, err := service.LoadPasskeyUser()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		credential, err := webAuthn.CreateCredential(user, *session, parsed)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		record, err := service.CreatePasskeyCredential(credential, strings.TrimSpace(request.DeviceName))
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(record, "passkey registered"))
	}
}

func passkeyLoginBeginHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		webAuthn := service.GetWebAuthn()
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

		sessionID, err := service.GenerateSessionToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		if cfg == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, "config not initialized"))
			return
		}

		if err := service.StorePasskeySession(sessionID, session, int64(cfg.Passkey.TempPassword.TTL)); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(passkeyLoginAssertionResponse{SessionID: sessionID, Data: assertion}, "passkey login begin"))
	}
}

func passkeyLoginFinishHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request passkeyFinishRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		webAuthn := service.GetWebAuthn()
		if webAuthn == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, "webauthn not initialized"))
			return
		}

		session, err := service.LoadPasskeySession(request.SessionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		parsed, err := protocol.ParseCredentialRequestResponseBytes(request.Response)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		_, credential, err := webAuthn.ValidatePasskeyLogin(service.LoadPasskeyUserByHandle, *session, parsed)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse(http.StatusUnauthorized, err.Error()))
			return
		}

		_ = service.UpdatePasskeyCredentialAuth(credential)

		token, err := service.GenerateSessionToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		if cfg == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, "config not initialized"))
			return
		}

		if err := service.StoreSessionToken(token, int64(cfg.Passkey.TokenTTL)); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(gin.H{"token": token, "token_type": "Bearer", "expires_in": cfg.Passkey.TokenTTL}, "login success"))
	}
}

func passkeyListCredentialsHandler(c *gin.Context) {
	credentials, err := service.ListPasskeyCredentials()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	dtos := make([]passkeyCredentialDTO, len(credentials))
	for i, cred := range credentials {
		dtos[i] = passkeyCredentialDTO{
			ID:         cred.ID,
			DeviceName: cred.DeviceName,
			CreatedAt:  cred.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, SuccessResponse(dtos, "passkey credentials"))
}

func passkeyDeleteCredentialHandler(c *gin.Context) {
	var id uint
	if err := parseUintParam(c, "id", &id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := service.DeletePasskeyCredential(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil, "passkey credential deleted"))
}
