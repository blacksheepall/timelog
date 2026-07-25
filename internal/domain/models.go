// Package domain defines pure business entities independent of persistence
// and transport concerns. These types are used by the service layer and
// the port interfaces, ensuring the core business logic is decoupled from
// GORM and protobuf.
package domain

import "time"

// --- TimeLog ---

type TimeLog struct {
	ID         int32
	UserID     int32
	StartTime  time.Time
	EndTime    *time.Time
	CategoryID int32
	TaskID     *int32
	Remark     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (t *TimeLog) Duration() time.Duration {
	if t.EndTime != nil {
		return t.EndTime.Sub(t.StartTime)
	}
	return time.Since(t.StartTime)
}

func (t *TimeLog) IsOngoing() bool {
	return t.EndTime == nil
}

// --- Category ---

type Category struct {
	ID          int32
	Name        string
	Color       string
	Description string
	ParentID    *int32
	Level       int32
	SortOrder   int32
	Path        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CategoryNode struct {
	Category Category
	Children []*CategoryNode
}

// --- Task ---

type Task struct {
	ID               int32
	Title            string
	Description      string
	CategoryID       int32
	DueDate          time.Time
	EstimatedMinutes int32
	IsCompleted      bool
	CompletedAt      *time.Time
	IsSuspended      bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// --- Constraint ---

type Constraint struct {
	ID                int32
	Description       string
	EndReason         string
	PunishmentQuote   string
	StartDate         time.Time
	EndDate           *time.Time
	IsActive          bool
	MetricID          *int32
	MetricOperator    *string
	MetricTargetValue *float64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// --- Metric ---

type Metric struct {
	ID             int32
	Name           string
	Description    string
	MetricType     string
	Unit           string
	CurrentValue   *float64
	LastRecordedAt *time.Time
	Extra          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MetricRecord struct {
	ID         int32
	MetricID   int32
	Value      float64
	Source     string
	RecordedAt *time.Time
	CreatedAt  time.Time
}

// MetricDataPoint represents a single external measurement ready to be
// recorded against a named metric.
type MetricDataPoint struct {
	MetricName string
	Value      float64
	RecordedAt time.Time
	Source     string
}

// --- Passkey / Auth ---

type PasskeyCredential struct {
	ID                            int32
	CredentialID                  []byte
	PublicKey                     []byte
	AttestationType               string
	Transport                     string
	DeviceName                    string
	UserPresent                   bool
	UserVerified                  bool
	BackupEligible                bool
	BackupState                   bool
	AuthenticatorAaguid           []byte
	AuthenticatorSignCount        int32
	AuthenticatorCloneWarning     bool
	AuthenticatorAttachment       string
	AttestationClientDataJSON     []byte
	AttestationClientDataHash     []byte
	AttestationAuthenticatorData  []byte
	AttestationPublicKeyAlgorithm int32
	AttestationObject             []byte
	CreatedAt                     time.Time
	UpdatedAt                     time.Time
}

// --- TempPassword ---

type TempPassword struct {
	ID           int32
	PasswordHash string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
