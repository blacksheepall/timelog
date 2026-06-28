package errs

import "errors"

// Category
var ErrInvalidParentID = errors.New("invalid parent_id: must be omitted for root category or a positive category id")

// Timelog
var ErrOngoingTimeLogExists = errors.New("ongoing timelog exists")

// Auth
var ErrAuthNoSession = errors.New("no session found")
var ErrAuthInvalidSession = errors.New("invalid session")

// Passkey
var ErrPasskeyConfigNotInitialized = errors.New("config not initialized")
var ErrPasskeyConfigIncomplete = errors.New("passkey config missing rp_id/rp_name/rp_origins")
var ErrPasskeySessionNil = errors.New("session is nil")
var ErrPasskeySessionNotFound = errors.New("session not found")
var ErrPasskeySessionInvalid = errors.New("invalid session data")
var ErrPasskeyCredentialNil = errors.New("credential is nil")
var ErrPasskeyUserNotFound = errors.New("user not found")
