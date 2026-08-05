-- Create webauthn credentials table
CREATE TABLE webauthn_credentials (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    credential_id BLOB NOT NULL UNIQUE,
    public_key BLOB NOT NULL,
    attestation_type TEXT NOT NULL,
    transport TEXT NOT NULL,
    device_name TEXT,
    user_present BOOLEAN NOT NULL DEFAULT 0,
    user_verified BOOLEAN NOT NULL DEFAULT 0,
    backup_eligible BOOLEAN NOT NULL DEFAULT 0,
    backup_state BOOLEAN NOT NULL DEFAULT 0,
    authenticator_aaguid BLOB NOT NULL,
    authenticator_sign_count INTEGER NOT NULL DEFAULT 0,
    authenticator_clone_warning BOOLEAN NOT NULL DEFAULT 0,
    authenticator_attachment TEXT NOT NULL,
    attestation_client_data_json BLOB NOT NULL,
    attestation_client_data_hash BLOB NOT NULL,
    attestation_authenticator_data BLOB NOT NULL,
    attestation_public_key_algorithm INTEGER NOT NULL,
    attestation_object BLOB NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE INDEX idx_webauthn_credentials_deleted_at ON webauthn_credentials(deleted_at);
